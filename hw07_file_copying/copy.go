package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	srcFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	fileSize := fileInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || offset+limit > fileSize {
		limit = fileSize - offset
	}

	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek in source file: %w", err)
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	return copyWithProgress(srcFile, dstFile, limit)
}

func copyWithProgress(src io.Reader, dst io.Writer, limit int64) error {
	buffer := make([]byte, 32*1024)
	var copied int64

	for copied < limit {
		toRead := limit - copied
		if toRead > int64(len(buffer)) {
			toRead = int64(len(buffer))
		}

		n, err := src.Read(buffer[:toRead])
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("read error: %w", err)
		}

		if n == 0 {
			break
		}

		if _, writeErr := dst.Write(buffer[:n]); writeErr != nil {
			return fmt.Errorf("write error: %w", writeErr)
		}

		copied += int64(n)

		if copied != limit {
			fmt.Printf("\rProgress: %.2f%%", float64(copied)/float64(limit)*100)
		}
	}

	return nil
}
