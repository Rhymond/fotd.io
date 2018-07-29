package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mgutz/ansi"
	"github.com/rhymond/fotd/s3"
)

type s3Getter func(bucket, key string) (res []byte, err error)

// Handler holds injected services
type Handler struct {
	S3Getter s3Getter
}

// New creates new Handler instance
func New(s3g s3Getter) *Handler {
	return &Handler{s3g}
}

// handle returns random fact from S3 bucket
func (h *Handler) handle() (events.APIGatewayProxyResponse, error) {
	res, err := h.S3Getter(os.Getenv("S3_BUCKET"), "facts")
	if err != nil {
		return h.error(err.Error())
	}

	facts := make([]string, 0)
	err = json.Unmarshal(res, &facts)
	if err != nil {
		return h.error(err.Error())
	}

	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(facts)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body: fmt.Sprintf(
			"%s\n%s\n%s",
			ansi.Color("Did you know?", "black:white"),
			ansi.Color(facts[n], "white+bh"),
			ansi.Color("source: FACTSlides.com", "yellow"),
		),
		Headers: map[string]string{"content-type": "text/html"},
	}, nil
}

// errors returns 400 API Gateway response error with given errors message
func (h *Handler) error(err string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       err,
		Headers:    map[string]string{"content-type": "text/html"},
	}, nil
}

func main() {
	s3Ses := s3.New(os.Getenv("S3_REGION"))
	h := New(s3Ses.Get)

	lambda.Start(h.handle)
}
