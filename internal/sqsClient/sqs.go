package sqsclient

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/rs/zerolog/log"
)

type SQSClient struct {
	svc *sqs.SQS
}

func NewSQSClient(svc *sqs.SQS) *SQSClient {
	return &SQSClient{
		svc: svc,
	}
}
func (c *SQSClient) Subscribe(queue string) {
	for {
		messages, err := c.ReceiveMessage(queue)
		if err != nil {
			log.Err(err).Msgf("Unable to read message for [queue: %s]", queue)
		}
		for _, msg := range messages {
			if msg == nil {
				continue
			}
			c.DeleteMessage(queue, msg.ReceiptHandle)
		}
	}
}

func (c *SQSClient) ReceiveMessage(queue string) ([]*sqs.Message, error) {
	receiveMessageInput := &sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(queue),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(3),
	}
	receiveMessageOutput, err := c.svc.ReceiveMessage(receiveMessageInput)
	if receiveMessageOutput.Messages == nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if receiveMessageOutput == nil || receiveMessageOutput.Messages == nil {
		return nil, err
	}
	return receiveMessageOutput.Messages, nil
}

func (c *SQSClient) DeleteMessage(queue string, handle *string) {
	deleteInput := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queue),
		ReceiptHandle: handle,
	}
	_, err := c.svc.DeleteMessage(deleteInput)
	if err != nil {
		log.Err(err).Msgf("Unable to delete message for [queue: %s]", queue)
		return
	}

}
