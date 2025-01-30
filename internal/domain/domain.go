package domain

import (
	"sort"
	"strconv"
	"strings"
	"time"
)

type LogRecord struct {
	localTime string
	status    int
	bytesSent int
	resource  string
}

type LogsStatistic struct {
	RequestsCount     int
	ResourceCount     map[string]int
	StatusCount       map[int]int
	AverageAnswerSize int
	AnswerSize95p     int
}

func NewLogRecord(localTime string, status, bytesSent int, resource string) *LogRecord {
	return &LogRecord{localTime, status, bytesSent, resource}
}

func NewLogsStatistic(reqCount int, resCount map[string]int, statusCount map[int]int, averageAnswerSize, answerSize95p int) *LogsStatistic {
	return &LogsStatistic{reqCount, resCount, statusCount, averageAnswerSize, answerSize95p}
}

func ParseLineToLog(line string) *LogRecord {
	splitLine := strings.Split(line, " ")

	var localTime string

	status, _ := strconv.Atoi(splitLine[8])
	bytesSent, _ := strconv.Atoi(splitLine[9])
	resource := splitLine[6]

	for _, el := range splitLine {
		if strings.Contains(el, "[") {
			for _, ch := range el {
				if ch != '[' && ch != ']' {
					localTime += string(ch)
				}
			}
		}
	}

	return NewLogRecord(localTime, status, bytesSent, resource)
}

func FilterLogRecords(record *LogRecord, fromTime, toTime time.Time) bool {
	if toTime.Before(fromTime) && !toTime.IsZero() {
		return false
	}

	logTime, _ := time.Parse("02/Jan/2006:15:04:05", record.localTime)

	if logTime.Equal(fromTime) || logTime.Equal(toTime) {
		return true
	}

	if (logTime.After(fromTime) && logTime.Before(toTime)) || (toTime.IsZero() && logTime.After(fromTime)) {
		return true
	}

	return false
}

func CalculateLogsStatistic(logsStatistic *LogsStatistic, record *LogRecord, answersSize *[]int) {
	logsStatistic.RequestsCount++
	if _, ok := logsStatistic.ResourceCount[record.resource]; !ok {
		logsStatistic.ResourceCount[record.resource] = 1
	} else {
		logsStatistic.ResourceCount[record.resource]++
	}

	if _, ok := logsStatistic.StatusCount[record.status]; !ok {
		logsStatistic.StatusCount[record.status] = 1
	} else {
		logsStatistic.StatusCount[record.status]++
	}

	*answersSize = append(*answersSize, record.bytesSent)
}

func CalculateAverageStatistic(logsStatistic *LogsStatistic, answersSize []int) {
	sort.Ints(answersSize)

	if len(answersSize) == 0 {
		logsStatistic.AnswerSize95p = 0
	} else {
		p := 0.95 * float64(logsStatistic.RequestsCount-1)
		indexFloor := int(p)
		indexCeil := indexFloor + 1

		logsStatistic.AnswerSize95p = answersSize[indexFloor]

		if indexCeil < logsStatistic.RequestsCount {
			logsStatistic.AnswerSize95p += int((p - float64(indexFloor)) * (float64(answersSize[indexCeil]) - float64(answersSize[indexFloor])))
		}
	}

	logsStatistic.AverageAnswerSize = answersSize[logsStatistic.RequestsCount/2]
}
