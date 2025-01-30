package application

import (
	"io"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/input"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/output"
)

type LogsReader interface {
	Read() (string, error)
	Close() error
}

type ReportWriter interface {
	Write(report string) error
}

func CreateReader(path string) (LogsReader, error) {
	if strings.Contains(path, "https") {
		urlReader, err := output.NewURLFileReader(&path)
		if err != nil {
			return nil, err
		}

		return urlReader, nil
	}

	localReader, err := output.NewLocalLogsReader(path)
	if err != nil {
		return nil, err
	}

	return localReader, nil
}

func StartApp(path, format, from, to string) error {
	reader, err := CreateReader(path)
	if err != nil {
		return err
	}

	defer func(reader LogsReader) {
		err := reader.Close()
		if err != nil {
			return
		}
	}(reader)

	logsStatistic := domain.NewLogsStatistic(0, make(map[string]int), make(map[int]int), 0, 0)
	answersSize := make([]int, 0)

	fromTime, toTime, err := input.ParseTime(from, to)

	if err != nil {
		return err
	}

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}

		logRecord := domain.ParseLineToLog(line)

		if domain.FilterLogRecords(logRecord, fromTime, toTime) {
			domain.CalculateLogsStatistic(logsStatistic, logRecord, &answersSize)
		}
	}

	if len(answersSize) > 0 {
		domain.CalculateAverageStatistic(logsStatistic, answersSize)
	}

	var writer ReportWriter

	var report string

	switch format {
	case "markdown":
		writer = output.NewMarkdownReportWriter()
		report = output.GenerateMarkdown(logsStatistic, path, from, to)

	case "adoc":
		writer = output.NewAdocReportWriter()
		report = output.GenerateAdoc(logsStatistic, path, from, to)

	default:
		writer = output.NewConsoleReportWriter()
		report = output.GenerateMarkdown(logsStatistic, path, from, to)
	}

	err = writer.Write(report)
	if err != nil {
		return err
	}

	return nil
}
