package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type testCase struct {
	expectedErr string
	s3Getter    s3Putter
	httpGetter  httpGetter
}

func TestHTTPHandlerError(t *testing.T) {
	run(t, testCase{
		expectedErr: "error-http-handler",
		httpGetter: func(url string) (resp *http.Response, err error) {
			return nil, errors.New("error-http-handler")
		},
	})
}

func TestHTTPGetterStatus400(t *testing.T) {
	run(t, testCase{
		expectedErr: "400",
		httpGetter: func(url string) (resp *http.Response, err error) {
			return &http.Response{
				StatusCode: 400,
				Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			}, nil
		},
	})
}

func run(t *testing.T, test testCase) {
	h := New(
		test.httpGetter,
		test.s3Getter,
	)

	err := h.handle()

	if !strings.Contains(err.Error(), test.expectedErr) {
		t.Errorf("Expected body %s to contain %s", err.Error(), test.expectedErr)
	}
}
