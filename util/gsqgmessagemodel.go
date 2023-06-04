package util

// GSGGMessageModel - Estrutura de Mensagem
// DataMessage - Deverá conter dados para classificar sua mensagem, como por exemplo Titulo e outros dados de sua necessidade
// ContentMessage - Deverá conter os atributos de sua mensagem
type GSQGMessageModel struct {
	DataMessage    interface{}
	ContentMessage interface{}
}
