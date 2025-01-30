package output

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type MarkdownReportWriter struct {
}

func NewMarkdownReportWriter() *MarkdownReportWriter {
	return &MarkdownReportWriter{}
}

func (w *MarkdownReportWriter) Write(report string) error {
	return os.WriteFile("../../report.md", []byte(report), 0o600)
}

func GenerateMarkdown(statistic *domain.LogsStatistic, path, from, to string) string {
	var markdown string

	if to == "" {
		to = "-"
	}

	if from == "" {
		from = "-"
	}

	markdown += "#### Общая информация\n\n| Метрика | Значение |\n|:----:|:---:|\n| Файл(-ы) | "
	if strings.Contains(path, "https") {
		markdown += path + " |\n"
	} else {
		matches, err := filepath.Glob(path)

		if err != nil {
			return ""
		}

		for _, match := range matches {
			markdown += match + " "
		}

		markdown += "|\n"
	}

	markdown += "| Начальная дата | " + from + " |\n" +
		"| Конечная дата | " + to + " |\n| Количество запросов | " + strconv.Itoa(statistic.RequestsCount) + " |\n| Средний размер ответа | " +
		strconv.Itoa(statistic.AverageAnswerSize) + "b |\n| 95p размера ответа | " +
		strconv.Itoa(statistic.AnswerSize95p) + "b |\n\n#### Запрашиваемые ресурсы\n\n| Ресурс | Количество |\n|:----:|:---:|\n"

	for resource, count := range statistic.ResourceCount {
		markdown += "| `" + resource + "` | " + strconv.Itoa(count) + " |\n"
	}

	markdown += "\n#### Коды ответа\n\n| Код | Количество |\n|:----:|:---:|\n"

	for status, count := range statistic.StatusCount {
		markdown += "| " + strconv.Itoa(status) + " | " + strconv.Itoa(count) + " |\n"
	}

	return markdown
}
