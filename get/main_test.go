package main

import (
	"errors"
	"strings"
	"testing"
)

type testCase struct {
	expectedBody string
	S3Getter     s3Getter
}

func TestOutput(t *testing.T) {
	run(t, testCase{
		expectedBody: "fact 1",
		S3Getter: func(bucket, key string) (res []byte, err error) {
			return []byte(`["fact 1"]`), nil
		},
	})
}

func TestS3Error(t *testing.T) {
	run(t, testCase{
		expectedBody: "error",
		S3Getter: func(bucket, key string) (res []byte, err error) {
			return nil, errors.New("error")
		},
	})
}

func TestS3WrongJSONFormat(t *testing.T) {
	run(t, testCase{
		expectedBody: "invalid character",
		S3Getter: func(bucket, key string) (res []byte, err error) {
			return []byte("random-string"), nil
		},
	})
}

func run(t *testing.T, test testCase) {
	h := New(
		test.S3Getter,
	)

	resp, err := h.handle()

	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}

	if !strings.Contains(resp.Body, test.expectedBody) {
		t.Errorf("Expected body %s to contain %s", resp.Body, test.expectedBody)
	}
}
