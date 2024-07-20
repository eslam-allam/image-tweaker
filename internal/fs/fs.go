package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Jeffail/tunny"
)

type pathtype int

const (
	File pathtype = iota
	Directory
)

func Exists(path string) (bool, pathtype, error) {
	info, err := os.Stat(path)
	if err == nil {
		var t pathtype
		if info.IsDir() {
			t = Directory
		} else {
			t = File
		}
		return true, t, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, 0, nil
	}

	return false, 0, err
}

type (
	FileTransformer  func(fileSrc string, fileDst string) error
	threadSubmission struct {
		srcDir string
		dstDir string
	}
)

func shadowWalkDirParallelFile(threads uint, srcDir, dstDir string, fn FileTransformer) error {
	if !filepath.IsAbs(srcDir) || !filepath.IsAbs(dstDir) {
		return errors.New("src and dst directories must be absolute")
	}

	exists, typ, err := Exists(dstDir)
	if err != nil {
		return err
	}
	if exists && typ == File {
		return errors.New("dst directory is a file")
	} else if !exists {
		err := os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	pool := tunny.NewFunc(int(threads), func(payload interface{}) interface{} {
		paths, ok := payload.(threadSubmission)
		if !ok {
			return errors.New("invalid payload for file transformer")
		}
		return fn(paths.srcDir, paths.dstDir)
	})
	defer pool.Close()

	err = filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("failed to process '%s'", path)
			return fs.SkipDir
		}
		dstComponent, _ := strings.CutPrefix(path, srcDir)
		dstPath := filepath.Join(dstDir, dstComponent)
		if d.IsDir() {
			if path != srcDir {
				err = os.MkdirAll(dstPath, os.ModePerm)
				if err != nil {
					fmt.Printf("failed to process '%s'", path)
					return fs.SkipDir
				}
			}
			return nil
		}
		result := pool.Process(threadSubmission{path, dstPath})
		if result == nil {
			return nil
		} else {
			fmt.Printf("failed to transform '%s': %s\n", path, result.(error))
		}
		return nil
	})
	return err
}

// parallelises on files only
func ApplyDirectoryParallel(threads uint, srcDir, dstDir string, fn FileTransformer) error {
	exists, typ, err := Exists(srcDir)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("directory '%s' does not exist", srcDir)
	}
	if typ != Directory {
		return fmt.Errorf("path '%s' is not a directory", dstDir)
	}
	return shadowWalkDirParallelFile(threads, srcDir, dstDir, fn)
}
