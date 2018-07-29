package s3

import (
	"bytes"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3 holds simple storage service session
type S3 struct {
	*session.Session
}

// New returns new instance of S3 session
func New(region string) *S3 {
	ses := session.Must(session.NewSession(&aws.Config{Region: aws.String(region)}))
	return &S3{ses}
}

// Get retrieves object by key from given S3 bucket
func (s *S3) Get(bucket, key string) (res []byte, err error) {
	obj, err := s3.New(s).GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return
	}

	defer obj.Body.Close()
	res, err = ioutil.ReadAll(obj.Body)
	if err != nil {
		return
	}

	return
}

// Put creates or updates given object by key in S3 bucket
func (s *S3) Put(bucket, key string, body []byte) (err error) {
	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ACL:    aws.String("private"),
		Body:   bytes.NewReader(body),
	})
	if err != nil {
		return
	}

	return nil
}
