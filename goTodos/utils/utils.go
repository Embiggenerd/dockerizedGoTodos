package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// RandHex simply creates random hex string
// To use to query session data
func RandHex(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)

	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func fileHash(filePath string) (string, error) {
	var fileHash string
	file, err := os.Open(filePath)
	defer file.Close()

	hash := md5.New()

	_, err = io.Copy(hash, file)

	hashInBytes := hash.Sum(nil)[:16]

	fileHash = hex.EncodeToString(hashInBytes)

	return fileHash, err
}

func copyFileToPublic(path string) error {
	source, err := os.Open(path)

	defer source.Close()
	_, filename := filepath.Split(path)
	newPath := "./public/" + filename

	destination, err := os.Create(newPath)
	defer destination.Close()

	_, err = io.Copy(destination, source)

	return err
}

func removeStaleFiles(prefix string) error {
	separated := strings.Split(prefix, ".")
	err := filepath.Walk("./public", func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		ok := strings.HasPrefix(file.Name(), separated[0])
		if ok {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// BustaCache creates a css file with its hash included in the name
func BustaCache(filename string) (string, error) {
	var filenamePlusHash string
	// Create 'public' filename if it doesn't exist
	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.Mkdir("public", 0755)
	}

	err := removeStaleFiles(filename)
	if err != nil {
		fmt.Println("Busta Cache Error", err)
		return filenamePlusHash, err
	}

	// Loop through the files in assets
	err = filepath.Walk("./assets", func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Look for the file with our filename
		ok := strings.HasPrefix(file.Name(), filename)

		if ok {
			// copy file from assets to publix
			err := copyFileToPublic(path)
			if err != nil {
				fmt.Println(err)
				return err
			}
			// Get the path to the file
			newPath := "./public/" + file.Name()
			// Generate a hash
			hash, err := fileHash(newPath)

			if err != nil {
				fmt.Println(err)
				return err
			}

			separated := strings.Split(file.Name(), ".")
			// Write the hash
			var b strings.Builder
			b.WriteString("./")
			b.WriteString(filepath.Dir(newPath))
			b.WriteString("/")
			b.WriteString(separated[0])
			b.WriteString(".")
			b.WriteString(hash)
			b.WriteString(".")
			b.WriteString(separated[1])

			_, filenamePlusHash = filepath.Split(b.String())

			//
			err = os.Rename(newPath, b.String())
			if err != nil {
				return err
			}
		}
		return err
	})

	return filenamePlusHash, err
}
