package input

import (
	"slices"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func GetParameters() (path, format, from, to string, err error) {
	pflag.String("path", "", "Path parameter")
	pflag.String("format", "", "Format parameter")
	pflag.String("from", "", "From parameter")
	pflag.String("to", "", "To parameter")
	pflag.Parse()

	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return "", "", "", "", err
	}

	return viper.GetString("path"), viper.GetString("format"), viper.GetString("from"), viper.GetString("to"), nil
}

func CheckParameters(path, format string) error {
	err := CheckPath(path)
	if err != nil {
		return err
	}

	formats := []string{"markdown", "adoc", ""}
	err = CheckFormat(format, formats)

	if err != nil {
		return err
	}

	return nil
}
func CheckPath(arg string) error {
	if arg == "" {
		return errors.PathError{}
	}

	return nil
}

func CheckFormat(arg string, formats []string) error {
	if !slices.Contains(formats, arg) {
		return errors.FormatError{}
	}

	return nil
}

func ParseTime(from, to string) (fromTime, toTime time.Time, err error) {
	if from != "" {
		fromTime, err := time.Parse("2006-01-02", from)
		if err != nil {
			return fromTime, toTime, errors.FromError{}
		}
	}

	if to != "" {
		toTime, err := time.Parse("2006-01-02", to)
		if err != nil {
			return fromTime, toTime, errors.ToError{}
		}
	}

	return fromTime, toTime, nil
}
