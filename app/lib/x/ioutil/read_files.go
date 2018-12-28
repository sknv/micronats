package ioutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const (
	errWalkFiles = "failed to walk files"
)

func ReadFilesWalk(root string, ext string) (string, error) {
	buf := new(bytes.Buffer)
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ext) {
			return errors.WithMessage(err, errWalkFiles) // returns error if exists or nil otherwise
		}

		// read the file
		file, err := os.Open(path)
		if err != nil {
			return errors.WithMessage(err, errWalkFiles)
		}
		defer file.Close()

		// copy the file content to the buffer
		_, err = io.Copy(buf, file)
		if err != nil {
			return errors.WithMessage(err, errWalkFiles)
		}
		return nil
	}); err != nil {
		return "", errors.WithMessagef(err, "failed to read files recursively in %s", root)
	}
	return buf.String(), nil
}
