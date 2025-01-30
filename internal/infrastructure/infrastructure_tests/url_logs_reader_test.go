package infrastructure_test

import (
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/output"
	"github.com/stretchr/testify/assert"
)

func TestUrlLogsRead(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name     string
		givenURL string
		expected []string
	}

	data1 := "[17/May/2015:08:05:32 +0000]"
	data2 := "[17/May/2015:08:05:23 +0000]"
	addr := "93.180.71.3 "
	res := "/downloads/product_1"

	testCases := []TestCase{{
		name:     "Test url",
		givenURL: "https://raw.githubusercontent.com/elastic/examples/master/Common%20Data%20Formats/nginx_logs/nginx_logs",
		expected: []string{addr + "- - " + data1 + " \"GET " + res + " HTTP/1.1\" 304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.21)\"",
			"93.180.71.3 - - " + data2 + " \"GET /downloads/product_1 HTTP/1.1\" 304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.21)\"",
		}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			tt.Parallel()

			reader, err := output.NewURLFileReader(&testCase.givenURL)

			if err != nil {
				tt.Fatal(err)
			}

			actual := make([]string, 0)

			for i := 0; i < len(testCase.expected); i++ {
				line, err := reader.Read()
				if err != nil {
					tt.Fatal(err)
				}

				actual = append(actual, line)
			}

			assert.Equal(tt, testCase.expected, actual)
		})
	}
}
