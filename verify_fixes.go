package main

import (
	"fmt"
	"os"

	"github.com/Jose-Merida-UVG/Compilers-1-CC3071/internal/yapar"
)

func main() {
	tests := []struct {
		name     string
		file     string
		shouldPass bool
	}{
		{"Nested Comments", "workspace/specs/test_nested_comments.yalp", true},
		{"Unclosed Comment", "workspace/specs/test_unclosed_comment.yalp", false},
		{"Multiple IGNORE", "workspace/specs/test_ignore_multiple.yalp", true},
		{"Production With Space", "workspace/specs/test_production_with_space.yalp", false},
		{"Empty Production", "workspace/specs/test_empty_production.yalp", false},
		{"Epsilon Good", "workspace/specs/test_epsilon_good.yalp", true},
	}

	fmt.Println("=== Yapar Fix Verification ===\n")

	passed := 0
	failed := 0

	for _, test := range tests {
		content, err := os.ReadFile(test.file)
		if err != nil {
			fmt.Printf("❌ SKIP: %s (file not found: %v)\n", test.name, err)
			continue
		}

		_, parseErr := yapar.ParseYalpContent(string(content))

		if test.shouldPass {
			if parseErr == nil {
				fmt.Printf("✅ PASS: %s\n", test.name)
				passed++
			} else {
				fmt.Printf("❌ FAIL: %s\n   Expected: success\n   Got: %v\n", test.name, parseErr)
				failed++
			}
		} else {
			if parseErr != nil {
				fmt.Printf("✅ PASS: %s\n   Error (expected): %v\n", test.name, parseErr)
				passed++
			} else {
				fmt.Printf("❌ FAIL: %s\n   Expected: error\n   Got: success\n", test.name)
				failed++
			}
		}
	}

	fmt.Printf("\n=== Results ===\n")
	fmt.Printf("Passed: %d\n", passed)
	fmt.Printf("Failed: %d\n", failed)

	if failed > 0 {
		os.Exit(1)
	}
}
