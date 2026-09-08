package main

import (
	"fmt"
	"github.com/rippreport/go-content-pipeline/internal/models"
	"github.com/rippreport/go-content-pipeline/internal/processor"
)

func main() {
	pm := processor.NewPostManager("/Volumes/1tb/rippreport/content/p")
	
	// Test reading a real post
	postPath := models.PostPath("/Volumes/1tb/rippreport/content/p/justice-for/index.md")
	
	post, err := pm.ReadPost(postPath)
	if err != nil {
		fmt.Printf("Error reading post: %v\n", err)
		return
	}
	
	fmt.Printf("Successfully read post!\n")
	fmt.Printf("Slug: %s\n", post.Slug)
	fmt.Printf("Title: %s\n", post.FrontMatter.Title)
	fmt.Printf("Date: %s\n", post.FrontMatter.Date)
	fmt.Printf("Categories: %v\n", post.FrontMatter.Categories)
	fmt.Printf("Tags: %v\n", post.FrontMatter.Tags)
	fmt.Printf("Body length: %d characters\n", len(post.Body))
	fmt.Printf("Content hash: %s\n", post.ContentHash)
	fmt.Printf("Related articles: %d\n", len(post.FrontMatter.RelatedArticles))
	
	if len(post.FrontMatter.RelatedArticles) > 0 {
		fmt.Printf("\nFirst related article:\n")
		ra := post.FrontMatter.RelatedArticles[0]
		fmt.Printf("  Slug: %s\n", ra.Slug)
		fmt.Printf("  Title: %s\n", ra.Title)
		fmt.Printf("  Score: %d\n", ra.Score)
		fmt.Printf("  Rank: %d\n", ra.Rank)
	}
}
