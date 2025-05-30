package publisher

import (
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/message"
)

// ==========================================================================
// Interfaces
// ==========================================================================

// IPublisher define a interface para publicação de mensagens
// O parâmetro genérico TData define o tipo de dado que será publicado
type IPublisher[TData any] interface {
	// PublishMessage publica uma mensagem no tópico especificado com o formato definido
	PublishMessage(topic string, message message.Message[TData], serialization enums.Serialization) error
}
