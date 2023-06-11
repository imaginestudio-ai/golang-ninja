package main

import (
	"fmt"
	"github.com/ImagineDevOps/Network-Automation-with-Go/ch02/ping"
)

func main() {
	s := ping.Send()
	fmt.Println(s)
}
