package main

import (
	"fmt"

	"github.com/ImagineDevOps/Chapter7/Examples/Performance/pkg"
)

func main() {
	fmt.Println(pkg.RecursiveFact(10))
	//debug.SetMaxStack(100)
	stackOverflowExample(0)
}
