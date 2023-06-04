package messagequeue

import "github.com/gofrs/uuid"

//
type SQGReceiveSimpleMessageQueue struct {
	sqgReceiveMsgQueueUUID uuid.UUID
	strQueueURL            *string
}

//
func NewSQGReceiveSimpleMessageQueue(refStrQueueURL *string) SQGReceiveSimpleMessageQueue {
	m_uuid, _ := uuid.NewV4()

	return SQGReceiveSimpleMessageQueue{
		sqgReceiveMsgQueueUUID: m_uuid,
	}
}

//
