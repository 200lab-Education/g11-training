package component

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	sctx "github.com/viettranx/service-context"
	"log"
	"net/http"
	"time"
)

type UploadProvider interface {
	SaveFileUploaded(ctx context.Context, data []byte, dst string) error
}

type s3Provider struct {
	id         string
	bucketName string
	region     string
	apiKey     string
	secret     string
	domain     string
	session    *session.Session
}

func (p *s3Provider) ID() string { return p.id }

func (p *s3Provider) InitFlags() {
	flag.StringVar(
		&p.secret,
		"aws-s3-secret",
		"",
		"ASW S3 secret key",
	)

	flag.StringVar(
		&p.apiKey,
		"aws-s3-api-key",
		"",
		"ASW S3 API key",
	)

	flag.StringVar(
		&p.bucketName,
		"aws-s3-bucket",
		"",
		"ASW S3 bucket name",
	)

	flag.StringVar(
		&p.region,
		"aws-s3-region",
		"ap-southeast-1",
		"ASW S3 region name",
	)

	flag.StringVar(
		&p.domain,
		"cdn-domain",
		"",
		"ASW S3 API key",
	)
}

func NewAWSS3Provider(id string) *s3Provider {
	return &s3Provider{id: id}
}

func NewS3Provider(bucketName string, region string, apiKey string, secret string, domain string) *s3Provider {
	provider := &s3Provider{
		bucketName: bucketName,
		region:     region,
		apiKey:     apiKey,
		secret:     secret,
		domain:     domain,
	}

	s3Session, err := session.NewSession(&aws.Config{
		Region: aws.String(provider.region),
		Credentials: credentials.NewStaticCredentials(
			provider.apiKey, // Access key ID
			provider.secret, // Secret access key
			""),             // Token can be ignore
	})

	if err != nil {
		log.Fatalln(err)
	}

	provider.session = s3Session

	return provider
}

func (p *s3Provider) Activate(_ sctx.ServiceContext) error {
	s3Session, err := session.NewSession(&aws.Config{
		Region: aws.String(p.region),
		Credentials: credentials.NewStaticCredentials(
			p.apiKey, // Access key ID
			p.secret, // Secret access key
			""),      // Token can be ignore
	})

	if err != nil {
		return err
	}

	p.session = s3Session
	return nil
}
func (p *s3Provider) Stop() error {
	return nil
}

func (p *s3Provider) SaveFileUploaded(ctx context.Context, data []byte, dst string) error {
	fileBytes := bytes.NewReader(data)
	fileType := http.DetectContentType(data)

	_, err := s3.New(p.session).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(p.bucketName),
		Key:         aws.String(dst),
		ACL:         aws.String("private"),
		ContentType: aws.String(fileType),
		Body:        fileBytes,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *s3Provider) GetUploadPresignedURL(ctx context.Context) string {
	req, _ := s3.New(p.session).PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(p.bucketName),
		Key:    aws.String(fmt.Sprintf("img/%d", time.Now().UnixNano())),
		ACL:    aws.String("private"),
	})
	//
	url, _ := req.Presign(time.Second * 60)

	return url
}

func (p *s3Provider) GetDomain() string { return p.domain }
func (*s3Provider) GetName() string     { return "aws_s3" }
