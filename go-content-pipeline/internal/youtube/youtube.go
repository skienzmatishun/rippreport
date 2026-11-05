package youtube

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// VideoInfo contains metadata fetched from a YouTube video
type VideoInfo struct {
	VideoID      string
	Title        string
	Description  string
	ThumbnailURL string
}

// ExtractVideoID extracts the YouTube video ID from various URL formats
// Supports:
//   - https://www.youtube.com/watch?v=VIDEO_ID
//   - https://youtu.be/VIDEO_ID
//   - https://www.youtube.com/embed/VIDEO_ID
//   - https://www.youtube.com/v/VIDEO_ID
//   - https://www.youtube.com/shorts/VIDEO_ID
//   - https://www.youtube.com/live/VIDEO_ID
//   - Bare 11-character video ID
func ExtractVideoID(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("empty YouTube URL or ID")
	}

	// Check if it's already a bare video ID (11 chars, alphanumeric + - _)
	bareIDRe := regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)
	if bareIDRe.MatchString(input) {
		return input, nil
	}

	// Parse as URL
	u, err := url.Parse(input)
	if err != nil {
		return "", fmt.Errorf("invalid YouTube URL: %w", err)
	}

	// youtu.be/VIDEO_ID
	if strings.EqualFold(u.Host, "youtu.be") {
		id := strings.Trim(u.Path, "/")
		if id == "" {
			return "", fmt.Errorf("no video ID found in youtu.be URL")
		}
		return id, nil
	}

	// youtube.com variants
	if strings.Contains(strings.ToLower(u.Host), "youtube.com") {
		// /watch?v=VIDEO_ID
		if v := u.Query().Get("v"); v != "" {
			return v, nil
		}

		// /embed/VIDEO_ID, /v/VIDEO_ID, /shorts/VIDEO_ID, /live/VIDEO_ID
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) >= 2 {
			switch strings.ToLower(parts[0]) {
			case "embed", "v", "shorts", "live":
				return parts[1], nil
			}
		}

		return "", fmt.Errorf("no video ID found in YouTube URL: %s", input)
	}

	return "", fmt.Errorf("unrecognized YouTube URL format: %s", input)
}

// ytDLPOutput represents the relevant fields from yt-dlp --dump-json
type ytDLPOutput struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Thumbnail   string `json:"thumbnail"`
}

// FetchVideoInfo retrieves video metadata using yt-dlp --dump-json
func FetchVideoInfo(youtubeURL string) (*VideoInfo, error) {
	videoID, err := ExtractVideoID(youtubeURL)
	if err != nil {
		return nil, fmt.Errorf("failed to extract video ID: %w", err)
	}

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return nil, fmt.Errorf("yt-dlp not found - install with: brew install yt-dlp")
	}

	videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	cmd := exec.Command("yt-dlp",
		"--dump-json",
		"--no-download",
		"--no-playlist",
		videoURL,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp failed to fetch video info: %w", err)
	}

	var ytInfo ytDLPOutput
	if err := json.Unmarshal(output, &ytInfo); err != nil {
		return nil, fmt.Errorf("failed to parse yt-dlp output: %w", err)
	}

	return &VideoInfo{
		VideoID:      videoID,
		Title:        ytInfo.Title,
		Description:  ytInfo.Description,
		ThumbnailURL: ytInfo.Thumbnail,
	}, nil
}

// DownloadThumbnail downloads the video thumbnail and saves it as a PNG file to destPath.
// It first attempts download & PNG conversion via yt-dlp.
// If that fails, it falls back to downloading highest resolution image via HTTP and encoding to PNG.
func DownloadThumbnail(videoID, destPath string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Method 1: Try yt-dlp with thumbnail conversion
	if _, err := exec.LookPath("yt-dlp"); err == nil {
		tempDir, err := os.MkdirTemp("", "yt-thumb-*")
		if err == nil {
			defer os.RemoveAll(tempDir)
			tempOutput := filepath.Join(tempDir, "thumb.%(ext)s")
			videoURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)

			cmd := exec.Command("yt-dlp",
				"--skip-download",
				"--write-thumbnail",
				"--convert-thumbnails", "png",
				"-o", tempOutput,
				videoURL,
			)
			if err := cmd.Run(); err == nil {
				// Find converted png file in tempDir
				matches, err := filepath.Glob(filepath.Join(tempDir, "thumb*.png"))
				if err == nil && len(matches) > 0 {
					data, err := os.ReadFile(matches[0])
					if err == nil && len(data) > 0 {
						if err := os.WriteFile(destPath, data, 0644); err == nil {
							return nil
						}
					}
				}
			}
		}
	}

	// Method 2: Fallback to HTTP download and Go image PNG encoding
	thumbnailURLs := []string{
		fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID),
		fmt.Sprintf("https://img.youtube.com/vi/%s/sddefault.jpg", videoID),
		fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", videoID),
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	for _, thumbURL := range thumbnailURLs {
		resp, err := httpClient.Get(thumbURL)
		if err != nil {
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}

		img, _, decodeErr := image.Decode(resp.Body)
		resp.Body.Close()

		if decodeErr == nil {
			outFile, err := os.Create(destPath)
			if err != nil {
				return fmt.Errorf("failed to create thumbnail file %s: %w", destPath, err)
			}
			encodeErr := png.Encode(outFile, img)
			outFile.Close()
			if encodeErr == nil {
				return nil
			}
		}

		// If decode failed, re-fetch raw bytes as fallback
		rawResp, err := httpClient.Get(thumbURL)
		if err == nil && rawResp.StatusCode == http.StatusOK {
			data, err := io.ReadAll(rawResp.Body)
			rawResp.Body.Close()
			if err == nil && len(data) > 0 {
				if err := os.WriteFile(destPath, data, 0644); err == nil {
					return nil
				}
			}
		}
	}

	return fmt.Errorf("failed to download thumbnail for video %s", videoID)
}
