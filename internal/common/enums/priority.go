package enums

// ==========================================================================
// Prioridades
// ==========================================================================

// ProducerOrderPriority define os níveis de prioridade disponíveis para produtores
// Controla o comportamento do produtor em relação a performance vs. garantias de entrega
type ProducerOrderPriority string

const (
	// PRODUCER_ORDER_PRIORITY_ORDER prioriza a ordem de chegada das mensagens
	// Garante que as mensagens sejam entregues na ordem em que foram enviadas
	// Modo mais seguro, mas com menor desempenho devido às confirmações síncronas
	PRODUCER_ORDER_PRIORITY_ORDER ProducerOrderPriority = "ORDER"

	// PRODUCER_ORDER_PRIORITY_BALANCED equilibra entre ordem e desempenho
	// Compromisso entre ordenação e throughput
	// Usa confirmações assíncronas, mas aguarda respostas antes de enviar lotes grandes
	PRODUCER_ORDER_PRIORITY_BALANCED ProducerOrderPriority = "BALANCED"

	// PRODUCER_ORDER_PRIORITY_HIGH_PERFORMANCE prioriza o desempenho sobre a ordem
	// Máximo throughput com menor latência, sem garantias rígidas de ordenação
	// Usa confirmações totalmente assíncronas para otimizar o desempenho
	PRODUCER_ORDER_PRIORITY_HIGH_PERFORMANCE ProducerOrderPriority = "HIGH_PERFORMANCE"
)

// ConsumerOrderPriority define os níveis de prioridade disponíveis para consumidores
// Controla o comportamento do consumidor em relação a performance vs. garantias de processamento
type ConsumerOrderPriority string

const (
	// CONSUMER_ORDER_PRIORITY_ORDER prioriza a ordem de processamento das mensagens
	// Processa uma mensagem por vez e só avança após confirmação
	// Garante processamento em ordem e sem perda, mas com throughput reduzido
	CONSUMER_ORDER_PRIORITY_ORDER ConsumerOrderPriority = "ORDER"

	// CONSUMER_ORDER_PRIORITY_BALANCED equilibra entre ordem e desempenho
	// Processa pequenos lotes de mensagens, confirmando-os em conjunto
	// Bom compromisso entre desempenho e confiabilidade
	CONSUMER_ORDER_PRIORITY_BALANCED ConsumerOrderPriority = "BALANCED"

	// CONSUMER_ORDER_PRIORITY_HIGH_PERFORMANCE prioriza o desempenho sobre a ordem
	// Processa grandes lotes de mensagens em paralelo para máximo throughput
	// Pode não garantir a ordem exata de processamento dentro de partições
	CONSUMER_ORDER_PRIORITY_HIGH_PERFORMANCE ConsumerOrderPriority = "HIGH_PERFORMANCE"

	// CONSUMER_ORDER_PRIORITY_RISKY prioriza máximo desempenho com risco de perda de mensagens
	// Modo extremo com confirmação automática e processamento paralelo agressivo
	// Adequado apenas para cenários onde perda ocasional de mensagens é aceitável
	// Não recomendado para dados críticos ou transacionais
	CONSUMER_ORDER_PRIORITY_RISKY ConsumerOrderPriority = "RISKY"
)
