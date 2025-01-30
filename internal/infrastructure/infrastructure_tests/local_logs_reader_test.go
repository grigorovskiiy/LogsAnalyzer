package infrastructure_test

import (
	"io"
	"os"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/output"
	"github.com/stretchr/testify/assert"
)

func Read(expected []string, reader *output.LocalLogsReader) ([]string, error) {
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		expected = append(expected, line)
	}

	return expected, nil
}

func TestLocalLogsRead(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name         string
		givenFiles   []string
		givenContent []string
		expected     []string
	}

	testCases := []TestCase{{
		name:         "Test dir/* pattern",
		givenFiles:   []string{"logs1.txt", "logs2.txt"},
		givenContent: []string{"logs1-content", "logs2-content"},
		expected: []string{
			"logs1-content", "logs2-content",
		}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			tt.Parallel()

			tempDir, err := os.MkdirTemp("", "temp_dir")

			if err != nil {
				t.Fatal(err)
			}

			for ind, fileName := range testCase.givenFiles {
				file, err := os.CreateTemp(tempDir, fileName)
				if err != nil {
					t.Fatal(err)
				}

				_, err = file.WriteString(testCase.givenContent[ind])

				if err != nil {
					t.Fatal(err)
				}
			}

			var reader *output.LocalLogsReader

			reader, err = output.NewLocalLogsReader(tempDir + "/*")
			if err != nil {
				t.Fatal(err)
			}

			mas := make([]string, 0)
			expected, err := Read(mas, reader)

			if err != nil {
				t.Fatal(err)
			}

			err = reader.Close()

			if err != nil {
				return
			}

			err = os.RemoveAll(tempDir)

			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, expected, testCase.expected)
		})
	}
}
