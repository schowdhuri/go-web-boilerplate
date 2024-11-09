package assets

import (
	"fmt"
	"log"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func FooBar() {
	fmt.Println("FooBar")
}

func BuildJs(watchMode bool) {
	// Set up build options
	options := esbuild.BuildOptions{
		EntryPoints:       []string{"assets/js/*.js"},
		Outdir:            "dist/js",
		Bundle:            true,
		Write:             true,
		MinifyWhitespace:  !watchMode,
		MinifyIdentifiers: !watchMode,
		MinifySyntax:      !watchMode,
	}

	// Build the JavaScript files
	result := esbuild.Build(options)
	if len(result.Errors) != 0 {
		for _, err := range result.Errors {
			log.Fatalf("Build failed: %v", err)
		}
	}

	// Watch for changes in development mode
	if watchMode {
		ctx, err := esbuild.Context(options)
		if err != nil {
			log.Fatalf("Build failed: %v", err)
		}

		err2 := ctx.Watch(esbuild.WatchOptions{})
		if err2 != nil {
			log.Fatalf("Unable to watch files: %v", err2)
		}
		fmt.Println("Watching JS files...")
		make(chan struct{}) <- struct{}{}
	} else {
		fmt.Println("JS bundles built")
	}
}
