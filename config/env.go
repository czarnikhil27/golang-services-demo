package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AwsAccessKeyID     string `yaml:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey string `yaml:"AWS_SECRET_ACCESS_KEY"`
}
type SNSTopics struct {
	Test string `SNS_TOPIC`
}

type SQSTopic struct {
	SqsTopic1 string `SQS_TOPIC`
	SqsTopic2 string `SQS_TOPIC_2`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
