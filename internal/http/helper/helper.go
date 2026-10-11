package helper

import (
	"net/http"
	"strconv"
)

func QueryInt(r *http.Request, key string, defaultValue int) (int, error) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return defaultValue, nil
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return defaultValue, err
	}
	return i, nil
}

func QueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
