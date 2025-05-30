package enums

// ==========================================================================
// Configurações de Offset
// ==========================================================================

// AutoOffsetReset define os comportamentos disponíveis para redefinição automática de offset
// Determina o comportamento do consumidor quando não há offset válido salvo para um grupo
type AutoOffsetReset string

const (
	// OFFSET_RESET_ERROR gera erro quando não há offset válido
	// A aplicação receberá um erro e deverá tratá-lo manualmente
	// Útil quando o processamento deve ser explicitamente controlado
	OFFSET_RESET_ERROR AutoOffsetReset = "error"

	// OFFSET_RESET_SMALLEST define offset para a menor posição disponível (legado)
	// Termo legado equivalente a "earliest"
	// Mantido para compatibilidade com versões antigas
	OFFSET_RESET_SMALLEST AutoOffsetReset = "smallest"

	// OFFSET_RESET_EARLIEST define offset para a posição mais antiga disponível
	// Inicia o consumo desde a mensagem mais antiga disponível
	// Garante que todas as mensagens disponíveis serão processadas
	OFFSET_RESET_EARLIEST AutoOffsetReset = "earliest"

	// OFFSET_RESET_BEGINNING define offset para o início da partição
	// Similar a "earliest", mas pode ter comportamento diferente em algumas implementações
	// Consulte a documentação do Kafka para detalhes específicos da versão
	OFFSET_RESET_BEGINNING AutoOffsetReset = "beginning"

	// OFFSET_RESET_LARGEST define offset para a maior posição disponível (legado)
	// Termo legado equivalente a "latest"
	// Mantido para compatibilidade com versões antigas
	OFFSET_RESET_LARGEST AutoOffsetReset = "largest"

	// OFFSET_RESET_LATEST define offset para a posição mais recente disponível
	// Inicia o consumo apenas de novas mensagens que chegarem após a conexão
	// Ignora todas as mensagens existentes no tópico
	OFFSET_RESET_LATEST AutoOffsetReset = "latest"

	// OFFSET_RESET_END define offset para o final da partição
	// Similar a "latest", mas pode ter comportamento diferente em algumas implementações
	// Consulte a documentação do Kafka para detalhes específicos da versão
	OFFSET_RESET_END AutoOffsetReset = "end"
)
