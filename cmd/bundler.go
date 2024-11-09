package main

import (
	"fmt"
	"sync"

	"viabl.ventures/gossr/internal/assets"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	fmt.Println("Building static assets...")

	go func() {
		defer wg.Done()
		assets.BuildJs(false)
	}()

	go func() {
		defer wg.Done()
		assets.BuildCss(false)
	}()

	go func() {
		defer wg.Done()
		assets.CopyPublicAssets(false)
	}()

	wg.Wait()

	fmt.Println("Static assets built")
}
