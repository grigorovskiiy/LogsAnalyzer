package main

import (
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/application"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/input"
)

func main() {
	path, format, from, to, err := input.GetParameters()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = input.CheckParameters(path, format)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = application.StartApp(path, format, from, to)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Отчет сформирован в соответствующем формату файле, либо в консоли")
}
