package main

import (
	"os"
	"time"

	clockface "github.com/tomaaron/learn-go-with-tests/16-maths"
)

func main() {
	t := time.Now()
	clockface.SVGWriter(os.Stdout, t)
}
