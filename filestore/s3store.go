package filestore

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type s3store struct {
	config *S3Settings
}

func (ss s3store) save(username string, filename string, data []byte, contentType string) error {
	svc, err := ss.configureService()
	if err != nil {
		return err
	}

	params := &s3.PutObjectInput{
		Bucket:        aws.String(ss.config.Bucket),
		Key:           aws.String("/" + username + "/" + filename),
		Body:          bytes.NewReader(data),
		ContentLength: aws.Int64(int64(len(data))),
		ContentType:   aws.String(contentType),
	}

	_, err = svc.PutObject(params)
	return err
}

func (ss s3store) retrieve(username string, filename string) ([]byte, error) {
	svc, err := ss.configureService()
	if err != nil {
		return nil, err
	}

	params := &s3.GetObjectInput{
		Bucket: aws.String(ss.config.Bucket),
		Key:    aws.String("/" + username + "/" + filename),
	}

	resp, err := svc.GetObject(params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	return buf.Bytes(), nil
}

func (ss s3store) configureService() (*s3.S3, error) {
	creds := credentials.NewStaticCredentials(ss.config.AccessKeyID, ss.config.SecretKey, "")
	_, err := creds.Get()
	if err != nil {
		return nil, err
	}

	cfg := aws.NewConfig().WithRegion("ap-southeast-1").WithCredentials(creds)
	svc := s3.New(session.New(), cfg)

	return svc, nil
}
