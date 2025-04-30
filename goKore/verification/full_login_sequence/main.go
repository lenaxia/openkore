package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting Full Login Sequence Test")

	// Run the full login sequence test
	test := NewLoginSequenceTest()
	if test == nil {
		fmt.Println("Failed to create test")
		os.Exit(1)
	}

	success := test.RunTest()
	if success {
		fmt.Println("Full login sequence test passed!")
		os.Exit(0)
	} else {
		fmt.Println("Full login sequence test failed!")
		os.Exit(1)
	}
}
