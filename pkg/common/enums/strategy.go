package enums

// ==========================================================================
// Estratégias de Deserialização
// ==========================================================================

// DeserializationStrategy define as estratégias disponíveis para lidar com falhas de deserialização
// durante o processamento de mensagens Kafka.
//
// Esta enum determina o comportamento do consumer quando encontra uma mensagem que não pode
// ser deserializada corretamente, permitindo diferentes níveis de tolerância a falhas.
type DeserializationStrategy int

const (
	// OnDeserializationFailedStopHost interrompe o processamento e para o host quando
	// uma mensagem não pode ser deserializada.
	//
	// Características:
	// - Para completamente o consumer ao encontrar erro de deserialização
	// - Mais restritivo e seguro para ambientes críticos
	// - Garante que nenhuma mensagem malformada seja ignorada
	// - Requer intervenção manual para resolver o problema
	//
	// Uso recomendado:
	// - Ambientes de produção críticos onde dados corrompidos são inaceitáveis
	// - Sistemas que requerem processamento garantido de todas as mensagens
	// - Cenários onde é preferível parar o sistema a processar dados inválidos
	OnDeserializationFailedStopHost DeserializationStrategy = iota

	// OnDeserializationIgnoreMessage ignora a mensagem que falhou na deserialização
	// e continua processando as próximas mensagens.
	//
	// Características:
	// - Pula mensagens com erro de deserialização e continua o processamento
	// - Mais tolerante e resiliente a falhas
	// - Permite que o sistema continue funcionando mesmo com mensagens malformadas
	// - Requer logging adequado para rastrear mensagens ignoradas
	//
	// Uso recomendado:
	// - Ambientes onde alta disponibilidade é mais importante que processamento garantido
	// - Sistemas que podem tolerar perda ocasional de mensagens malformadas
	// - Cenários de desenvolvimento e teste
	// - Processamento em batch onde algumas falhas são aceitáveis
	OnDeserializationIgnoreMessage DeserializationStrategy = 1
)
