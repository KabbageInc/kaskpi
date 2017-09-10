package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func getEmitter() io.Writer {
	if conf.Event.EventEmitterName == "HTTP" {
		return RetryWrapperEmitter{Writer: HttpEventEmitter{}}
	}
	return RetryWrapperEmitter{Writer: SnsEventEmitter{}}
}

type SnsEventEmitter struct {
}

func (w SnsEventEmitter) Write(p []byte) (int, error) {
	return len(p), sendSnsMessageString(string(p))
}

func serializeMessage(message interface{}) []byte {
	jsonMsg, _ := json.Marshal(message)
	return jsonMsg
}

func sendSnsMessageStruct(message interface{}) {
	jsonMsg := serializeMessage(message)
	sendSnsMessageString(string(jsonMsg))
}

func sendSnsMessageString(message string) error {
	sess := session.Must(session.NewSession())

	creds := credentials.NewStaticCredentials(conf.Aws.AccessKeyId, conf.Aws.SecretAccessKey, "")

	svc := sns.New(sess, &aws.Config{
		Region:      aws.String(endpoints.UsEast1RegionID),
		Credentials: creds,
	})

	params := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(conf.Aws.SnsArn),
	}

	resp, err := svc.Publish(params)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("sent SNS message: " + message)
	fmt.Println(resp)

	return nil
}

type HttpEventEmitter struct {
}

func (w HttpEventEmitter) Write(p []byte) (int, error) {
	url := conf.Event.HttpEmitter.PostUrl
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(p))
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil || response.StatusCode >= 300 {
		return 0, err
	}
	defer response.Body.Close()

	return len(p), nil
}

type RetryWrapperEmitter struct {
	Writer io.Writer
}

func (w RetryWrapperEmitter) Write(p []byte) (int, error) {
	count, err := w.Writer.Write(p)
	if err != nil {
		//write to db for retry

	}
	return count, err
}
