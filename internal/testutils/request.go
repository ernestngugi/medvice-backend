package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
)

func DoRequest(r http.Handler, method, path string, form any) (*httptest.ResponseRecorder, error) {

	w := httptest.NewRecorder()

	var body io.Reader

	if form != nil {
		b, err := json.Marshal(form)
		if err != nil {
			return w, err
		}

		body = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return w, err
	}

	r.ServeHTTP(w, req)
	return w, nil
}
