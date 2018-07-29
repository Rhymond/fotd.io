package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin/json"
	"github.com/rhymond/fotd/s3"
)

// Injection types.
type httpGetter func(url string) (resp *http.Response, err error)
type s3Putter func(bucket, key string, body []byte) error

// Handler holds injected services
type Handler struct {
	HTTPGetter httpGetter
	S3Putter   s3Putter
}

// New creates new Handler instance
func New(hg httpGetter, s3p s3Putter) *Handler {
	return &Handler{hg, s3p}
}

// handle fetches facts from given source and stores to S3 bucket
func (h *Handler) handle() error {
	res, err := h.HTTPGetter("https://www.factslides.com")
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	var facts []string
	doc.Find(".factTools").Each(func(i int, s *goquery.Selection) {
		text := s.ParentFiltered("div").Text()
		if text != "" {
			facts = append(facts, strings.Join(strings.Fields(text), " "))
		}
	})

	factsJSON, err := json.Marshal(facts)
	if err != nil {
		return err
	}

	err = h.S3Putter(os.Getenv("S3_BUCKET"), "facts", factsJSON)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	s3Ses := s3.New(os.Getenv("S3_REGION"))
	h := New(http.Get, s3Ses.Put)

	lambda.Start(h.handle)
}
