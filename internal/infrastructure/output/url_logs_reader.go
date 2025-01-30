package output

import (
	"bufio"
	"io"
	"net/http"

	"github.com/es-debug/backend-academy-2024-go-template/internal/errors"
)

type URLFileReader struct {
	responseBody io.ReadCloser
	scanner      *bufio.Scanner
}

func NewURLFileReader(url *string) (*URLFileReader, error) {
	res, err := http.Get(*url)
	if res.StatusCode != 200 {
		return nil, errors.URLError{}
	}

	if err != nil {
		_ = res.Body.Close()
		return nil, errors.URLError{}
	}

	scanner := bufio.NewScanner(res.Body)
	reader := &URLFileReader{res.Body, scanner}

	return reader, err
}

func (r *URLFileReader) Read() (string, error) {
	if r.scanner.Scan() {
		return r.scanner.Text(), nil
	}

	return "", io.EOF
}

func (r *URLFileReader) Close() error {
	return r.responseBody.Close()
}
