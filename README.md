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
| **KAFKA_SECURITY_PROTOCOL**        | Protocolo de segurança                                       | plaintext, sasl_plaintext, ssl, sasl_ssl   | sasl_plaintext         | Não          | Usa default             |
| **KAFKA_SCHEMA_REGISTRY_URL**      | URL do Schema Registry                                       | http(s)://host:8081                        | —                      | Sim          | erro                    |
| **KAFKA_SCHEMA_REGISTRY_USERNAME** | Usuário do Schema Registry                                   | string                                     | —                      | Condicional† | erro†                   |
| **KAFKA_SCHEMA_REGISTRY_PASSWORD** | Senha do Schema Registry                                     | string                                     | —                      | Condicional† | erro†                   |
| **KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE**| Fonte de credencial do Schema Registry                      | USER_INFO, SASL_INHERIT, ""                | USER_INFO              | Não          | Usa default             |
| **KAFKA_TIMEOUT**                  | Timeout de requisição (ms)                                   | inteiro > 0                                | 5000                   | Não          | Usa default             |
| **KAFKA_PRODUCER_PRIORITY**        | Prioridade do producer                                       | ORDER, BALANCED, HIGH_PERFORMANCE          | BALANCED               | Não          | Usa default             |
| **KAFKA_CONSUMER_PRIORITY**        | Prioridade do consumer                                       | ORDER, BALANCED, HIGH_PERFORMANCE, RISKY   | BALANCED               | Não          | Usa default             |
| **KAFKA_AUTO_OFFSET_RESET**        | Offset inicial                                               | earliest, latest, beginning, end, etc       | latest                 | Não          | Usa default             |

> \* Obrigatório apenas se o protocolo SASL exigir autenticação (ex: PLAIN, SCRAM, etc). Para protocolos sem autenticação (plaintext), essas variáveis são ignoradas.
>
> † A obrigatoriedade das credenciais do Schema Registry segue estas regras:
> - Se KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE="USER_INFO" (valor padrão), então KAFKA_SCHEMA_REGISTRY_USERNAME e KAFKA_SCHEMA_REGISTRY_PASSWORD são obrigatórios
> - Se KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE="SASL_INHERIT", as credenciais SASL do Kafka (KAFKA_USERNAME e KAFKA_PASSWORD) serão utilizadas, e as credenciais específicas do Schema Registry são ignoradas
> - Se o Schema Registry não requerer autenticação, use KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE="" (string vazia)
> - Se KAFKA_SCHEMA_REGISTRY_USERNAME e KAFKA_SCHEMA_REGISTRY_PASSWORD forem informados sem especificar KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE, o valor "USER_INFO" será assumido por padrão

### Detalhes e Observações
- **Obrigatórios**: KAFKA_BROKERS, KAFKA_GROUPID, KAFKA_SCHEMA_REGISTRY_URL e, dependendo do protocolo, KAFKA_USERNAME/KAFKA_PASSWORD e KAFKA_SCHEMA_REGISTRY_USERNAME/KAFKA_SCHEMA_REGISTRY_PASSWORD.
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
A biblioteca realiza verificação automática dos valores informados:

- Para o **KAFKA_AUTO_OFFSET_RESET**, os valores válidos são:
  - `earliest`, `beginning`, `smallest` (início do tópico)
  - `latest`, `end`, `largest` (fim do tópico)
  - `error` (gera erro se não houver offset inicial)

- Para o **KAFKA_SASL_MECHANISM**, os valores válidos são:
  - `PLAIN`, `SCRAM-SHA-256`, `SCRAM-SHA-512`, `GSSAPI`, `OAUTHBEARER`, `NONE`

- Para o **KAFKA_SECURITY_PROTOCOL**, os valores válidos são:
  - `plaintext`, `sasl_plaintext`, `ssl`, `sasl_ssl`

- Para o **KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE**, os valores válidos são:
  - `USER_INFO` (usa as credenciais fornecidas em KAFKA_SCHEMA_REGISTRY_USERNAME e KAFKA_SCHEMA_REGISTRY_PASSWORD)
  - `SASL_INHERIT` (herda as credenciais SASL do Kafka - KAFKA_USERNAME e KAFKA_PASSWORD)
  - `""` (string vazia para Schema Registry sem autenticação)

- Para o **KAFKA_PRODUCER_PRIORITY** e **KAFKA_CONSUMER_PRIORITY**, consulte as seções específicas abaixo para detalhes sobre cada opção.

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

err := consumer.ConsumeMessage[MyPayload](ctx, "meu-topico", enums.JsonDeserialization, handler)
if err != nil {
    panic(err)
}
```

## Formatos Suportados

### Serialização/Deserialização
- **Json**
  - enums.JsonSerialization / enums.JsonDeserialization
  - Usa JSON Schema com Schema Registry (validação, versionamento, contratos).
  - Requer Schema Registry configurado.
  - Exemplo:
    ```go
    enums.JsonSerialization // para publicar
    enums.JsonDeserialization // para consumir
    ```
- **PureJson**
  - enums.PureJsonSerialization / enums.PureJsonDeserialization
  - Serialização/deserialização JSON puro, sem Schema Registry.
  - Útil para integração simples, sistemas legados ou ambientes de desenvolvimento.
  - Exemplo:
    ```go
    enums.PureJsonSerialization // para publicar
    enums.PureJsonDeserialization // para consumir
    ```
- **Avro**
  - enums.AvroSerialization / enums.AvroDeserialization
  - Ideal para contratos rígidos, versionamento de schema e compressão eficiente.
  - Requer Schema Registry configurado.
  - Exemplo:
    ```go
    enums.AvroSerialization
    enums.AvroDeserialization
    ```
- **Protobuf**
  - enums.ProtobufSerialization / enums.ProtobufDeserialization
  - Ótimo para payloads binários, alta performance e contratos multi-linguagem.
  - Requer Schema Registry configurado.
  - Exemplo:
    ```go
    enums.ProtobufSerialization
    enums.ProtobufDeserialization
    ```

## Adaptador Protobuf

### Problema de Compatibilidade Resolvido
A biblioteca inclui um adaptador que resolve automaticamente a incompatibilidade entre diferentes implementações de protobuf:

1. A biblioteca Confluent Kafka utiliza a implementação antiga do protobuf (`github.com/golang/protobuf`)
2. Aplicações modernas geralmente usam a implementação mais recente (`google.golang.org/protobuf`)
3. Embora ambas tenham a interface `proto.Message`, elas não são diretamente compatíveis

Existe um adaptador na API que atua como uma ponte, permitindo que você use mensagens protobuf modernas sem precisar modificar seu código para acomodar a biblioteca Confluent.

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
err := consumer.ConsumeMessage(ctx, "meu-topico", enums.ProtobufDeserialization, handler)
```

### Quando o Adaptador é Necessário
Este adaptador é necessário nas seguintes situações:
1. Você está usando a implementação moderna do protobuf (`google.golang.org/protobuf`)
2. Seu sistema precisa se comunicar com o Kafka usando mensagens protobuf
3. Você recebe o erro: "serialization target must be a protobuf message"

### Prioridades Producer
- **ORDER**
  - Garante ordem de entrega e idempotência.
  - Uso: sistemas financeiros, logs ordenados, eventos críticos.
  - Pode impactar performance.
- **BALANCED** (padrão)
  - Equilíbrio entre performance e consistência.
  - Uso geral.
- **HIGH_PERFORMANCE**
  - Máximo throughput, menos garantias de ordem e duplicidade.
  - Uso: telemetria, métricas, grandes volumes de dados.

### Prioridades Consumer
- **ORDER**
  - Consumo ordenado, maior segurança e controle.
  - Uso: processamento sequencial, workflows dependentes.
- **BALANCED** (padrão)
  - Equilíbrio entre performance e consistência.
  - Uso geral.
- **HIGH_PERFORMANCE**
  - Consumo rápido, menos garantias de ordem.
  - Uso: processamento paralelo, baixa latência.
- **RISKY**
  - Máximo throughput, mínimo de garantias (auto-commit, menos consistência).
  - Uso: analytics, logs, cenários onde perder mensagens é aceitável.

### Offset
- **earliest**
  - Consome desde o início do tópico.
  - Útil para reprocessamento ou bootstrap de dados.
- **latest** (padrão)
  - Consome apenas novas mensagens.
  - Uso comum em produção.
- **beginning**
  - Sinônimo de earliest.
- **end**
  - Sinônimo de latest.
- **smallest**
  - Offset mais antigo disponível (pode variar conforme retenção do tópico).
- **largest**
  - Offset mais recente disponível.
- **error**
  - Gera erro se não houver offset inicial (útil para debugging e controle estrito).

### Protocolos de Segurança
- **plaintext**
  - Sem criptografia. Não recomendado para ambientes de produção.
- **sasl_plaintext**
  - SASL sem criptografia. Use apenas em ambientes controlados.
- **ssl**
  - Comunicação criptografada via SSL/TLS.
  - Recomendado para produção.
- **sasl_ssl**
  - SASL sobre SSL/TLS. Máxima segurança e autenticação.
  - Recomendado para ambientes sensíveis.

### SASL Mechanisms
- **PLAIN**
  - Usuário/Senha simples. Fácil de configurar, menos seguro.
- **SCRAM-SHA-256**
  - Autenticação forte baseada em hash SHA-256.
- **SCRAM-SHA-512**
  - Autenticação forte baseada em hash SHA-512.
- **GSSAPI**
  - Suporte a Kerberos (ambientes corporativos).
- **OAUTHBEARER**
  - OAuth2 para autenticação federada.
- **NONE**
  - Sem autenticação (apenas para testes locais).

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
- **JSON Puro**: `playground/json_pure_example/main.go`

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
go run playground/json_pure_example/main.go
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

Dependendo do seu ambiente, você pode precisar configurar a autenticação do Schema Registry de diferentes formas:

- **Ambiente com autenticação básica**: Use `KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=USER_INFO` e forneça as credenciais específicas
  ```env
  KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=USER_INFO
  KAFKA_SCHEMA_REGISTRY_USERNAME=usuario
  KAFKA_SCHEMA_REGISTRY_PASSWORD=senha
  ```

- **Ambiente usando as mesmas credenciais do Kafka**: Use `KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=SASL_INHERIT`
  ```env
  KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=SASL_INHERIT
  # Neste caso, as credenciais KAFKA_USERNAME e KAFKA_PASSWORD serão usadas
  ```

- **Ambiente sem autenticação no Schema Registry**: Use `KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=""` (string vazia)
  ```env
  KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE=""
  # Não é necessário fornecer KAFKA_SCHEMA_REGISTRY_USERNAME e KAFKA_SCHEMA_REGISTRY_PASSWORD
  ```

## Exemplo de Uso