package message

import (
	"github.com/google/uuid"
)

// ==========================================================================
// Interfaces
// ==========================================================================

// IMessage define a interface para operações de mensagens
// O parâmetro genérico TData define o tipo de dado que será encapsulado na mensagem
type IMessage[TData any] interface {
	// NewForData cria uma nova mensagem com os dados especificados
	NewForData(correlationId uuid.UUID, data TData, metadata map[string][]byte) (Message[TData], error)

	// NewForDataWithKey cria uma nova mensagem com chave e dados especificados
	NewForDataWithKey(correlationId uuid.UUID, data TData, key string, metadata map[string][]byte) (Message[TData], error)
}
