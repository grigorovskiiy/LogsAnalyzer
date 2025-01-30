package domain_test

import (
	"testing"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestUrlLogsRead(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name     string
		given    string
		expected *domain.LogRecord
	}

	data := "[17/May/2015:08:05:34 +0000]"

	testCases := []TestCase{{
		name:     "Test logs parser",
		given:    "217.168.17.5 - - " + data + " \"GET /downloads/product_1 HTTP/1.1\" 200 490 \"-\" \"Debian APT-HTTP/1.3 (0.8.10.3)\"",
		expected: domain.NewLogRecord("17/May/2015:08:05:34", 200, 490, "/downloads/product_1"),
	}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			tt.Parallel()

			actual := domain.ParseLineToLog(testCase.given)

			assert.Equal(tt, actual, testCase.expected)
		})
	}
}

func TestFilterLogRecords(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name         string
		givenRecords []*domain.LogRecord
		from         string
		to           string
		expected     []bool
	}

	testCases := []TestCase{{
		name: "Test from to log records filter",
		givenRecords: []*domain.LogRecord{domain.NewLogRecord("17/May/2015:08:05:34", 200, 490, "/downloads/product_1"),
			domain.NewLogRecord("21/May/2015:18:05:01", 404, 336, "/downloads/product_1"),
		},
		from:     "2015-05-17",
		to:       "2015-05-21",
		expected: []bool{true, false},
	},
		{
			name: "Test from - log records filter",
			givenRecords: []*domain.LogRecord{domain.NewLogRecord("17/May/2015:08:05:34", 200, 490, "/downloads/product_1"),
				domain.NewLogRecord("21/May/2015:18:05:01", 404, 336, "/downloads/product_1"),
			},
			from:     "2015-05-10",
			to:       "",
			expected: []bool{true, true},
		},
		{
			name: "Test - to log records filter",
			givenRecords: []*domain.LogRecord{domain.NewLogRecord("17/May/2015:08:05:34", 200, 490, "/downloads/product_1"),
				domain.NewLogRecord("21/May/2015:18:05:01", 404, 336, "/downloads/product_1"),
			},
			from:     "",
			to:       "2015-05-18",
			expected: []bool{true, false},
		},
		{
			name: "Test - - log records filter",
			givenRecords: []*domain.LogRecord{domain.NewLogRecord("17/May/2015:08:05:34", 200, 490, "/downloads/product_1"),
				domain.NewLogRecord("21/May/2015:18:05:01", 404, 336, "/downloads/product_1"),
			},
			from:     "",
			to:       "",
			expected: []bool{true, true},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			tt.Parallel()

			actual := make([]bool, len(testCase.givenRecords))

			for ind, record := range testCase.givenRecords {
				fromTime, _ := time.Parse("2006-01-02", testCase.from)
				toTime, _ := time.Parse("2006-01-02", testCase.to)
				actual[ind] = domain.FilterLogRecords(record, fromTime, toTime)
			}

			assert.Equal(tt, testCase.expected, actual)
		})
	}
}

func TestCalculateLogsStatistic(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name         string
		givenRecords []*domain.LogRecord
		expected     *domain.LogsStatistic
	}

	resources := map[string]int{"/downloads/product_1": 2, "/downloads/product_3": 1}
	testCases := []TestCase{{
		name: "Test calculate logs statistic",
		givenRecords: []*domain.LogRecord{
			domain.NewLogRecord("17/May/2015:08:05:32", 304, 0, "/downloads/product_1"),
			domain.NewLogRecord("30/May/2015:00:05:59", 404, 328, "/downloads/product_1"),
			domain.NewLogRecord("30/May/2015:05:05:06", 200, 1768, "/downloads/product_3"),
		},
		expected: domain.NewLogsStatistic(3, resources, map[int]int{304: 1, 404: 1, 200: 1}, 328, 1623),
	}}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(tt *testing.T) {
			tt.Parallel()

			logsStatistic := domain.NewLogsStatistic(0, make(map[string]int), make(map[int]int), 0, 0)
			answersSize := make([]int, 0)

			for _, record := range testCase.givenRecords {
				domain.CalculateLogsStatistic(logsStatistic, record, &answersSize)
			}

			domain.CalculateAverageStatistic(logsStatistic, answersSize)

			assert.Equal(tt, logsStatistic, testCase.expected)
		})
	}
}
