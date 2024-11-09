package assets

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func FooBar2() {
	fmt.Println("FooBar")
}

func BuildCss(watchMode bool) {
	inputDir := "assets/css"
	outputDir := "dist/css"

	err := build(inputDir, outputDir)
	if err != nil {
		log.Fatal(err)
	}

	if watchMode {
		err = watch(inputDir, outputDir)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Println("CSS built")
	}
}

func build(inputDir, outputDir string) error {
	files, err := os.ReadDir(inputDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		inputFile := filepath.Join(inputDir, file.Name())
		outputFile := filepath.Join(outputDir, file.Name())

		cmd := exec.Command("./node_modules/.bin/tailwindcss", "-i", inputFile, "-o", outputFile)
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func watch(inputDir string, outputDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(inputDir)
	if err != nil {
		return err
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("Event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("File changed, rebuilding CSS...")
					build(inputDir, outputDir)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	<-done
	return nil
}
