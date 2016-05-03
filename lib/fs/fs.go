package fs // import "github.com/zentrope/tools/lib/fs"

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Copy a file at source to target. Assumes target's parent directory
// is present.
func CopyFile(source, target string) error {

	var err error

	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	defer reader.Close()

	info, err := reader.Stat()
	if err != nil {
		return err
	}

	writer, err := os.Create(target)
	if err != nil {
		return err
	}

	defer writer.Close()

	_, err = io.Copy(writer, reader)
	writer.Chmod(info.Mode())

	return err
}

// Copy a file from `sourcePath` to `targetPath` creating the
// directory path to the `target` if not present.
func CopyFileAll(sourcePath, targetPath string) error {

	parent := targetPath

	if _, ok := IsDir(parent); !ok {
		parent = filepath.Dir(targetPath)
	}

	err := os.MkdirAll(parent, 0755)
	if err != nil {
		return err
	}
	return CopyFile(sourcePath, targetPath)
}

// Copy the file tree in `sourceDir` to `targetDir`.
func CopyDir(sourceDir, targetDir string) error {

	process := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			targetFile := strings.Replace(path, sourceDir, "", 1)
			targetFile = filepath.Join(targetDir, targetFile)

			fmt.Printf(" - copy: %s -> %s\n", path, targetFile)
			CopyFileAll(path, targetFile)
		}

		return nil
	}

	err := filepath.Walk(sourceDir, process)
	if err != nil {
		return err
	}
	return nil
}

func IsDir(path string) (error, bool) {
	f, err := os.Open(path)
	if err != nil {
		return err, false
	}

	info, err := f.Stat()
	if err != nil {
		return err, false
	}
	return nil, info.IsDir()
}
