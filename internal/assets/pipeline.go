package assets

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type AssetPipeline struct {
	IsDevelopment bool
	cache         map[string]string
	mutex         sync.RWMutex
}

func NewAssetPipeline(isDev bool) *AssetPipeline {
	return &AssetPipeline{
		IsDevelopment: isDev,
		cache:         make(map[string]string),
	}
}

func (ap *AssetPipeline) GetAssetURL(path string) string {
	path = filepath.Join("/", "dist", path)

	if ap.IsDevelopment {
		return path
	}

	ap.mutex.RLock()
	hash, exists := ap.cache[path]
	ap.mutex.RUnlock()

	if exists {
		return path + "?v=" + hash
	}

	hash, _ = calculateFileHash(path)

	ap.mutex.Lock()
	ap.cache[path] = hash
	ap.mutex.Unlock()

	return path + "?v=" + hash
}

func calculateFileHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", err
	}
	defer f.Close()

	// Create a SHA-256 hash
	h := sha256.New()

	// Copy the file content to the hash
	if _, err := io.Copy(h, f); err != nil {
		fmt.Println("Error hashing file:", err)
		return "", err
	}

	// Get the hash sum
	hash := h.Sum(nil)

	// hash in hexadecimal format
	return fmt.Sprintf("%x", hash), nil
}
