package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	httpsProxy, err := url.Parse(os.Getenv("HTTPS_PROXY"))
	if err != nil {
		panic(err)
	}
	defer launchProxy(":" + httpsProxy.Port())()

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	// log.Println("GET https://example.com")
	// resp1, err := http.Get("https://example.com")
	// if err != nil {
	// 	panic(err)
	// }
	// body1, _ := io.ReadAll(resp1.Body)
	// log.Println(string(body1))

	// fmt.Println()

	log.Println("GET https://google.com/foo?search=bar")
	resp2, err := http.Get("https://google.com/foo?search=bar")
	if err != nil {
		panic(err)
	}
	body2, _ := io.ReadAll(resp2.Body)
	log.Println(string(body2))

	testS3(http.DefaultClient)
}

func launchProxy(addr string) func() {
	cmd := exec.Command("./_proxy", addr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Millisecond * 500)
	return func() {
		err := cmd.Process.Kill()
		log.Println("kill proxy:", err)
	}
}

func testS3(client *http.Client) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("KEY", "SECRET", "SESSION")),
		config.WithHTTPClient(client),
	)
	if err != nil {
		log.Panicf("unable to load SDK config, %v", err)
	}

	cli := s3.NewFromConfig(cfg)

	_, err = cli.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String("hey-0"),
	})
	if err != nil {
		log.Panic(err)
	}
	_, err = cli.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		log.Panic(err)
	}

	_, err = cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("hey-0"),
		Key:    aws.String("object1"),
		Body:   strings.NewReader("value1"),
	})
	if err != nil {
		log.Panic(err)
	}

	content, err := cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String("hey-0"),
		Key:    aws.String("object1"),
	})
	if err != nil {
		log.Panic(err)
	}
	data, err := io.ReadAll(content.Body)
	if err != nil {
		log.Panic(err)
	}

	log.Println(string(data))
}
