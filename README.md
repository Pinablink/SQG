# sqg

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Pinablink/tmdbGoTutorial?style=plastic)
![Alt text](https://img.shields.io/badge/AWS-SQS-blue?style=plastic)


É uma lib escrita em linguagem **Golang** e tem como objetivo, abstrair algumas configurações para o SQS da AWS. Desse modo toda a atenção e empenho podem estar na codificação da sua solução especifica.

### Pontos de desenvolvimento

Uso da Lib da AWS
Uso do conceito de **reflection** e **manipulação do stream** dos dados 

### Dependências

A lib **sqg** utiliza recursos de outros pacotes. O principal pacote é da **AWS**. Mas existe também um pacote especifíco para geração de um código uuid para identificação de referências.


	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gofrs/uuid"

### Como utilizar

Você pode acessar um tutorial criado no medium https://medium.com/@weberasantos/publicando-e-consumindo-mensagem-no-sqs-aws-com-golang-6970a0e7581e, que demonstra a utilização dessa lib na aplicação testeSQSGsqgLambda, que pode ser encontrado em https://github.com/Pinablink/testeSQSGsqgLambda
