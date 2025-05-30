# Benchmark de Produção - Kafka Toolkit Lib

Este benchmark testa o desempenho de **produção** da biblioteca Kafka Toolkit Lib, medindo throughput, latência e taxa de sucesso dos producers.

## 🎯 Características Principais

✅ **Benchmark de Produção**: Foco exclusivo na performance de envio de mensagens  
✅ **Métricas Detalhadas**: Throughput, latência (média, P95, P99) e taxa de sucesso  
✅ **Análise Inteligente**: Classificação automática de performance e recomendações  
✅ **Uso Padrão da Biblioteca**: Segue exatamente o padrão do `json_example`  
✅ **Mensagens JSON**: Serialização JSON simples (sem Schema Registry)  
✅ **Relatórios Detalhados**: Estatísticas completas com análises e recomendações  
✅ **Estimativa de Capacidade**: Projeções de throughput por hora/dia  

## 📊 Métricas Coletadas

### Produção (Producer)
- **Throughput**: Mensagens/segundo produzidas
- **Latência de Produção**: Tempo para enviar ao Kafka (média, P95, P99)
- **Taxa de Sucesso**: Percentual de mensagens enviadas com sucesso
- **Taxa de Falhas**: Mensagens que falharam no envio

### Análise de Performance
- **Classificação de Throughput**: Alto (≥1000 msg/s), Médio (≥500 msg/s), Baixo (<500 msg/s)
- **Classificação de Latência**: Baixa (≤5ms), Média (≤20ms), Alta (>20ms)
- **Classificação de Sucesso**: Excelente (≥99.5%), Boa (≥95%), Baixa (<95%)
- **Recomendações Automáticas**: Sugestões específicas de otimização

## 🚀 Como Usar

### Execução Básica
```bash
cd playground/benchmark/example
go run main.go
```

### Configuração Personalizada
```go
config := benchmark.NewDefaultConfig()
config.ProducerWorkers = 8                      // 8 producers
config.Duration = 60 * time.Second              // 1 minuto
config.TopicName = "meu-topico-benchmark"       // Tópico customizado

result := benchmark.RunProducerOnlyBenchmark(ctx, config)
```

## 📋 Exemplo de Relatório

```
======================================================================
                 RESULTADOS DO BENCHMARK DE PRODUÇÃO
======================================================================

🚀 PRODUÇÃO:
   Total de mensagens: 15,642
   Mensagens com sucesso: 15,642 (100.00%)
   Mensagens falharam: 0 (0.00%)
   Throughput: 521.40 mensagens/segundo
   Latência média: 2.15 ms
   Latência P95: 8.34 ms
   Latência P99: 15.67 ms

📊 ANÁLISE DE PERFORMANCE:
   Taxa de Sucesso: ✅ EXCELENTE (100.00%)
   Throughput: ⚠️  MÉDIO (521 msg/s)
   Latência: ✅ BAIXA (2.15 ms)

💡 RECOMENDAÇÕES:
   - Considere ajustar batch.size e linger.ms para melhor throughput
   - Verifique se o número de workers está adequado para sua CPU

   Tempo total: 30s
======================================================================

🔍 ANÁLISE DETALHADA:
✅ Taxa de sucesso excelente - sistema funcionando perfeitamente
⚠️  Throughput médio (521 msg/s) - há espaço para otimização

💡 RECOMENDAÇÕES ESPECÍFICAS:
📈 Para melhorar throughput:
   - Ajuste batch.size (ex: 16384 ou 32768)
   - Configure linger.ms (ex: 5-10ms)
   - Considere aumentar buffer.memory
   - Verifique se o número de workers está adequado

📊 ESTIMATIVA DE CAPACIDADE:
   - Por hora: ~1,877,040 mensagens
   - Por dia: ~45,048,960 mensagens

🎯 Benchmark de produção concluído com sucesso em 30s!
```

## ⚙️ Configurações Disponíveis

| Parâmetro | Padrão | Descrição |
|-----------|--------|-----------|
| `ProducerWorkers` | `runtime.NumCPU()` | Número de workers produtores |
| `Duration` | `30s` | Duração do teste |
| `TopicName` | `benchmark-producer-topic` | Nome do tópico Kafka |
| `ReportInterval` | `5s` | Intervalo para relatórios (futuro uso) |
| `PreWarmup` | `true` | Indica se deve fazer warm-up (futuro uso) |

## 🎯 Interpretação dos Resultados

### Taxa de Sucesso
- **≥ 99.5%**: ✅ Excelente - Sistema funcionando perfeitamente
- **≥ 95.0%**: ⚠️ Boa - Algumas falhas ocasionais
- **< 95.0%**: ❌ Baixa - Necessita investigação

### Throughput (Mensagens/segundo)
- **≥ 1000**: ✅ Alto - Excelente performance
- **≥ 500**: ⚠️ Médio - Há espaço para otimização  
- **< 500**: ❌ Baixo - Necessita otimização urgente

### Latência Média
- **≤ 5ms**: ✅ Baixa - Resposta rápida
- **≤ 20ms**: ⚠️ Média - Aceitável para maioria dos casos
- **> 20ms**: ❌ Alta - Pode impactar performance

### Recomendações Automáticas
O benchmark analisa automaticamente os resultados e fornece sugestões específicas:

#### Para Melhorar Throughput:
- Ajustar `batch.size` (ex: 16384 ou 32768)
- Configurar `linger.ms` (ex: 5-10ms)
- Aumentar `buffer.memory`
- Otimizar número de workers

#### Para Reduzir Latência P99:
- Ajustar `request.timeout.ms`
- Configurar `delivery.timeout.ms`
- Verificar configurações de `acks`

#### Para Reduzir Falhas:
- Verificar logs detalhados do Kafka
- Ajustar configurações de retry
- Verificar conectividade de rede

## 🔧 Dependências

O benchmark usa as mesmas dependências da biblioteca principal:
- Kafka Go Client (Confluent)
- Inicialização via IoC Container
- Serialização JSON nativa

## 📝 Estrutura das Mensagens

```go
type JsonMessage struct {
    HoraMensagem string `json:"hora_mensagem"`  // Timestamp formatado
    Producer     string `json:"producer"`       // ID do worker produtor
    MessageID    string `json:"message_id"`     // UUID único
    Counter      int64  `json:"counter"`        // Contador sequencial
    Timestamp    int64  `json:"timestamp"`      // Nanosegundos (para futuras medições)
}
```

## 📊 Estimativa de Capacidade

O benchmark fornece projeções automáticas baseadas no throughput medido:
- **Mensagens por hora**: `throughput × 3600`
- **Mensagens por dia**: `throughput × 86400`

Essas estimativas ajudam no planejamento de capacidade e dimensionamento da infraestrutura.

## 🚨 Considerações

⚠️ **Ambiente de Teste**: Execute em ambiente controlado  
⚠️ **Recursos**: Monitore CPU/memória durante execução  
⚠️ **Kafka**: Certifique-se que o cluster suporte a carga  
⚠️ **Limpeza**: O tópico pode acumular mensagens - considere limpeza após testes  
⚠️ **Baseline**: Execute múltiplas vezes para estabelecer baseline confiável  

## 🔮 Futuras Funcionalidades

- Benchmark de consumo (será adicionado posteriormente)
- Medição de latência end-to-end
- Suporte a diferentes formatos de serialização
- Relatórios em formato JSON/CSV
- Integração com ferramentas de monitoramento 