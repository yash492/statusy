package postgres

import (
	"bytes"
	"text/template"
)

type transformFunc[T any] func(params T) map[string]any

func ParseQueryWithValues[T any](
	templString string,
	values T,
	transform transformFunc[T],
) (string, error) {
	var buffer bytes.Buffer
	queryValues := transform(values)
	tmpl, err := template.New("query").Parse(templString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buffer, queryValues)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func ParseQuery(
	templString string,
) (string, error) {
	var buffer bytes.Buffer
	tmpl, err := template.New("query").Parse(templString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buffer, nil)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
