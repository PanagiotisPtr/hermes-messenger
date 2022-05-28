package handler

import (
	"fmt"
	"net/http"
	"net/url"
)

type httpHandler func(http.ResponseWriter, *http.Request)

func getFormValue(formValues url.Values, key string) (string, error) {
	values, ok := formValues[key]
	if !ok || len(values) != 1 {
		return "", fmt.Errorf("Could not find a single value for key: %s", key)
	}

	return values[0], nil
}
