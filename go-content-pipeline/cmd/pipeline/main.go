package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Go Content Pipeline v0.1.0")
	fmt.Println("High-performance Hugo content recommendation system")
	fmt.Println()
	fmt.Println("Usage: pipeline <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  generate-related  Generate related article recommendations")
	fmt.Println("  generate-hooks    Generate contextual sidebar hooks")
	fmt.Println("  pipeline          Run complete pipeline (related + hooks)")
	fmt.Println("  apply-staged      Apply staged rankings to posts")
	fmt.Println("  clear-staged      Clear staged ranking files")
	fmt.Println()
	fmt.Println("Run 'pipeline <command> --help' for more information on a command.")
	
	os.Exit(0)
}
