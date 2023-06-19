package sqg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Pinablink/sqg/messagequeue"
	"github.com/Pinablink/sqg/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gofrs/uuid"
)

//
type SQG struct {
	sqgUUID   uuid.UUID
	nameQueue string
	sQGStruct util.GSQGMessageModel
}

//
func NewSQG(strnameQueue string) *SQG {

	m_uuid, _ := uuid.NewV4()

	return &SQG{
		sqgUUID:   m_uuid,
		nameQueue: strnameQueue,
	}

}

//
func (ref *SQG) SetGSQGMessageModel(refMessageModel util.GSQGMessageModel) {
	ref.sQGStruct = refMessageModel
}

//
func (ref *SQG) parseStruct() (error, string) {

	dataObjectByte, err := json.Marshal(ref.sQGStruct.ContentMessage)

	if err != nil {

		messageErrorReturn := fmt.Sprintf(util.NOT_SERIALIZE_STR_BASE64, ref.sqgUUID)
		return errors.New(messageErrorReturn), ""
	}

	return nil, string(dataObjectByte)
}

// JoinTheQueue : Adiciona os dados de util.GSQGMessageModel a fila SQS
func (ref *SQG) JoinTheQueue() (*string, error) {

	var refCtxAwsConfig aws.Config
	var client *sqs.Client
	var sqsQueueUrlOutput *sqs.GetQueueUrlOutput
	var mapAttrMsgEnvelope map[string]types.MessageAttributeValue
	var simpleMessageQueue messagequeue.SQGSimpleMessageQueue
	var messageInput *sqs.SendMessageInput
	var messageOutput *sqs.SendMessageOutput
	var idmessage *string

	mErr, strContentMessageBody := ref.parseStruct()

	if mErr == nil {

		mapAttrMsgEnvelope, mErr = messagequeue.ParserMessage(ref.sQGStruct.DataMessage)

		if mErr == nil {

			refCtxAwsConfig, mErr = config.LoadDefaultConfig(context.TODO())
			client = sqs.NewFromConfig(refCtxAwsConfig)

			queueURLInput := &sqs.GetQueueUrlInput{
				QueueName: &ref.nameQueue,
			}

			sqsQueueUrlOutput, mErr = client.GetQueueUrl(context.TODO(), queueURLInput)

			if mErr == nil {
				ref_queueURL := sqsQueueUrlOutput.QueueUrl
				simpleMessageQueue = messagequeue.NewSimpleSQGMessageQueue(strContentMessageBody, ref_queueURL, mapAttrMsgEnvelope)
				messageInput = simpleMessageQueue.ProduceMessage()
				messageOutput, mErr = client.SendMessage(context.TODO(), messageInput)
				idmessage = messageOutput.MessageId
			}

		}

	}

	return idmessage, mErr
}

// getDeleteMessage Recupera a mensagem da fila SQS
// contentRef: Referência do Conteúdo da Mensagem
// dataAttrRef: Dados de cabeçalho da Mensagem
// reviewMessageOK: Informa se ocorreu leitura de uma mensagem
// deleteMessageOk: Informa se a mensagem foi deletada
// errorGetDelMessage: Erro de processamento
func (ref *SQG) getDeleteMessage(contentRef interface{}, dataAttrRef interface{}) (reviewMessageOK bool, deleteMessageOk bool, errorGetDelMessage error) {
	var iErr error
	var refCtxAwsConfig aws.Config
	var client *sqs.Client
	var urlRef *sqs.GetQueueUrlOutput
	var messagesOutput *sqs.ReceiveMessageOutput
	var reviewMessage bool = false
	var deleteMessageOK bool = false

	refCtxAwsConfig, iErr = config.LoadDefaultConfig(context.TODO())

	if iErr == nil {
		client = sqs.NewFromConfig(refCtxAwsConfig)

		queueURLInput := &sqs.GetQueueUrlInput{
			QueueName: &ref.nameQueue,
		}

		urlRef, iErr = client.GetQueueUrl(context.TODO(), queueURLInput)

		if iErr == nil {

			strQueueURL := urlRef.QueueUrl

			receiveMessages := &sqs.ReceiveMessageInput{
				MessageAttributeNames: []string{
					string(types.QueueAttributeNameAll),
				},
				QueueUrl:            strQueueURL,
				MaxNumberOfMessages: int32(1),
			}

			messagesOutput, iErr = client.ReceiveMessage(context.TODO(), receiveMessages)

			if iErr == nil {

				var mMessages []types.Message = messagesOutput.Messages

				if mMessages != nil {
					var strReceiptBody *string = mMessages[0].Body
					var mapAtributes map[string]types.MessageAttributeValue = mMessages[0].MessageAttributes
					var strReceiptHandler *string = mMessages[0].ReceiptHandle

					unMarshalError := messagequeue.ReturnData(contentRef, dataAttrRef, mapAtributes, []byte(*strReceiptBody))

					if unMarshalError == nil {
						var deleteMessageInput *sqs.DeleteMessageInput = &sqs.DeleteMessageInput{
							QueueUrl:      strQueueURL,
							ReceiptHandle: strReceiptHandler,
						}

						reviewMessage = true

						_, delError := client.DeleteMessage(context.TODO(), deleteMessageInput)

						if delError == nil {
							deleteMessageOK = true
						}

					} else {
						iErr = unMarshalError
					}

				}

			} else {
				var strMessageError = fmt.Sprintf(util.NOT_FOUND_MESSAGE_PROCESS, iErr.Error())
				iErr = errors.New(strMessageError)
			}

		}

	}

	return reviewMessage, deleteMessageOK, iErr
}

// GetMessage Obtêm a mensagem da fila e é um Wrapper para a func getDeleteMessage
// contentRef: Referência do Conteúdo da Mensagem
// dataAttrRef: Dados de cabeçalho da Mensagem
// reviewMessageOK: Informa se ocorreu leitura de uma mensagem
// deleteMessageOk: Informa se a mensagem foi deletada
// errorGetDelMessage: Erro de processamento
func (ref *SQG) GetMsgInQueue(contentRef interface{}, dataAttrRef interface{}) (reviewMessageOK bool, deleteMessageOk bool, errorGetDelMessage error) {
	return ref.getDeleteMessage(contentRef, dataAttrRef)
}
