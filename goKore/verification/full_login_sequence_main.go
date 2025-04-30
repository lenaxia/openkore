package main

import (
	"fmt"
	"os"

	"github.com/lenaxia/goKore/verification"
)

func main() {
	fmt.Println("Starting Full Login Sequence Test")

	// Run the full login sequence test
	success := verification.RunFullLoginSequenceTest()

	if success {
		fmt.Println("Full login sequence test passed!")
		os.Exit(0)
	} else {
		fmt.Println("Full login sequence test failed!")
		os.Exit(1)
	}
}
