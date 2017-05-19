package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func sendSnsMessageStruct(message interface{}) {
	jsonMsg, _ := json.Marshal(message)
	sendSnsMessageString(string(jsonMsg))
}

func sendSnsMessageString(message string) {
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
		return
	}

	fmt.Println("sent SNS message: " + message)
	fmt.Println(resp)
}
