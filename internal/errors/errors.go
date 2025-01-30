package errors

type FormatError struct{}

func (err FormatError) Error() string {
	return "Такой формат не поддерживается."
}

type PathError struct{}

func (err PathError) Error() string {
	return "Путь не был введен."
}

type FromError struct{}

func (err FromError) Error() string {
	return "Начальная дата введена в неподдерживаемом формате. Введите формате YYYY-MM-DD."
}

type ToError struct{}

func (err ToError) Error() string {
	return "Конечная дата введена в неподдерживаемом формате. Введите формате YYYY-MM-DD."
}

type URLError struct{}

func (err URLError) Error() string {
	return "Ошибка при получении данных по Url"
}

type PatternError struct{}

func (err PatternError) Error() string {
	return "Ни один файл не удовлетворяет введенному шаблону"
}
