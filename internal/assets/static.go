package assets

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func CopyPublicAssets(watchMode bool) {
	srcDir := "public"
	destDir := "dist"

	err := copyDir(srcDir, destDir)
	if err != nil {
		log.Fatal(err)
	}

	if watchMode {
		err = watchDir(srcDir, destDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func copyFile(src, dst string) error {
	// Open the source file
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Create the destination file
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source '%s' is not a directory", src)
	}

	err = os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	// Iterate over the entries in the source directory
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the source directory itself
		if d.IsDir() && path == src {
			return nil
		}

		srcPath := path
		dstPath := filepath.Join(dst, d.Name())

		if d.IsDir() {
			return copyDir(srcPath, dstPath)
		} else {
			return copyFile(srcPath, dstPath)
		}
	})

	return err
}

func watchDir(srcDir string, destDir string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(srcDir)
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
					copyDir(srcDir, destDir)
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
