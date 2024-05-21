package snsclient

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/rs/zerolog/log"
)

type snsClient struct {
	svc *sns.SNS
}

func NewSNSClient(svc *sns.SNS) *snsClient {
	return &snsClient{
		svc: svc,
	}
}

func (p *snsClient) SendMessage(topic string, message string) {

	pubInput := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topic),
	}
	output, err := p.svc.Publish(pubInput)
	if err != nil {
		log.Err(err).Msgf("Unable to publish message to [topic: %s], [message: %s]", topic, message)
		return
	}
	log.Info().Msgf("Published message to [topic: %s], [message: %s], [output message id: %s]", topic, *output.MessageId)

}
