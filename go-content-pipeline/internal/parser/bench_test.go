package parser

import (
	"testing"
	"time"

	"github.com/rippreport/go-content-pipeline/internal/models"
)

var sampleMarkdownWithFrontMatter = `---
title: "Benchmark Post on Municipal Policy and Zoning"
date: "2024-03-15T14:30:00Z"
description: "A detailed analysis of municipal zoning regulations and economic development."
categories:
  - government
  - policy
tags:
  - zoning
  - development
  - municipal
related_articles:
  - slug: related-article-1
    title: "Prior Analysis on County Zoning"
    score: 95
    rank: 1
  - slug: related-article-2
    title: "Regional Planning Commission Report"
    score: 88
    rank: 2
  - slug: related-article-3
    title: "Public Infrastructure Funding"
    score: 82
    rank: 3
---

# Introduction

This document provides a comprehensive overview of urban development policies.
Throughout the past decade, several shifts in regulatory frameworks have altered
the landscape of municipal planning.

## Key Observations

1. Density requirements have shifted.
2. Environmental constraints are prioritized.
3. Transit connectivity remains vital.
`

// Requirement 23.1: Benchmark YAML front matter parsing
func BenchmarkParseFrontMatter(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _, _, err := ParseFrontMatter(sampleMarkdownWithFrontMatter)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Requirement 23.1: Benchmark YAML front matter serialization
func BenchmarkSerializeFrontMatter(b *testing.B) {
	fm := &models.FrontMatter{
		Title:       "Benchmark Post on Municipal Policy and Zoning",
		Date:        time.Date(2024, 3, 15, 14, 30, 0, 0, time.UTC),
		Description: "A detailed analysis of municipal zoning regulations and economic development.",
		Categories:  []string{"government", "policy"},
		Tags:        []string{"zoning", "development", "municipal"},
		RelatedArticles: []models.RelatedArticle{
			{Slug: "related-article-1", Title: "Prior Analysis on County Zoning", Score: 95, Rank: 1},
			{Slug: "related-article-2", Title: "Regional Planning Commission Report", Score: 88, Rank: 2},
			{Slug: "related-article-3", Title: "Public Infrastructure Funding", Score: 82, Rank: 3},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := SerializeFrontMatter(fm, models.DelimiterYAML)
		if err != nil {
			b.Fatal(err)
		}
	}
}
