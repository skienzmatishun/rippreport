package refiner

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/processor"
	"github.com/rippreport/go-content-pipeline/pkg/gemini"
)

type Embedder interface {
	GenerateEmbedding(context.Context, string) ([]float32, error)
}
type Refiner struct {
	geminiClient     *gemini.Client
	postManager      *processor.PostManager
	embedder         Embedder
	prepareEmbedding func(context.Context) (Embedder, error)
}
type Config struct {
	EmbeddingCacheFile, EmbeddingModel string
	ClearCache                         bool
	GeminiAPIKey, PostsDir             string
	Embedder                           Embedder
	PrepareEmbedding                   func(context.Context) (Embedder, error)
}
type RefineOptions struct {
	SpecificSlugs        []string
	Force, DryRun, Stage bool
	Model                string
	BatchSize            int
	SimilarityThreshold  float64
}
type RefineResult struct {
	ProcessedPosts, FailedPosts, SkippedPosts int
	CollisionPairs, RewrittenHooks            int
	Duration                                  time.Duration
}

func New(cfg Config) (*Refiner, error) {
	var c *gemini.Client
	if cfg.GeminiAPIKey != "" {
		var err error
		c, err = gemini.NewClient(gemini.Config{APIKey: cfg.GeminiAPIKey})
		if err != nil {
			return nil, err
		}
	}
	r := &Refiner{geminiClient: c, postManager: processor.NewPostManager(cfg.PostsDir), embedder: cfg.Embedder, prepareEmbedding: cfg.PrepareEmbedding}
	if cfg.EmbeddingCacheFile != "" {
		cached, err := newCachedEmbedder(cfg.EmbeddingCacheFile, cfg.EmbeddingModel, cfg.ClearCache, cfg.Embedder, cfg.PrepareEmbedding)
		if err != nil {
			return nil, err
		}
		r.embedder = cached
		r.prepareEmbedding = nil // Start the server only on a cache miss.
	}
	return r, nil
}

type hookEntry struct{ path, slug, text string }

// RefineHooks polishes hooks, then compares the entire run for semantic clones.
// Successful batches are saved before continuing; source hashes allow resuming.
func (r *Refiner) RefineHooks(ctx context.Context, opts RefineOptions) (*RefineResult, error) {
	start := time.Now()
	result := &RefineResult{}
	defer func() { result.Duration = time.Since(start) }()
	threshold := opts.SimilarityThreshold
	if threshold == 0 {
		threshold = 0.85
	}
	if math.IsNaN(threshold) || threshold <= 0 || threshold >= 1 {
		return result, fmt.Errorf("similarity threshold must be between 0 and 1")
	}
	if opts.BatchSize < 0 {
		return result, fmt.Errorf("batch size must be nonnegative")
	}
	paths, err := r.postManager.DiscoverPosts(models.PostFilter{SpecificSlugs: opts.SpecificSlugs})
	if err != nil {
		return result, err
	}
	sort.Slice(paths, func(i, j int) bool { return paths[i] < paths[j] })
	docs := map[string]*processor.HooksDocument{}
	var entries []hookEntry
	var groups [][2]int
	pending := map[int]bool{}
	for _, p := range paths {
		path := filepath.Join(filepath.Dir(string(p)), "sidebar-hooks.yaml")
		outputPath := path
		if opts.Stage {
			outputPath = filepath.Join(filepath.Dir(string(p)), "sidebar-hooks.staged.yaml")
			if _, err := os.Stat(outputPath); err == nil {
				path = outputPath
			} else if !os.IsNotExist(err) {
				return result, fmt.Errorf("stat %s: %w", outputPath, err)
			}
		}
		if _, err := os.Stat(path); os.IsNotExist(err) {
			result.SkippedPosts++
			continue
		}
		doc, err := processor.LoadHooksYAML(path)
		if err != nil {
			return result, fmt.Errorf("load %s: %w", path, err)
		}
		path = outputPath
		alreadyRefined := doc.Metadata.RefinedBy == "gemini" && !opts.Force
		copied := false
		var slugs []string
		for slug, h := range doc.Hooks {
			if h == nil || strings.TrimSpace(h.LLMText) == "" {
				continue
			}
			if directPublishHook(h, doc.Metadata.ModelName) {
				if h.Text != h.LLMText {
					h.Text = h.LLMText
					copied = true
				}
				continue
			}
			if !alreadyRefined && shouldRefineHook(h, doc.Metadata.ModelName) {
				slugs = append(slugs, slug)
			}
		}
		sort.Strings(slugs)
		if copied {
			docs[path] = doc
		}
		if len(slugs) == 0 {
			if copied {
				continue
			}
			result.SkippedPosts++
			continue
		}
		begin := len(entries)
		for _, slug := range slugs {
			h := doc.Hooks[slug]
			text := h.LLMText
			if !opts.Force && h.RefinedSource == sourceHash(h) && strings.TrimSpace(h.Text) != "" {
				text = h.Text
			} else {
				pending[len(entries)] = true
			}
			entries = append(entries, hookEntry{path, slug, text})
		}
		groups = append(groups, [2]int{begin, len(entries)})
		docs[path] = doc
	}
	if len(docs) == 0 {
		return result, nil
	}
	if opts.DryRun {
		fmt.Printf("[DRY RUN] Would save %d files, polishing/comparing %d hooks (threshold > %.2f)\n", len(docs), len(entries), threshold)
		result.SkippedPosts += len(docs)
		return result, nil
	}
	if len(entries) > 0 {
		if r.geminiClient == nil {
			return result, fmt.Errorf("Gemini API key not configured")
		}
		if r.prepareEmbedding != nil {
			r.embedder, err = withProgress("Preparing embedding server", func() (Embedder, error) { return r.prepareEmbedding(ctx) })
			if err != nil {
				return result, err
			}
		}
		if r.embedder == nil {
			return result, fmt.Errorf("local embedding client not configured")
		}
		// Save direct-publication copies before remote operations.
		for path, doc := range docs {
			if err := processor.SaveHooksYAML(path, doc); err != nil {
				return result, err
			}
		}
		if opts.BatchSize > 0 {
			groups = nil
			for i := 0; i < len(entries); i += opts.BatchSize {
				end := i + opts.BatchSize
				if end > len(entries) {
					end = len(entries)
				}
				groups = append(groups, [2]int{i, end})
			}
		}
		for _, g := range groups {
			var items []gemini.HookItem
			for i := g[0]; i < g[1]; i++ {
				if pending[i] {
					items = append(items, gemini.HookItem{Slug: fmt.Sprint(i), Brief: entries[i].text})
				}
			}
			if len(items) == 0 {
				continue
			}
			refined, err := withProgress(fmt.Sprintf("Gemini polishing %d hooks", len(items)), func() ([]gemini.RefinedHookItem, error) { return r.geminiClient.RefineHooks(ctx, items, opts.Model) })
			if err != nil {
				return result, err
			}
			texts, err := validateResponse(items, refined)
			if err != nil {
				return result, err
			}
			for i := g[0]; i < g[1]; i++ {
				if pending[i] {
					entries[i].text = texts[fmt.Sprint(i)]
				}
			}
			saved := map[string]bool{}
			for i := g[0]; i < g[1]; i++ {
				if !pending[i] {
					continue
				}
				h := entries[i]
				hook := docs[h.path].Hooks[h.slug]
				hook.Text = h.text
				hook.RefinedSource = sourceHash(hook)
				saved[h.path] = true
			}
			for path := range saved {
				if err := processor.SaveHooksYAML(path, docs[path]); err != nil {
					return result, err
				}
				fmt.Printf("Saved checkpoint: %s\n", path)
			}
			fmt.Printf("Polished and saved %d hooks\n", len(items))
		}
		vectors, err := r.embedHookTexts(ctx, entries)
		if err != nil {
			return result, err
		}

		if err := r.resolveCollisions(ctx, entries, vectors, docs, threshold, opts.Model, result, func(h hookEntry) error {
			hook := docs[h.path].Hooks[h.slug]
			hook.Text = h.text
			hook.RefinedSource = sourceHash(hook)
			return processor.SaveHooksYAML(h.path, docs[h.path])
		}); err != nil {
			return result, err
		}

	}

	for _, h := range entries {
		doc := docs[h.path]
		hook := doc.Hooks[h.slug]
		hook.Text = h.text
		doc.Hooks[h.slug] = hook
	}
	refinedPaths := map[string]bool{}
	for _, h := range entries {
		refinedPaths[h.path] = true
	}
	var files []string
	for path := range docs {
		files = append(files, path)
	}
	sort.Strings(files)
	for _, path := range files {
		doc := docs[path]
		if refinedPaths[path] {
			doc.Metadata.RefinedBy = "gemini"
			doc.Metadata.RefinedAt = time.Now().UTC()
		}
		if err := processor.SaveHooksYAML(path, doc); err != nil {
			result.FailedPosts++
			return result, err
		}
		result.ProcessedPosts++
	}
	return result, nil
}
func validateResponse(items []gemini.HookItem, response []gemini.RefinedHookItem) (map[string]string, error) {
	expected := map[string]bool{}
	for _, item := range items {
		expected[item.Slug] = true
	}
	texts := map[string]string{}
	for _, item := range response {
		if !expected[item.Slug] || texts[item.Slug] != "" || strings.TrimSpace(item.Text) == "" {
			return nil, fmt.Errorf("invalid Gemini response for hook %q", item.Slug)
		}
		texts[item.Slug] = item.Text
	}
	if len(texts) != len(expected) {
		return nil, fmt.Errorf("incomplete Gemini response: got %d of %d hooks", len(texts), len(expected))
	}
	return texts, nil
}
func normalize(v []float32) ([]float64, error) {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	if sum == 0 || math.IsNaN(sum) || math.IsInf(sum, 0) {
		return nil, fmt.Errorf("invalid embedding: empty, zero, or non-finite vector")
	}
	out := make([]float64, len(v))
	for i, x := range v {
		out[i] = float64(x) / math.Sqrt(sum)
	}
	return out, nil
}
func dot(a, b []float64) float64 {
	var sum float64
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// Per-hook provenance takes precedence over the document default. Force only
// overrides prior refinement, never the direct-publication model exclusion.
func shouldRefineHook(h *models.Hook, defaultModel string) bool {
	if h == nil || strings.TrimSpace(h.LLMText) == "" {
		return false
	}
	model := strings.TrimSpace(h.ModelName)
	if model == "" {
		model = strings.TrimSpace(defaultModel)
	}
	return !directPublishModel(model)
}

func directPublishHook(h *models.Hook, defaultModel string) bool {
	model := strings.TrimSpace(h.ModelName)
	if model == "" {
		model = strings.TrimSpace(defaultModel)
	}
	return directPublishModel(model)
}
func directPublishModel(model string) bool {
	model = strings.ToLower(strings.TrimSpace(model))
	return strings.Contains(model, "gemini") || model == "recursive/flint-v1.0:free"
}

// Each rewrite is checked against every other hook, including later hooks.
// Bound retries so a stubborn collision cannot consume unlimited API calls.
const maxCollisionRewrites = 3

func (r *Refiner) resolveCollisions(ctx context.Context, entries []hookEntry, vectors [][]float64, docs map[string]*processor.HooksDocument, threshold float64, model string, result *RefineResult, checkpoints ...func(hookEntry) error) error {
	return r.resolveTargetCollisions(ctx, entries, vectors, docs, threshold, model, result, nil, checkpoints...)
}

func (r *Refiner) resolveTargetCollisions(ctx context.Context, entries []hookEntry, vectors [][]float64, docs map[string]*processor.HooksDocument, threshold float64, model string, result *RefineResult, targets map[int]bool, checkpoints ...func(hookEntry) error) error {
	rejected := make([][]string, len(entries))
	attempts := make([]int, len(entries))
	total := len(entries)
	if targets != nil {
		total = len(targets)
	}
	checked := 0
	fmt.Printf("Checking similarity for %d target hooks against %d hooks...\n", total, len(entries))
	for j := 0; j < len(entries); j++ {
		if targets != nil && !targets[j] {
			continue
		}
		checked++
		if targets != nil || checked == 1 || checked%100 == 0 || checked == total {
			fmt.Printf("Checking target %d/%d [%s / %s]\n", checked, total, entries[j].path, entries[j].slug)
		}
		limit := j
		if targets != nil {
			limit = len(entries)
		}
		first := -1
		for i := 0; i < limit; i++ {
			if i != j && dot(vectors[i], vectors[j]) > threshold {
				first = i
				break
			}
		}
		if first < 0 {
			if targets != nil {
				fmt.Println("  No collision above threshold")
			}
			continue
		}
		for {
			if err := ctx.Err(); err != nil {
				return err
			}
			if attempts[j] >= maxCollisionRewrites {
				return fmt.Errorf("hook %s in %s still exceeds similarity threshold %.2f after %d rewrites; completed work remains saved", entries[j].slug, entries[j].path, threshold, maxCollisionRewrites)
			}
			rejected[j] = append(rejected[j], entries[j].text)
			result.CollisionPairs++
			fmt.Printf("Semantic collision %d / %d: %.4f (rewrite %d/%d)\n", first, j, dot(vectors[first], vectors[j]), attempts[j]+1, maxCollisionRewrites)
			reference := gemini.HookItem{Slug: fmt.Sprint(first), Brief: docs[entries[first].path].Hooks[entries[first].slug].Brief, Text: entries[first].text}
			target := gemini.HookItem{Slug: fmt.Sprint(j), Brief: docs[entries[j].path].Hooks[entries[j].slug].Brief, Text: entries[j].text, AvoidExamples: append([]string(nil), rejected[j]...)}
			response, err := withProgress(fmt.Sprintf("Gemini rewrite %d/%d [%s / %s]", attempts[j]+1, maxCollisionRewrites, entries[j].path, entries[j].slug), func() ([]gemini.RefinedHookItem, error) {
				return r.geminiClient.RewriteCollision(ctx, reference, target, model)
			})
			if err != nil {
				return err
			}
			texts, err := validateResponse([]gemini.HookItem{target}, response)
			if err != nil {
				return err
			}
			entries[j].text = texts[target.Slug]
			for _, checkpoint := range checkpoints {
				if err := checkpoint(entries[j]); err != nil {
					return err
				}
			}
			attempts[j]++
			result.RewrittenHooks++
			v, err := withProgress("Embedding rewritten hook", func() ([]float32, error) { return r.embedder.GenerateEmbedding(ctx, entries[j].text) })
			if err != nil {
				return err
			}
			vectors[j], err = normalize(v)
			if err != nil {
				return err
			}
			if len(vectors[j]) != len(vectors[0]) {
				return fmt.Errorf("inconsistent embedding dimensions")
			}
			// Restart the scan after every rewrite: it may collide with a different hook.
			first = -1
			for i := range entries {
				if i != j && dot(vectors[i], vectors[j]) > threshold {
					first = i
					break
				}
			}
			if first < 0 {
				fmt.Println("  Rewrite passed similarity check")
				break
			}
		}
	}
	fmt.Printf("Similarity check complete: %d targets checked, %d rewrites\n", checked, result.RewrittenHooks)
	return nil
}

func sourceHash(h *models.Hook) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(h.LLMText))) }
