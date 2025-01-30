package output

import (
	"io"
	"os"
)

type ConsoleReportWriter struct{}

func NewConsoleReportWriter() *ConsoleReportWriter { return &ConsoleReportWriter{} }

func (w *ConsoleReportWriter) Write(report string) error {
	_, err := io.Writer(os.Stdout).Write([]byte(report))
	return err
}
