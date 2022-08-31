package main

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	ErrOffsetAndLimitShouldBePositive = errors.New("offset and limit should be positive")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromFile, fromFileSize, err := validateParams(fromPath, toPath, offset, limit)
	if err != nil {
		return err
	}

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

func validateParams(fromPath, toPath string, offset, limit int64) (*os.File, int64, error) {
	if offset < 0 || limit < 0 {
		return nil, 0, ErrOffsetAndLimitShouldBePositive
	}

	if strings.Index(fromPath, "/dev/") == 0 || strings.Index(toPath, "/dev/") == 0 {
		return nil, 0, ErrUnsupportedFile
	}

	fromFile, err := os.Open(fromPath)
	if err != nil {
		return nil, 0, err
	}

	fromFileInfo, err := fromFile.Stat()
	if err != nil {
		return nil, 0, err
	}
	fromFileSize := fromFileInfo.Size()
	if offset > fromFileSize {
		return nil, 0, ErrOffsetExceedsFileSize
	}

	return fromFile, fromFileSize, nil
}
