# Kafka Toolkit Lib

Biblioteca Go para integração robusta, tipada e flexível com Apache Kafka, suportando serialização Avro, JSON e Protobuf, além de integração com Schema Registry. Focada em produtividade, segurança e performance, abstrai detalhes de configuração e oferece APIs simples para publicação e consumo de mensagens.

## Principais Recursos
- **Publicação e consumo fortemente tipados** (genéricos)
- **Serialização/Deserialização**: Avro, JSON, Protobuf
- **Integração transparente com Schema Registry**
- **Configuração via variáveis de ambiente** (usando Viper)
- **Gerenciamento de prioridades de performance e consistência** para producers e consumers
- **Thread-safe e singleton** para publishers e consumers

## Instalação

Adicione ao seu projeto Go:
```sh
go get github.com/Dieg657/kafka-toolkit-lib
```

## Configuração

A biblioteca utiliza variáveis de ambiente para configuração. Veja abaixo o detalhamento de cada uma:

| Variável                           | Descrição                                                    | Valores Possíveis / Exemplo                | Default                | Obrigatório? | Comportamento se ausente |
|------------------------------------|--------------------------------------------------------------|--------------------------------------------|------------------------|--------------|-------------------------|
| **KAFKA_BROKERS**                  | Lista de brokers Kafka (endpoints)                           | host1:9092,host2:9092                      | —                      | Sim          | erro                    |
| **KAFKA_GROUPID**                  | Identificador do grupo de consumidores                       | string                                     | —                      | Sim          | erro                    |
| **KAFKA_USERNAME**                 | Usuário SASL para autenticação                               | string                                     | —                      | Sim*         | erro*                   |
| **KAFKA_PASSWORD**                 | Senha SASL para autenticação                                 | string                                     | —                      | Sim*         | erro*                   |
| **KAFKA_SASL_MECHANISM**           | Mecanismo SASL                                               | PLAIN, SCRAM-SHA-256, SCRAM-SHA-512, etc   | PLAIN                  | Não          | Usa default             |
| **KAFKA_SECURITY_PROTOCOL**        | Protocolo de segurança                                       | PLAINTEXT, SASL_PLAINTEXT, SSL, SASL_SSL   | PLAINTEXT              | Não          | Usa default             |
| **KAFKA_SCHEMA_REGISTRY_URL**      | URL do Schema Registry                                       | http(s)://host:8081                        | —                      | Sim††        | erro                    |
| **KAFKA_SCHEMA_REGISTRY_USERNAME** | Usuário do Schema Registry                                   | string                                     | —                      | Condicional† | erro†                   |
| **KAFKA_SCHEMA_REGISTRY_PASSWORD** | Senha do Schema Registry                                     | string                                     | —                      | Condicional† | erro†                   |
| **KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE**| Fonte de credencial do Schema Registry                     | USER_INFO, SASL_INHERIT, NONE              | USER_INFO              | Não          | Usa default             |
| **KAFKA_TIMEOUT**                  | Timeout de requisição (ms)                                   | inteiro > 0                                | 5000                   | Não          | Usa default             |
| **KAFKA_PRODUCER_PRIORITY**        | Prioridade do producer                                       | ORDER, BALANCED, HIGH_PERFORMANCE          | ORDER                  | Não          | Usa default             |
| **KAFKA_CONSUMER_PRIORITY**        | Prioridade do consumer                                       | ORDER, BALANCED, HIGH_PERFORMANCE, RISKY   | ORDER                  | Não          | Usa default             |
| **KAFKA_AUTO_OFFSET_RESET**        | Offset inicial                                               | EARLIEST, LATEST, BEGINNING, END, etc      | LATEST                 | Não          | Usa default             |

> \* Obrigatório apenas se o protocolo SASL exigir autenticação (ex: PLAIN, SCRAM, etc). Para protocolos sem autenticação (plaintext), essas variáveis são ignoradas.
>
> † A obrigatoriedade das credenciais do Schema Registry segue estas regras:
> - Se KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE="USER_INFO" (valor padrão), então KAFKA_SCHEMA_REGISTRY_USERNAME e KAFKA_SCHEMA_REGISTRY_PASSWORD são obrigatórios
> - Se KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE="SASL_INHERIT", as credenciais SASL do Kafka (KAFKA_USERNAME e KAFKA_PASSWORD) serão utilizadas, e as credenciais específicas do Schema Registry são ignoradas
> - Se o Schema Registry não requerer autenticação, use KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE="" (string vazia)
> - Se KAFKA_SCHEMA_REGISTRY_USERNAME e KAFKA_SCHEMA_REGISTRY_PASSWORD forem informados sem especificar KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE, o valor "USER_INFO" será assumido por padrão
>
> †† **Obrigatório apenas para serialização Avro, Protobuf e JSON Schema.** Para JSON puro (`JsonSerialization`/`JsonDeserialization`), não é necessário configurar o Schema Registry.

> **🟢 DICA DE USO: Preenchimento Flexível e Robusto!**
>
> A biblioteca agora inclui **normalização automática de parâmetros**, permitindo que você preencha as variáveis de ambiente em **maiúsculas**, **minúsculas** ou **misturando** (ex: `ORDER`, `order`, `Order`).
> A biblioteca faz o mapeamento automaticamente para o valor esperado pelo Kafka, tornando a configuração muito mais amigável e à prova de erros de digitação/capitalização.
>
> **Exemplos práticos de preenchimento:**
>
> - Prioridade do Producer/Consumer:
>   - `KAFKA_PRODUCER_PRIORITY=ORDER`
>   - `KAFKA_PRODUCER_PRIORITY=order`
>   - `KAFKA_PRODUCER_PRIORITY=Order`
>   - `KAFKA_CONSUMER_PRIORITY=HIGH_PERFORMANCE`
>   - `KAFKA_CONSUMER_PRIORITY=high_performance`
>   - `KAFKA_CONSUMER_PRIORITY=High_Performance`
>
> - Protocolo de Segurança:
>   - `KAFKA_SECURITY_PROTOCOL=PLAINTEXT`
>   - `KAFKA_SECURITY_PROTOCOL=plaintext`
>   - `KAFKA_SECURITY_PROTOCOL=Plaintext`
>
> - Mecanismo SASL:
>   - `KAFKA_SASL_MECHANISM=PLAIN`
>   - `KAFKA_SASL_MECHANISM=plain`
>   - `KAFKA_SASL_MECHANISM=Plain`
>
> - Offset:
>   - `KAFKA_AUTO_OFFSET_RESET=EARLIEST`
>   - `KAFKA_AUTO_OFFSET_RESET=earliest`
>   - `KAFKA_AUTO_OFFSET_RESET=Earliest`
>
> **Não importa a capitalização!**
> 
> A nova implementação garante a normalização completa dos valores, tornando a configuração à prova de erros e permitindo maior flexibilidade na integração com diferentes sistemas de configuração e ambientes.

### Detalhes e Observações
- **Obrigatórios**: KAFKA_BROKERS, KAFKA_GROUPID, e, dependendo do modo de serialização, KAFKA_SCHEMA_REGISTRY_URL e credenciais do Schema Registry.
- **Para JSON puro** (`JsonSerialization`/`JsonDeserialization`):
  - **NÃO** é necessário configurar o Schema Registry nem suas credenciais.
- **Para Avro, Protobuf e JSON Schema**:
  - O Schema Registry e suas credenciais podem ser obrigatórios conforme explicado acima.
- **Erros**: Se uma variável obrigatória estiver ausente, a inicialização retorna um erro com mensagem descritiva.
- **Defaults**: Variáveis não obrigatórias assumem valores padrão seguros para facilitar o uso em ambientes de desenvolvimento.
- **Valores Inválidos**: Se um valor inválido for informado, o sistema tenta usar o default ou retorna erro se não for possível.
- **Exemplo de configuração mínima para ambiente sem autenticação**:
  ```env
  KAFKA_BROKERS=localhost:9092
  KAFKA_GROUPID=meu-grupo
  KAFKA_SCHEMA_REGISTRY_URL=http://localhost:8081
  KAFKA_SECURITY_PROTOCOL=plaintext
  ```
- **Exemplo de configuração mínima para JSON puro (sem Schema Registry)**:
  ```env
  KAFKA_BROKERS=localhost:9092
  KAFKA_GROUPID=meu-grupo
  KAFKA_SECURITY_PROTOCOL=plaintext
  ```
- **Exemplo de configuração completa com autenticação**:
  ```env
  KAFKA_BROKERS=broker1:9092,broker2:9092
  KAFKA_GROUPID=app-prod
  KAFKA_SCHEMA_REGISTRY_URL=https://schema-registry:8081
  KAFKA_USERNAME=usuario
  KAFKA_PASSWORD=senha
  KAFKA_SASL_MECHANISM=SCRAM-SHA-256
  KAFKA_SECURITY_PROTOCOL=sasl_ssl
  KAFKA_SCHEMA_REGISTRY_USERNAME=usuario
  KAFKA_SCHEMA_REGISTRY_PASSWORD=senha
  KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=USER_INFO
  KAFKA_TIMEOUT=10000
  KAFKA_PRODUCER_PRIORITY=ORDER
  KAFKA_CONSUMER_PRIORITY=HIGH_PERFORMANCE
  KAFKA_AUTO_OFFSET_RESET=earliest
  ```

### Verificação Automática de Valores
A biblioteca realiza verificação automática dos valores informados e faz o mapeamento de valores amigáveis para os valores esperados pelo Kafka. Por exemplo:

- Para o **KAFKA_AUTO_OFFSET_RESET**, os valores válidos são:
  - `EARLIEST`, `BEGINNING`, `SMALLEST` (início do tópico)
  - `LATEST`, `END`, `LARGEST` (fim do tópico)
  - `ERROR` (gera erro se não houver offset inicial)

- Para o **KAFKA_SASL_MECHANISM**, os valores válidos são:
  - `PLAIN`, `SCRAM-SHA-256`, `SCRAM-SHA-512`, `GSSAPI`, `OAUTHBEARER`, `NONE`

- Para o **KAFKA_SECURITY_PROTOCOL**, os valores válidos são:
  - `PLAINTEXT`, `SASL_PLAINTEXT`, `SSL`, `SASL_SSL`

- Para o **KAFKA_PRODUCER_PRIORITY** e **KAFKA_CONSUMER_PRIORITY**, os valores válidos são:
  - `ORDER`, `BALANCED`, `HIGH_PERFORMANCE`, `RISKY` (consumer)
  - `ORDER`, `BALANCED`, `HIGH_PERFORMANCE` (producer)

> **Você pode usar qualquer capitalização (ex: `order`, `ORDER`, `Order`). A biblioteca converte automaticamente para o valor correto.**

## Exemplo de Uso

### 1. Inicialização do Container de Dependências
```go
import (
    "context"
    "github.com/Dieg657/kafka-toolkit-lib/internal/common/ioc"
    "github.com/Dieg657/kafka-toolkit-lib/internal/common/constants"
)

ctx := context.Background()
iocContainer, err := ioc.NewKafkaIoC(ctx)
if err != nil {
    panic(err)
}
ctx = context.WithValue(ctx, constants.IocKey, iocContainer)
```

### 2. Publicando Mensagens
```go
import (
    "github.com/Dieg657/kafka-toolkit-lib/internal/publisher"
    "github.com/Dieg657/kafka-toolkit-lib/internal/common/message"
    "github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
    "github.com/google/uuid"
)

type MyPayload struct {
    Field1 string
    Field2 int
}

payload := MyPayload{Field1: "foo", Field2: 42}
correlationId := uuid.New()
msg, _ := message.NewForData(correlationId, payload, map[string][]byte{})
err := publisher.PublishMessage(ctx, "meu-topico", msg, enums.JsonSerialization)
if err != nil {
    panic(err)
}
```

### 3. Consumindo Mensagens
```go
import (
    "github.com/Dieg657/kafka-toolkit-lib/internal/consumer"
    "github.com/Dieg657/kafka-toolkit-lib/internal/common/message"
    "github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
)

type MyPayload struct {
    Field1 string
    Field2 int
}

handler := func(msg message.Message[MyPayload]) error {
    // Processa a mensagem
    fmt.Println(msg.Data.Field1, msg.Data.Field2)
    return nil
}

err := consumer.ConsumeMessage[MyPayload](ctx, "meu-topico", enums.JsonDeserialization, enums.OnDeserializationIgnoreMessage, handler)
if err != nil {
    panic(err)
}
```

## Estratégias de Deserialização

A biblioteca oferece estratégias flexíveis para lidar com falhas de deserialização durante o processamento de mensagens, permitindo diferentes níveis de tolerância a falhas conforme a criticidade do seu sistema.

### Estratégias Disponíveis

#### `OnDeserializationFailedStopHost`
**Comportamento**: Para completamente o consumer ao encontrar erro de deserialização.

**Características**:
- Mais restritivo e seguro para ambientes críticos
- Garante que nenhuma mensagem malformada seja ignorada
- Requer intervenção manual para resolver o problema

**Uso recomendado**:
- Ambientes de produção críticos onde dados corrompidos são inaceitáveis
- Sistemas que requerem processamento garantido de todas as mensagens
- Cenários onde é preferível parar o sistema a processar dados inválidos

**Exemplo**:
```go
err := consumer.ConsumeMessage(ctx, "meu-topico", 
    enums.JsonDeserialization, 
    enums.OnDeserializationFailedStopHost, // Para ao encontrar erro
    handler)
```

#### `OnDeserializationIgnoreMessage`
**Comportamento**: Ignora a mensagem que falhou na deserialização e continua processando.

**Características**:
- Mais tolerante e resiliente a falhas
- Permite que o sistema continue funcionando mesmo com mensagens malformadas
- Requer logging adequado para rastrear mensagens ignoradas

**Uso recomendado**:
- Ambientes onde alta disponibilidade é mais importante que processamento garantido
- Sistemas que podem tolerar perda ocasional de mensagens malformadas
- Cenários de desenvolvimento e teste
- Processamento em batch onde algumas falhas são aceitáveis

**Exemplo**:
```go
err := consumer.ConsumeMessage(ctx, "meu-topico", 
    enums.JsonDeserialization, 
    enums.OnDeserializationIgnoreMessage, // Ignora erros e continua
    handler)
```

### Quando Usar Cada Estratégia

| Cenário | Estratégia Recomendada | Justificativa |
|---------|------------------------|---------------|
| **Sistema Financeiro** | `OnDeserializationFailedStopHost` | Dados corrompidos podem causar inconsistências críticas |
| **Analytics/Logs** | `OnDeserializationIgnoreMessage` | Volume alto, algumas perdas são aceitáveis |
| **E-commerce (Carrinho)** | `OnDeserializationFailedStopHost` | Transações requerem integridade completa |
| **Métricas/Monitoramento** | `OnDeserializationIgnoreMessage` | Disponibilidade é mais importante que precisão absoluta |
| **Ambiente de Desenvolvimento** | `OnDeserializationIgnoreMessage` | Facilita testes e não interrompe desenvolvimento |
| **Compliance/Auditoria** | `OnDeserializationFailedStopHost` | Regulamentações exigem processamento completo |

### Melhores Práticas

1. **Logging**: Sempre implemente logging adequado para rastrear mensagens ignoradas
2. **Monitoramento**: Configure alertas para falhas de deserialização frequentes
3. **Ambiente**: Use estratégias mais restritivas em produção
4. **Schema Evolution**: Planeje evolução de schemas para minimizar incompatibilidades
5. **Testes**: Teste ambas as estratégias em ambiente de desenvolvimento

## Formatos Suportados

### Serialização/Deserialização
- **Json**
  - enums.JsonSerialization / enums.JsonDeserialization
  - **JSON puro, sem Schema Registry**.
  - Dados são convertidos para JSON sem validação de schema.
  - **Vantagens**: Facilidade de uso, flexibilidade, legibilidade humana.
  - **Desvantagens**: Sem validação estrutural, maior tamanho de payload, sem controle de schema.
  - **Recomendado para**: Integração simples, sistemas legados, ambientes de desenvolvimento, ou quando flexibilidade é mais importante que validação.
  - Exemplo:
    ```go
    enums.JsonSerialization // para publicar
    enums.JsonDeserialization // para consumir
    ```
- **JSON Schema**
  - enums.JsonSchemaSerialization / enums.JsonSchemaDeserialization
  - Dados JSON validados contra um schema JSON definido e registrado.
  - **Vantagens**: Validação estrutural, legibilidade humana, evolução controlada.
  - **Desvantagens**: Tamanho maior que formatos binários, performance moderada.
  - **Requer**: Schema Registry configurado.
  - **Recomendado para**: APIs REST, integrações web, quando é necessário equilíbrio entre legibilidade e validação de schema.
  - Exemplo:
    ```go
    enums.JsonSchemaSerialization // para publicar
    enums.JsonSchemaDeserialization // para consumir
    ```
- **Avro**
  - enums.AvroSerialization / enums.AvroDeserialization
  - Formato binário compacto com schemas definidos, ótimo desempenho e economia de espaço.
  - **Vantagens**: Alta compactação, evolução de schema retrocompatível, desempenho excelente.
  - **Desvantagens**: Não é legível por humanos, requer ferramentas específicas para visualização.
  - **Requer**: Schema Registry configurado.
  - **Recomendado para**: Alto volume de dados, sistemas de processamento de dados, analytics, quando eficiência de armazenamento é crítica.
  - Exemplo:
    ```go
    enums.AvroSerialization // para publicar
    enums.AvroDeserialization // para consumir
    ```
- **Protobuf**
  - enums.ProtobufSerialization / enums.ProtobufDeserialization
  - Formato binário compacto e eficiente baseado em schemas (.proto).
  - **Vantagens**: Extremamente rápido, tamanho compacto, forte tipagem, retrocompatibilidade, suporte multi-linguagem.
  - **Desvantagens**: Requer ferramentas específicas, curva de aprendizado inicial.
  - **Requer**: Schema Registry configurado.
  - **Recomendado para**: Sistemas de alta performance, comunicação entre serviços, microservices, aplicações sensíveis à latência.
  - Exemplo:
    ```go
    enums.ProtobufSerialization // para publicar
    enums.ProtobufDeserialization // para consumir
    ```

## Adaptador Protobuf

### Problemas comuns com Protobuf em Go e como a biblioteca resolve para você

Ao integrar Protobuf com Kafka em Go, muitos desenvolvedores enfrentam problemas como:

- **Incompatibilidade entre implementações de Protobuf**: O ecossistema Go possui duas principais implementações (`github.com/golang/protobuf` e `google.golang.org/protobuf`), que não são 100% compatíveis entre si. Isso pode gerar erros difíceis de diagnosticar.
- **Mensagens de erro confusas**: É comum encontrar erros como:
  - `serialization target must be a protobuf message`
  - `cannot use *MyProto as type proto.Message`
- **Necessidade de adaptação manual**: Em outras bibliotecas, o desenvolvedor precisa adaptar manualmente as mensagens ou converter entre tipos, o que gera código extra, manutenção difícil e risco de bugs sutis.

#### Como a Kafka Toolkit Lib resolve para você

- **Adaptação automática**: A biblioteca detecta automaticamente qual implementação de Protobuf você está usando e faz toda a ponte necessária, sem exigir nenhuma configuração extra do usuário.
- **Transparência total**: Você pode usar suas structs Protobuf normalmente, sem se preocupar com detalhes de compatibilidade.
- **Sem código extra**: Não é necessário criar adaptadores, wrappers ou conversões manuais.
- **Erros resolvidos**: Aqueles erros típicos de incompatibilidade simplesmente não acontecem aqui.

> **Resumo:**
> Basta usar suas mensagens Protobuf normalmente.

### Funcionamento Automático
Na maioria dos casos, o adaptador funcionará automaticamente sem configuração adicional:

- Detecta automaticamente tipos protobuf usando reflexão e heurísticas
- Realiza cópias de campos entre diferentes implementações
- Mantém cache para otimizar o desempenho

### Exemplo de Uso com Protobuf
```go
import (
    "github.com/Dieg657/kafka-toolkit-lib/internal/publisher"
    "github.com/Dieg657/kafka-toolkit-lib/internal/common/message"
    "github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
    "github.com/google/uuid"
    pb "seu/pacote/protobuf"
)

// Publicando
payload := &pb.MeuProtobuf{
    Campo1: "teste",
    Campo2: 123,
}
correlationId := uuid.New()
msg, _ := message.NewForData(correlationId, payload, nil)
err := publisher.PublishMessage(ctx, "meu-topico", msg, enums.ProtobufSerialization)

// Consumindo
handler := func(msg message.Message[*pb.MeuProtobuf]) error {
    fmt.Printf("Mensagem recebida: %+v\n", msg.Data)
    return nil
}
err := consumer.ConsumeMessage(ctx, "meu-topico", enums.ProtobufDeserialization, enums.OnDeserializationIgnoreMessage, handler)
```

### Prioridades Producer

> **⚠️ Nota sobre Valores Default:**
> O valor padrão para tanto `KAFKA_PRODUCER_PRIORITY` quanto `KAFKA_CONSUMER_PRIORITY` é **ORDER**. 
> Isto prioriza segurança e consistência por padrão, garantindo ordem de entrega e processamento.
> Para aplicações que precisam de mais performance, configure explicitamente como `BALANCED` ou `HIGH_PERFORMANCE`.

- **ORDER** (padrão)
  - Garante ordem de entrega e idempotência.
  - Uso: sistemas financeiros, logs ordenados, eventos críticos.
  - **Configuração automática da biblioteca:**
    ```yaml
    acks: all
    enable.idempotence: true
    max.in.flight.requests.per.connection: 1
    retries: 5
    linger.ms: 0
    batch.num.messages: 1000
    compression.type: snappy
    message.timeout.ms: 120000
    request.timeout.ms: 30000
    queue.buffering.max.messages: 100000
    queue.buffering.max.kbytes: 1048576
    retry.backoff.ms: 100
    queue.buffering.max.ms: 0
    ```
  - **O que cada parâmetro faz:**
    - `acks: all` — Garante que todas as réplicas confirmem a mensagem antes de considerar entregue (máxima confiabilidade).
    - `enable.idempotence: true` — Evita duplicidade de mensagens mesmo em falhas de rede.
    - `max.in.flight.requests.per.connection: 1` — Garante ordem absoluta das mensagens.
    - `retries: 5` — Tenta reenviar mensagens em caso de falha.
    - `linger.ms: 0` — Não espera para formar lotes, priorizando baixa latência.
    - `batch.num.messages: 1000` — Limita o tamanho dos lotes para controle de memória.
    - `compression.type: snappy` — Compressão eficiente para reduzir uso de rede.
    - `message.timeout.ms: 120000` — Tempo máximo para tentar entregar uma mensagem.
    - `request.timeout.ms: 30000` — Timeout para requisições ao broker.
    - `queue.buffering.max.messages`/`max.kbytes` — Controlam o buffer local do produtor.
    - `retry.backoff.ms: 100` — Tempo de espera entre tentativas de reenvio.
    - `queue.buffering.max.ms: 0` — Não atrasa entregas para otimizar lotes.
  - **Exemplo:**
    ```env
    KAFKA_PRODUCER_PRIORITY=ORDER
    ```
  - **Quando escolher:**
    - Processos que não podem tolerar duplicidade ou reordenação (ex: débito em conta, emissão de nota fiscal).
    - Workflows que dependem de causalidade estrita.

- **BALANCED**
  - Equilíbrio entre performance e consistência.
  - **Configuração automática da biblioteca:**
    ```yaml
    acks: all
    enable.idempotence: true
    max.in.flight.requests.per.connection: 5
    retries: 3
    linger.ms: 5
    batch.num.messages: 10000
    compression.type: lz4
    message.timeout.ms: 120000
    request.timeout.ms: 30000
    queue.buffering.max.messages: 200000
    queue.buffering.max.kbytes: 2097152
    retry.backoff.ms: 100
    queue.buffering.max.ms: 5
    ```
  - **O que cada parâmetro faz:**
    - `acks: all` — Confirmação de todas as réplicas, mas com mais paralelismo.
    - `enable.idempotence: true` — Segurança contra duplicidade.
    - `max.in.flight.requests.per.connection: 5` — Permite paralelismo moderado, balanceando ordem e performance.
    - `retries: 3` — Tolerância a falhas.
    - `linger.ms: 5` — Pequeno atraso para formar lotes maiores (melhor throughput).
    - `batch.num.messages: 10000` — Lotes maiores para eficiência.
    - `compression.type: lz4` — Compressão rápida e eficiente.
    - `message.timeout.ms`/`request.timeout.ms` — Controle de tempo de entrega.
    - `queue.buffering.max.messages`/`max.kbytes` — Buffers maiores para mais throughput.
    - `retry.backoff.ms: 100` — Espera entre tentativas.
    - `queue.buffering.max.ms: 5` — Pequeno atraso para otimizar lotes.
  - **Exemplo:**
    ```env
    KAFKA_PRODUCER_PRIORITY=BALANCED
    ```
  - **Quando escolher:**
    - Aplicações de uso geral, integrações, eventos de negócio não críticos.

- **HIGH_PERFORMANCE**
  - Máximo throughput, menos garantias de ordem e duplicidade.
  - **Configuração automática da biblioteca:**
    ```yaml
    acks: 1
    enable.idempotence: false
    max.in.flight.requests.per.connection: 10
    retries: 1
    linger.ms: 10
    batch.num.messages: 50000
    compression.type: lz4
    message.timeout.ms: 60000
    request.timeout.ms: 20000
    queue.buffering.max.messages: 500000
    queue.buffering.max.kbytes: 4194304
    retry.backoff.ms: 100
    socket.keepalive.enable: true
    queue.buffering.max.ms: 10
    ```
  - **O que cada parâmetro faz:**
    - `acks: 1` — Confirmação só do líder, priorizando velocidade.
    - `enable.idempotence: false` — Permite duplicidade para máximo throughput.
    - `max.in.flight.requests.per.connection: 10` — Altíssimo paralelismo.
    - `retries: 1` — Pouca tolerância a falhas.
    - `linger.ms: 10` — Espera para formar grandes lotes.
    - `batch.num.messages: 50000` — Lotes enormes para eficiência máxima.
    - `compression.type: lz4` — Compressão rápida.
    - `message.timeout.ms`/`request.timeout.ms` — Tempos menores para agilidade.
    - `queue.buffering.max.messages`/`max.kbytes` — Buffers gigantes para throughput.
    - `retry.backoff.ms: 100` — Espera curta entre tentativas.
    - `socket.keepalive.enable: true` — Mantém conexões ativas.
    - `queue.buffering.max.ms: 10` — Pequeno atraso para otimizar lotes.
  - **Exemplo:**
    ```env
    KAFKA_PRODUCER_PRIORITY=HIGH_PERFORMANCE
    ```
  - **Quando escolher:**
    - Coleta de logs, métricas, analytics, pipelines de dados massivos.
    - Situações onde performance é mais importante que ordem ou unicidade.
  - **Dica:** Combine com particionamento customizado para priorizar buckets ou classes de mensagens ([Bucket Priority Pattern](https://www.confluent.io/blog/prioritize-messages-in-kafka/)).

### Prioridades Consumer
- **ORDER** (padrão)
  - Consumo ordenado, maior segurança e controle.
  - **Configuração automática da biblioteca:**
    ```yaml
    enable.auto.commit: false
    auto.offset.reset: EARLIEST
    isolation.level: read_committed
    max.poll.interval.ms: 300000
    session.timeout.ms: 10000
    heartbeat.interval.ms: 3000
    fetch.min.bytes: 1
    fetch.wait.max.ms: 500
    retry.backoff.ms: 100
    fetch.error.backoff.ms: 500
    max.partition.fetch.bytes: 1048576
    fetch.max.bytes: 5242880
    ```
  - **O que cada parâmetro faz:**
    - `enable.auto.commit: false` — Commit manual, máxima segurança.
    - `auto.offset.reset: EARLIEST` — Consome desde o início se não houver offset salvo.
    - `isolation.level: read_committed` — Lê apenas mensagens confirmadas.
    - `max.poll.interval.ms: 300000` — Tempo máximo para processar lote.
    - `session.timeout.ms: 10000` — Timeout de sessão.
    - `heartbeat.interval.ms: 3000` — Frequência de heartbeat.
    - `fetch.min.bytes: 1` — Busca mensagem assim que disponível.
    - `fetch.wait.max.ms: 500` — Espera máxima para buscar.
    - `retry.backoff.ms: 100` — Espera entre tentativas.
    - `fetch.error.backoff.ms: 500` — Espera após erro de fetch.
    - `max.partition.fetch.bytes`/`fetch.max.bytes` — Controlam tamanho dos lotes e uso de memória.
  - **Exemplo:**
    ```env
    KAFKA_CONSUMER_PRIORITY=ORDER
    ```
  - **Quando escolher:**
    - Processamento de eventos financeiros, pipelines ETL sensíveis à ordem.

- **BALANCED**
  - Equilíbrio entre performance e consistência.
  - **Configuração automática da biblioteca:**
    ```yaml
    enable.auto.commit: false
    auto.offset.reset: EARLIEST
    isolation.level: read_committed
    max.poll.interval.ms: 300000
    session.timeout.ms: 30000
    heartbeat.interval.ms: 10000
    fetch.min.bytes: 1024
    fetch.wait.max.ms: 1000
    retry.backoff.ms: 200
    fetch.error.backoff.ms: 500
    max.partition.fetch.bytes: 1048576
    fetch.max.bytes: 10485760
    fetch.message.max.bytes: 262144
    ```
  - **O que cada parâmetro faz:**
    - `enable.auto.commit: false` — Commit manual, mais controle.
    - `auto.offset.reset: EARLIEST` — Consome desde o início se necessário.
    - `isolation.level: read_committed` — Lê apenas mensagens confirmadas.
    - `max.poll.interval.ms`/`session.timeout.ms` — Controlam tempo de processamento e sessão.
    - `heartbeat.interval.ms` — Frequência de heartbeat.
    - `fetch.min.bytes`/`fetch.wait.max.ms` — Controlam batching e latência.
    - `retry.backoff.ms`/`fetch.error.backoff.ms` — Resiliência a falhas.
    - `max.partition.fetch.bytes`/`fetch.max.bytes`/`fetch.message.max.bytes` — Controlam tamanho dos lotes e uso de memória.
  - **Exemplo:**
    ```env
    KAFKA_CONSUMER_PRIORITY=BALANCED
    ```
  - **Quando escolher:**
    - Consumo de eventos de negócio, integrações, sistemas de workflow.

- **HIGH_PERFORMANCE**
  - Consumo rápido, menos garantias de ordem.
  - **Configuração automática da biblioteca:**
    ```yaml
    enable.auto.commit: true
    auto.commit.interval.ms: 5000
    auto.offset.reset: LATEST
    isolation.level: read_uncommitted
    max.poll.interval.ms: 600000
    session.timeout.ms: 60000
    heartbeat.interval.ms: 20000
    fetch.min.bytes: 65536
    fetch.wait.max.ms: 100
    retry.backoff.ms: 50
    fetch.error.backoff.ms: 200
    max.partition.fetch.bytes: 10485760
    fetch.max.bytes: 52428800
    fetch.message.max.bytes: 1048576
    ```
  - **O que cada parâmetro faz:**
    - `enable.auto.commit: true` — Commit automático, mais performance.
    - `auto.commit.interval.ms: 5000` — Commit frequente.
    - `auto.offset.reset: LATEST` — Consome apenas novas mensagens.
    - `isolation.level: read_uncommitted` — Lê todas as mensagens, mesmo não confirmadas.
    - `max.poll.interval.ms`/`session.timeout.ms` — Tempos maiores para processar grandes lotes.
    - `heartbeat.interval.ms` — Heartbeat menos frequente.
    - `fetch.min.bytes`/`fetch.wait.max.ms` — Lotes grandes, menor latência.
    - `retry.backoff.ms`/`fetch.error.backoff.ms` — Resiliência a falhas.
    - `max.partition.fetch.bytes`/`fetch.max.bytes`/`fetch.message.max.bytes` — Lotes e buffers grandes para throughput.
  - **Exemplo:**
    ```env
    KAFKA_CONSUMER_PRIORITY=HIGH_PERFORMANCE
    ```
  - **Quando escolher:**
    - Analytics, logs, ingestão massiva de dados.
    - Situações onde performance é mais importante que ordem ou confiabilidade.
  - **Dica:** Combine com estratégias de particionamento e assignors customizados para priorização de buckets ([Bucket Priority Pattern](https://www.confluent.io/blog/prioritize-messages-in-kafka/)).

- **RISKY**
  - Máximo throughput, mínimo de garantias (auto-commit, menos consistência).
  - **Configuração automática da biblioteca:**
    ```yaml
    enable.auto.commit: true
    auto.commit.interval.ms: 1000
    auto.offset.reset: LATEST
    isolation.level: read_uncommitted
    max.poll.interval.ms: 900000
    session.timeout.ms: 120000
    heartbeat.interval.ms: 40000
    fetch.min.bytes: 131072
    fetch.wait.max.ms: 50
    retry.backoff.ms: 10
    fetch.error.backoff.ms: 100
    reconnect.backoff.ms: 10
    max.partition.fetch.bytes: 52428800
    fetch.max.bytes: 104857600
    fetch.message.max.bytes: 4194304
    ```
  - **O que cada parâmetro faz:**
    - `enable.auto.commit: true` — Commit automático, máxima velocidade.
    - `auto.commit.interval.ms: 1000` — Commit muito frequente.
    - `auto.offset.reset: LATEST` — Consome apenas novas mensagens.
    - `isolation.level: read_uncommitted` — Lê todas as mensagens.
    - `max.poll.interval.ms`/`session.timeout.ms` — Tempos longos para processar grandes lotes.
    - `heartbeat.interval.ms` — Heartbeat menos frequente.
    - `fetch.min.bytes`/`fetch.wait.max.ms` — Lotes enormes, mínima latência.
    - `retry.backoff.ms`/`fetch.error.backoff.ms`/`reconnect.backoff.ms` — Resiliência mínima, foco em velocidade.
    - `max.partition.fetch.bytes`/`fetch.max.bytes`/`fetch.message.max.bytes` — Buffers e lotes gigantes para máximo throughput.
  - **Exemplo:**
    ```env
    KAFKA_CONSUMER_PRIORITY=RISKY
    ```
  - **Quando escolher:**
    - Monitoramento, métricas, logs de auditoria não críticos.
    - Processamento best-effort, onde o volume é mais importante que a precisão.

> **Dica Avançada:**
> Para cenários de priorização real de mensagens, utilize padrões como custom partitioners e assignors (ex: Bucket Priority Pattern). Isso permite que diferentes consumidores ou buckets recebam fatias específicas do throughput, mesmo dentro do mesmo consumer group, otimizando recursos e garantindo SLAs diferenciados. Saiba mais em: [Prioritize Messages in Kafka](https://www.confluent.io/blog/prioritize-messages-in-kafka/)

### Offset
- **EARLIEST**
  - Consome desde o início do tópico.
  - Inicia o consumo desde a mensagem mais antiga disponível no tópico.
  - Garante que todas as mensagens disponíveis dentro da política de retenção sejam processadas.
  - **Recomendado para**: Reprocessamento, bootstrap de dados, ou inicialização de novos consumidores que precisam processar todo o histórico.
- **LATEST** (padrão)
  - Consome apenas novas mensagens que chegarem após a conexão do consumidor.
  - Ignora todas as mensagens existentes no tópico antes da conexão.
  - **Recomendado para**: A maioria dos casos de uso em produção, especialmente para processamento em tempo real onde o histórico não é relevante.
- **BEGINNING**
  - Sinônimo de EARLIEST.
  - Similar a "earliest", mas pode ter comportamento diferente em algumas implementações.
  - **Nota técnica**: Consulte a documentação do Kafka para detalhes específicos da sua versão.
- **END**
  - Sinônimo de LATEST.
  - Similar a "latest", mas pode ter comportamento diferente em algumas implementações.
  - **Nota técnica**: Consulte a documentação do Kafka para detalhes específicos da sua versão.
- **SMALLEST**
  - Termo legado equivalente a "earliest".
  - Mantido para compatibilidade com versões antigas do Kafka.
  - **Recomendado para**: Código legado ou compatibilidade com clientes antigos.
- **LARGEST**
  - Termo legado equivalente a "latest".
  - Mantido para compatibilidade com versões antigas do Kafka.
  - **Recomendado para**: Código legado ou compatibilidade com clientes antigos.
- **ERROR**
  - Gera erro se não houver offset válido salvo para o consumer group.
  - A aplicação receberá um erro e deverá tratá-lo manualmente.
  - **Recomendado para**: Debugging, validação de fluxos, e cenários onde o processamento precisa ser explicitamente controlado.

### Protocolos de Segurança
- **PLAINTEXT**
  - Sem criptografia ou autenticação.
  - **ATENÇÃO**: NÃO RECOMENDADO para ambientes de produção ou qualquer ambiente exposto à rede pública.
  - **Recomendado apenas para**: Ambientes de desenvolvimento totalmente isolados e testes locais.
- **SASL_PLAINTEXT**
  - Implementa autenticação SASL mas sem criptografia.
  - As credenciais são autenticadas, mas as mensagens trafegam em texto claro.
  - **Recomendado apenas para**: Redes internas totalmente isoladas e seguras.
- **SSL**
  - Comunicação criptografada via SSL/TLS.
  - Fornece criptografia do tráfego, mas sem autenticação SASL.
  - A autenticação pode ser implementada via certificados de cliente.
  - **Recomendado para**: Ambientes de produção onde a autenticação é feita por outros meios.
- **SASL_SSL**
  - Combina autenticação SASL e criptografia SSL/TLS.
  - Oferece o máximo de segurança com autenticação robusta e tráfego criptografado.
  - **Recomendado para**: Ambientes de produção, dados sensíveis, e comunicação através de redes públicas.

### SASL Mechanisms
- **PLAIN**
  - Usuário/Senha simples. Fácil de configurar, menos seguro.
  - **Atenção**: Credenciais são transmitidas em texto claro. Utilize apenas em conjunto com SSL/TLS.
  - **Recomendado para**: Ambientes de desenvolvimento ou quando a segurança é implementada na camada de transporte.
- **SCRAM-SHA-256**
  - Autenticação forte baseada em hash SHA-256.
  - Oferece melhor segurança que PLAIN sem transmitir senhas em texto claro.
  - **Recomendado para**: Maioria dos ambientes de produção com requisitos moderados de segurança.
- **SCRAM-SHA-512**
  - Autenticação forte baseada em hash SHA-512.
  - Oferece segurança superior ao SHA-256, mas maior sobrecarga computacional.
  - **Recomendado para**: Ambientes de alta segurança e dados sensíveis.
- **GSSAPI**
  - Implementa autenticação Kerberos/GSSAPI.
  - Adequado para ambientes com infraestrutura Kerberos existente.
  - **Recomendado para**: Ambientes corporativos e integração com Active Directory/MIT Kerberos.
- **OAUTHBEARER**
  - Implementa autenticação baseada em OAuth 2.0.
  - Adequado para integração com provedores de identidade externos.
  - **Recomendado para**: Ambientes cloud e sistemas que já utilizam OAuth para autenticação federada.
- **NONE**
  - Sem mecanismo SASL, apenas para conexões sem autenticação.
  - **Recomendado para**: Apenas ambientes isolados, desenvolvimento local e testes.

### Schema Registry Auth Source
- **USER_INFO**
  - Usa usuário/senha informados nas variáveis de ambiente.
- **SASL_INHERIT**
  - Herda autenticação SASL do Kafka.
- **""** (vazio)
  - Sem autenticação (apenas para Schema Registry aberto).

> **Dica:** Sempre valide as opções de configuração conforme o ambiente (dev, staging, prod) e as políticas de segurança da sua organização.

## Playground - Teste antes de integrar

A biblioteca inclui um ambiente de playground completo para que você possa experimentar todas as funcionalidades sem precisar integrar imediatamente ao seu projeto. O playground permite testar a biblioteca em um ambiente local controlado com Kafka, Zookeeper, Schema Registry e Kafka UI.

### Ambientes Disponíveis

O playground oferece duas configurações diferentes:

1. **Ambiente Simples (sem SASL)** - `docker-compose.yaml`
   - Sem autenticação (usa PLAINTEXT)
   - Ideal para testes rápidos e demonstrações
   - Tópicos são criados automaticamente

2. **Ambiente com SASL** - `docker-compose.sasl.yaml`
   - Com autenticação SASL/PLAIN
   - Similar a ambientes corporativos/produtivos
   - Demonstra cenários de autenticação

### Exemplos Prontos para Uso

O playground inclui exemplos para todos os formatos de serialização suportados:

- **Avro**: `playground/avro_example/main.go`
- **Protobuf**: `playground/protobuf_example/main.go`
- **JSON Schema**: `playground/json_schema_example/main.go`
- **JSON**: `playground/json_example/main.go`

### Benchmark de Produção

**🆕 NOVIDADE**: A biblioteca agora inclui um benchmark especializado para medir a performance de **produção** de mensagens Kafka.

- **Localização**: `playground/benchmark/`
- **Foco**: Exclusivamente na performance de producers
- **Métricas**: Throughput, latência (média, P95, P99), taxa de sucesso
- **Análise Inteligente**: Classificação automática de performance
- **Relatórios**: Detalhados com recomendações de otimização
- **Estimativas**: Projeções de capacidade por hora/dia

#### Como usar o Benchmark:
```sh
cd playground/benchmark/example
go run main.go
```

#### Exemplo de Relatório:
```
🚀 PRODUÇÃO:
   Total de mensagens: 15,642
   Mensagens com sucesso: 15,642 (100.00%)
   Throughput: 521.40 mensagens/segundo
   Latência média: 2.15 ms

📊 ANÁLISE DE PERFORMANCE:
   Taxa de Sucesso: ✅ EXCELENTE (100.00%)
   Throughput: ⚠️  MÉDIO (521 msg/s)
   Latência: ✅ BAIXA (2.15 ms)

💡 RECOMENDAÇÕES:
   - Ajuste batch.size (ex: 16384 ou 32768)
   - Configure linger.ms (ex: 5-10ms)

📊 ESTIMATIVA DE CAPACIDADE:
   - Por hora: ~1,877,040 mensagens
   - Por dia: ~45,048,960 mensagens
```

#### Configuração Personalizada:
```go
config := benchmark.NewDefaultConfig()
config.ProducerWorkers = 8                      // 8 producers
config.Duration = 60 * time.Second              // 1 minuto
config.TopicName = "meu-topico-benchmark"       // Tópico customizado

result := benchmark.RunProducerOnlyBenchmark(ctx, config)
```

> Para documentação completa do benchmark, consulte: `playground/benchmark/README.md`

### Como Utilizar o Playground

#### 1. Clone o repositório
```sh
git clone https://github.com/Dieg657/kafka-toolkit-lib.git
cd kafka-toolkit-lib
```

#### 2. Inicie o ambiente Docker
```sh
# Ambiente sem autenticação
docker compose -f playground/docker-compose.yaml up -d

# OU para ambiente com autenticação
docker compose -f playground/docker-compose.sasl.yaml up -d
```

#### 3. Configure as variáveis de ambiente

Para ambiente sem autenticação:
```sh
# Linux/macOS
export KAFKA_BROKERS=localhost:29093
export KAFKA_GROUPID=my-test-group
export KAFKA_SCHEMA_REGISTRY_URL=http://localhost:8081
export KAFKA_SECURITY_PROTOCOL=plaintext
export KAFKA_SASL_MECHANISM=PLAIN
export KAFKA_SCHEMA_REGISTRY_USERNAME=schema-registry
export KAFKA_SCHEMA_REGISTRY_PASSWORD=schema-registry-password
export KAFKA_SHEMA_REGISTRY_AUTH_SOURCE=""

# Windows (PowerShell)
$env:KAFKA_BROKERS="localhost:29093"
$env:KAFKA_GROUPID="my-test-group"
$env:KAFKA_SCHEMA_REGISTRY_URL="http://localhost:8081"
$env:KAFKA_SECURITY_PROTOCOL="plaintext"
$env:KAFKA_SASL_MECHANISM="PLAIN"
$env:KAFKA_SCHEMA_REGISTRY_USERNAME="schema-registry"
$env:KAFKA_SCHEMA_REGISTRY_PASSWORD="schema-registry-password"
```

#### 4. Execute um dos exemplos
```sh
go run playground/avro_example/main.go
# OU
go run playground/protobuf_example/main.go
# OU
go run playground/json_schema_example/main.go
# OU
go run playground/json_example/main.go
# OU - NOVO BENCHMARK
go run playground/benchmark/example/main.go
```

#### 5. Explore a interface web
O ambiente inclui uma interface web Kafka UI para visualizar tópicos e mensagens:
- Acesse: http://localhost:8080

### Benefícios do Playground

- **Teste Completo**: Experimente todos os recursos em um ambiente isolado
- **Validação de Conceitos**: Entenda como a biblioteca funciona antes de integrar
- **Ambiente Realista**: Inclui todas as dependências (Kafka, Schema Registry)
- **Aprendizado Prático**: Exemplos prontos e funcionais para cada formato
- **Avaliação Técnica**: Valide se a biblioteca atende às necessidades do seu projeto
- **🆕 Benchmark de Performance**: Teste e otimize a performance de produção

> Para mais detalhes e configurações avançadas, consulte a documentação do playground: `/playground/README.md`

## Estrutura da Mensagem

A biblioteca utiliza uma estrutura genérica tipada para representar mensagens Kafka:

```go
type Message[TData any] struct {
    CorrelationId uuid.UUID
    Data          TData
    Metadata      map[string][]byte
}
```

### Detalhamento dos Campos

#### CorrelationId (uuid.UUID)
- **Finalidade**: Identificador único de correlação que permite rastrear uma mensagem através de diferentes sistemas e componentes
- **Importância**: Fundamental para observabilidade, rastreamento de fluxos e depuração em arquiteturas distribuídas
- **Geração**: Normalmente gerado com `uuid.New()` durante a criação da mensagem
- **Uso**: Deve ser propagado entre sistemas, permitindo correlacionar logs, métricas e traces
- **Caso de Uso**: Ideal para acompanhar o caminho de uma requisição que passa por múltiplos serviços

#### Data (TData)
- **Finalidade**: Contém o payload principal da mensagem com tipo fortemente definido
- **Flexibilidade**: Utiliza genéricos do Go para garantir tipagem segura sem comprometer flexibilidade
- **Compatibilidade**: Suporta qualquer tipo Go que possa ser serializado (struct, map, slice, tipos primitivos)
- **Vantagens**: 
  - Evita type assertions e conversões manuais
  - Garante segurança de tipos em tempo de compilação
  - Melhora legibilidade e manutenção do código
- **Caso de Uso**: Transferência de dados estruturados entre produtores e consumidores

#### Metadata (map[string][]byte)
- **Finalidade**: Armazenar informações adicionais que não pertencem ao payload principal, mas são relevantes para processamento, rastreamento ou contexto
- **Conteúdo Típico**:
  - Timestamps de processamento
  - Identificadores de origem/destino
  - Informações de roteamento
  - Flags de controle
  - Dados de instrumentação
  - Versões de schema/API
- **Natureza**: Opcional, pode ser nil ou vazio quando não necessário
- **Flexibilidade**: Campos em bytes permitem qualquer tipo de dado serializado, inclusive binários
- **Caso de Uso**: Passar headers HTTP originais, dados de autenticação, ou informações de rastreamento

### Exemplos de Uso da Estrutura

#### Exemplo Básico
```go
// Criando uma mensagem simples
payload := MyStruct{Field1: "valor", Field2: 42}
correlationId := uuid.New()
msg, err := message.NewForData(correlationId, payload, nil)
```

#### Utilizando Metadados
```go
// Criando uma mensagem com metadados
payload := MyStruct{Field1: "valor", Field2: 42}
metadata := map[string][]byte{
    "source":      []byte("api-gateway"),
    "request-id":  []byte("req-123456"),
    "timestamp":   []byte(time.Now().Format(time.RFC3339)),
    "trace-token": []byte("trace-abc-xyz"),
}

correlationId := uuid.New()
msg, err := message.NewForData(correlationId, payload, metadata)
```

#### Acessando Metadados no Consumidor
```go
handler := func(msg message.Message[MyStruct]) error {
    // Acessando o payload tipado
    fmt.Println("Dados:", msg.Data.Field1, msg.Data.Field2)
    
    // Acessando metadados
    if source, ok := msg.Metadata["source"]; ok {
        fmt.Println("Origem:", string(source))
    }
    
    if traceToken, ok := msg.Metadata["trace-token"]; ok {
        // Usar o token para correlacionar com sistemas de tracing
        fmt.Println("Token de trace:", string(traceToken))
    }
    
    return nil
}
```

### Melhores Práticas

1. **Correlação**: Sempre reutilize o CorrelationId ao repassar mensagens para outros sistemas
2. **Consistência**: Defina padrões de nomenclatura para chaves de metadados em toda sua organização
3. **Tamanho**: Mantenha metadados pequenos e focados em informações essenciais
4. **Evolução**: Projete sua estrutura de Data pensando em evolução e versionamento
5. **Segurança**: Não armazene informações sensíveis em metadados, a menos que estejam criptografadas

## Compatibilidade

A biblioteca foi testada e é compatível com as seguintes versões:

- Apache Kafka: 2.x, 3.x
- Confluent Schema Registry: 5.x, 6.x, 7.x
- Confluent Kafka Go Client: v2.x
- Go: 1.18+

Para garantir a melhor experiência, recomendamos:
- Confluent Platform 7.x ou superior
- Go 1.19 ou superior
- Schema Registry com suporte a múltiplos formatos de serialização (Avro, Protobuf, JSON Schema)

### Configuração do Schema Registry

> **Atenção:**
> Se você informar `KAFKA_SCHEMA_REGISTRY_USERNAME` e `KAFKA_SCHEMA_REGISTRY_PASSWORD` mas **não definir explicitamente** o valor de `KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE`, a biblioteca irá **assumir automaticamente** o valor `USER_INFO` como default.
>
> Isso garante que as credenciais fornecidas sejam usadas corretamente para autenticação básica no Schema Registry, sem necessidade de configuração extra.
>
> **Exemplo prático:**
> ```env
> KAFKA_SCHEMA_REGISTRY_URL=http://localhost:8081
> KAFKA_SCHEMA_REGISTRY_USERNAME=meu-usuario
> KAFKA_SCHEMA_REGISTRY_PASSWORD=minha-senha
> # KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE não informado
> # A biblioteca assume USER_INFO automaticamente
> ```

#### Valores possíveis para `KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE` e seus efeitos

- **`USER_INFO`** (default se username/senha forem informados)
  - **O que faz:** Usa as credenciais fornecidas em `KAFKA_SCHEMA_REGISTRY_USERNAME` e `KAFKA_SCHEMA_REGISTRY_PASSWORD` para autenticação básica no Schema Registry.
  - **Quando usar:** Quando o Schema Registry tem usuário/senha próprios, diferentes do Kafka.
  - **Exemplo:**
    ```env
    KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=USER_INFO
    KAFKA_SCHEMA_REGISTRY_USERNAME=usuario
    KAFKA_SCHEMA_REGISTRY_PASSWORD=senha
    ```

- **`SASL_INHERIT`**
  - **O que faz:** Herda as credenciais do Kafka (`KAFKA_USERNAME` e `KAFKA_PASSWORD`) para autenticação no Schema Registry.
  - **Quando usar:** Quando o Schema Registry compartilha as mesmas credenciais do cluster Kafka (ambientes corporativos, SSO, etc).
  - **Exemplo:**
    ```env
    KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=SASL_INHERIT
    KAFKA_USERNAME=usuario
    KAFKA_PASSWORD=senha
    # Não precisa informar KAFKA_SCHEMA_REGISTRY_USERNAME/KAFKA_SCHEMA_REGISTRY_PASSWORD
    ```

- **`""`** (string vazia)
  - **O que faz:** Não utiliza autenticação para o Schema Registry (acesso aberto).
  - **Quando usar:** Quando o Schema Registry está aberto/publicamente acessível e não exige autenticação.
  - **Exemplo:**
    ```env
    KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=""
    # Não precisa informar usuário/senha
    ```

> **Resumo:**
> - Informe `USER_INFO` para usar credenciais específicas do Schema Registry.
> - Use `SASL_INHERIT` para herdar as credenciais do Kafka.
> - Use `""` (vazio) para Schema Registry sem autenticação.

## Exemplo de Uso

## Thread-safety e Singleton: Publisher e Consumer

> **🚦 Informação Crítica: Publisher e Consumer são Singleton e Thread-safe!**
>
> A biblioteca garante que tanto o **Publisher** quanto o **Consumer** são:
> - **Singleton por tipo de dado** (e por tópico, quando aplicável)
> - **Thread-safe**: podem ser usados em múltiplas goroutines sem risco de múltiplas instâncias ou condições de corrida
>
> Isso é possível graças ao uso de `sync.Map` e do padrão `LoadOrStore`, que garantem que:
> - Só existe **um Publisher por tipo** (`TData`) em toda a aplicação
> - Só existe **um Producer por tipo+tópico**
> - Só existe **um Consumer por tipo+tópico**
> - Só existe **um engine consumer por tipo+tópico**
>
> **Por que isso é importante?**
> - Evita múltiplas conexões desnecessárias com o Kafka
> - Garante uso eficiente de recursos
> - Elimina bugs de concorrência e race conditions
> - Permite máxima performance e escalabilidade
> - Facilita a integração em aplicações de alta concorrência (microserviços, workers, etc)
>
> **Exemplo prático:**
>
> ```go
> // Pode usar em quantas goroutines quiser, sempre será singleton e thread-safe!
> for i := 0; i < 100; i++ {
>     go func() {
>         publisher.PublishMessage(ctx, "meu-topico", msg, enums.JsonSerialization)
>         consumer.ConsumeMessage(ctx, "meu-topico", enums.JsonDeserialization, handler)
>     }()
> }
> ```
>
> **Diferencial:**
> - Esse design robusto é um dos grandes diferenciais da biblioteca e pode ser decisivo na escolha para projetos que exigem alta confiabilidade e performance.

## Testes de Singleton e Concorrência

A biblioteca inclui testes automatizados que garantem o padrão singleton mesmo sob concorrência intensa. Os testes lançam múltiplas goroutines que tentam obter o IoC container, o producer e o consumer simultaneamente, e validam que **todas as instâncias retornadas são idênticas**.

### Como rodar os testes

Para rodar todos os testes do projeto:
```sh
go test ./...
```

Para rodar apenas os testes do IoC container:
```sh
go test ./internal/common/ioc
```

Para ver o output detalhado:
```sh
go test -v ./internal/common/ioc
```

> Se algum teste de singleton falhar, isso indica que múltiplas instâncias estão sendo criadas, o que não deve acontecer.

## 🆕 Últimas Atualizações

### Benchmark de Produção Kafka
- **Localização**: `playground/benchmark/`
- **Funcionalidade**: Benchmark especializado para medir performance de producers
- **Métricas**: Throughput, latência (média, P95, P99), taxa de sucesso
- **Análise Inteligente**: Classificação automática de performance com recomendações
- **Estimativas**: Projeções de capacidade por hora/dia
- **Relatórios**: Detalhados com sugestões de otimização específicas

### Estratégias de Deserialização
- **Funcionalidade**: Controle flexível de comportamento em falhas de deserialização
- **Estratégias Disponíveis**:
  - `OnDeserializationFailedStopHost`: Para o sistema em caso de erro
  - `OnDeserializationIgnoreMessage`: Ignora mensagens com erro e continua
- **Uso**: Permite configurar tolerância a falhas conforme criticidade do sistema
- **Documentação**: Guia completo com cenários de uso e melhores práticas

### Melhorias na API
- **Consumer**: Assinatura atualizada para incluir estratégia de deserialização obrigatória
- **Documentação**: Enums completamente documentados seguindo padrões da biblioteca
- **Exemplos**: Todos os exemplos do playground atualizados com as novas funcionalidades

### Como Atualizar

Para usar as novas funcionalidades, atualize suas chamadas de `ConsumeMessage`:

**Antes:**
```go
err := consumer.ConsumeMessage(ctx, "topico", enums.JsonDeserialization, handler)
```

**Agora:**
```go
err := consumer.ConsumeMessage(ctx, "topico", enums.JsonDeserialization, enums.OnDeserializationIgnoreMessage, handler)
```

> **Nota**: A estratégia de deserialização agora é **obrigatória** para maior controle e segurança.