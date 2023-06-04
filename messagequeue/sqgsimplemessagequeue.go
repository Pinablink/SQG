package messagequeue

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gofrs/uuid"
)

//
type SQGSimpleMessageQueue struct {
	sqgMessageQueueUUID uuid.UUID
	strQueueURL         *string
	messageEnvelope     map[string]types.MessageAttributeValue
	messageBody         string
}

//
func NewSimpleSQGMessageQueue(strMessageBody string, refStrQueueURL *string, mapMessageEnvelope map[string]types.MessageAttributeValue) SQGSimpleMessageQueue {
	m_uuid, _ := uuid.NewV4()

	return SQGSimpleMessageQueue{
		sqgMessageQueueUUID: m_uuid,
		strQueueURL:         refStrQueueURL,
		messageEnvelope:     mapMessageEnvelope,
		messageBody:         strMessageBody,
	}
}

//
func (ref SQGSimpleMessageQueue) ProduceMessage() *sqs.SendMessageInput {
	return &sqs.SendMessageInput{
		DelaySeconds:      10,
		MessageAttributes: ref.messageEnvelope,
		MessageBody:       aws.String(ref.messageBody),
		QueueUrl:          ref.strQueueURL,
	}
}
