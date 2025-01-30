package output

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain"
)

type AdocReportWriter struct {
}

func NewAdocReportWriter() *AdocReportWriter {
	return &AdocReportWriter{}
}

func (w *AdocReportWriter) Write(report string) error {
	return os.WriteFile("../../report.adoc", []byte(report), 0o600)
}

func GenerateAdoc(statistic *domain.LogsStatistic, path, from, to string) string {
	var adoc string

	if to == "" {
		to = "-"
	}

	if from == "" {
		from = "-"
	}

	adoc += "#### Общая информация\n\n| Метрика | Значение |\n\n| Файл(-ы) | "
	if strings.Contains(path, "https") {
		adoc += path + " |\n\n"
	} else {
		matches, err := filepath.Glob(path)
		if err != nil {
			return ""
		}

		for _, match := range matches {
			adoc += match + " "
		}

		adoc += "|\n\n"
	}

	requests := strconv.Itoa(statistic.RequestsCount)
	adoc += "| Начальная дата | " + from + " |\n\n" +
		"| Конечная дата | " + to + " |\n\n| Количество запросов | " + requests + " |\n\n| Средний размер ответа | " +
		strconv.Itoa(statistic.AverageAnswerSize) + "b |\n\n| 95p размера ответа | " +
		strconv.Itoa(statistic.AnswerSize95p) + "b |\n\n#### Запрашиваемые ресурсы\n\n| Ресурс | Количество |\n\n"

	for resource, count := range statistic.ResourceCount {
		adoc += "| `" + resource + "` | " + strconv.Itoa(count) + " |\n\n"
	}

	adoc += "\n\n#### Коды ответа\n\n| Код | Количество |\n\n"

	for status, count := range statistic.StatusCount {
		adoc += "| " + strconv.Itoa(status) + " | " + strconv.Itoa(count) + " |\n\n"
	}

	return adoc
}
