package sqg

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

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
func (ref *SQG) serializeStruct() (error, string) {
	var mBuffer bytes.Buffer
	var strData string
	encoder := base64.NewEncoder(base64.StdEncoding, &mBuffer)

	dataObjectByte, err := json.Marshal(ref.sQGStruct.ContentMessage)

	if err != nil {

		messageErrorReturn := fmt.Sprintf(util.NOT_SERIALIZE_STR_BASE64, ref.sqgUUID)
		return errors.New(messageErrorReturn), ""

	} else {

		defer encoder.Close()
		strData = string(dataObjectByte)
		encoder.Write([]byte(strData))

	}

	return nil, mBuffer.String()
}

//
func (ref *SQG) deserializeStruct(strBase64Ref string) {

	var aDataByte = []byte(strBase64Ref)
	var reader *bytes.Reader = bytes.NewReader(aDataByte)

	decoder := base64.NewDecoder(base64.StdEncoding, reader)

	// RETORNAR ESSA ESTRUTURA
	//dataByte, err := ioutil.ReadAll(decoder)
	_, err := ioutil.ReadAll(decoder)

	if err != nil {
		// Implementar retorno do erro
		fmt.Println("Erro no Decode")
	}

}

// JoinTheQueue : Adiciona os dados de util.GSQGMessageModel a fila SQS em formato String64
func (ref *SQG) JoinTheQueue() (*string, error) {

	var refCtxAwsConfig aws.Config
	var client *sqs.Client
	var sqsQueueUrlOutput *sqs.GetQueueUrlOutput
	var mapAttrMsgEnvelope map[string]types.MessageAttributeValue
	var simpleMessageQueue messagequeue.SQGSimpleMessageQueue
	var messageInput *sqs.SendMessageInput
	var messageOutput *sqs.SendMessageOutput
	var idmessage *string

	mErr, strContentMessageBody := ref.serializeStruct()

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

//
func (ref *SQG) getDeleteMessage(iQtGetMsgProcess int, v any) bool {
	var iErr error
	var refCtxAwsConfig aws.Config
	var client *sqs.Client
	var urlRef *sqs.GetQueueUrlOutput
	var messagesOutput *sqs.ReceiveMessageOutput
	var reviewMessage bool = true

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
				QueueUrl: strQueueURL,
			}

			for i := 0; i <= iQtGetMsgProcess; i++ {

				messagesOutput, iErr = client.ReceiveMessage(context.TODO(), receiveMessages)

				if iErr == nil {

					var mMessages []types.Message = messagesOutput.Messages

					if mMessages != nil {
						var strReceiptBody *string = mMessages[0].Body
						var strReceiptHandler *string = mMessages[0].ReceiptHandle

						ref.deserializeStruct(*strReceiptBody)

						var deleteMessageInput *sqs.DeleteMessageInput = &sqs.DeleteMessageInput{
							QueueUrl:      strQueueURL,
							ReceiptHandle: strReceiptHandler,
						}

						client.DeleteMessage(context.TODO(), deleteMessageInput)

					} else {
						break
					}

				} else {
					// Implementar aqui um tratamento de erro
					fmt.Println("Ocorreu um ero na obtencao da mensagem")
				}

			}

		}

	}

	return reviewMessage
}

// GetMessage:
// EM PROCESSO DE DESENVOLVIMENTO
func (ref *SQG) GetMessages(qtMessages int, v any) {
	ref.getDeleteMessage(qtMessages, v)
}
