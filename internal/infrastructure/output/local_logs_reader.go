package output

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/errors"
)

type LocalLogsReader struct {
	scanners []*bufio.Scanner
	files    []*os.File
}

func NewLocalLogsReader(pattern string) (*LocalLogsReader, error) {
	var scanners []*bufio.Scanner

	var files []*os.File

	if strings.Contains(pattern, "**") {
		return DirsPattern(files, scanners, pattern)
	}

	return FilesPattern(files, scanners, pattern)
}

func (r *LocalLogsReader) Read() (string, error) {
	for _, scanner := range r.scanners {
		if scanner.Scan() {
			return scanner.Text(), nil
		}

		if err := scanner.Err(); err != nil {
			return "", err
		}
	}

	return "", io.EOF
}

func (r *LocalLogsReader) Close() error {
	for _, file := range r.files {
		err := file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func FilesPattern(files []*os.File, scanners []*bufio.Scanner, pattern string) (*LocalLogsReader, error) {
	matches, err := filepath.Glob(pattern)
	if len(matches) == 0 {
		return nil, errors.PatternError{}
	}

	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		file, err := os.Open(match)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
		scanner := bufio.NewScanner(file)
		scanners = append(scanners, scanner)
	}

	return &LocalLogsReader{files: files, scanners: scanners}, nil
}

func DirsPattern(files []*os.File, scanners []*bufio.Scanner, pattern string) (*LocalLogsReader, error) {
	splitPattern := strings.Split(pattern, "/")
	root := ""

	var suffix string

	for ind := range splitPattern {
		if splitPattern[ind] == "**" {
			for ind1 := ind + 1; ind1 < len(splitPattern); ind1++ {
				suffix += "/" + splitPattern[ind1]
			}

			break
		} else if splitPattern[ind] != "" {
			root += splitPattern[ind] + "/"
		}
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type().IsRegular() && strings.HasSuffix(path, suffix) {
			f, err := os.Open(path)
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(f)
			scanners = append(scanners, scanner)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &LocalLogsReader{files: files, scanners: scanners}, nil
}
