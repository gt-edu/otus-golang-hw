package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	ErrOffsetAndLimitShouldBePositive = errors.New("offset and limit should be positive")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFileSize, err := validateParams(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer closeResource(fromFile)

	if offset > 0 {
		_, err := fromFile.Seek(offset, io.SeekStart)
		if err != nil {
			return err
		}
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer closeResource(toFile)

	copyUntilEnd := limit == 0 || (offset+limit) > fromFileSize
	bar, barReader := setupProgressBar(fromFile, fromFileSize, copyUntilEnd, offset, limit)
	if copyUntilEnd {
		_, err = io.Copy(toFile, barReader)
	} else {
		_, err = io.CopyN(toFile, barReader, limit)
	}

	if err != nil && errors.Is(err, io.EOF) {
		return err
	}

	bar.Finish()
	return nil
}

func closeResource(closer io.Closer) {
	if err := closer.Close(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error during closing file: %v\n", err)
	}
}

func setupProgressBar(
	fromFile *os.File, fromFileSize int64, copyUntilEnd bool, offset int64, limit int64,
) (*pb.ProgressBar, *pb.Reader) {
	pbCount := fromFileSize
	if copyUntilEnd {
		pbCount -= offset
	} else {
		pbCount = limit
	}
	bar := pb.Full.Start64(pbCount)
	barReader := bar.NewProxyReader(fromFile)
	return bar, barReader
}

func validateParams(fromPath, toPath string, offset, limit int64) (int64, error) {
	if offset < 0 || limit < 0 {
		return 0, ErrOffsetAndLimitShouldBePositive
	}

	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return 0, err
	}

	if !fromFileInfo.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
	}

	toFileInfo, err := os.Stat(toPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return 0, err
	}

	if toFileInfo != nil && !toFileInfo.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
	}

	fromFileSize := fromFileInfo.Size()
	if offset > fromFileSize {
		return 0, ErrOffsetExceedsFileSize
	}

	return fromFileSize, nil
}
