# Kafka Toolkit Library - Playground

Este playground permite testar a biblioteca Kafka Toolkit em um ambiente local com Kafka, Zookeeper, Schema Registry e Kafka UI.

## Ambientes disponíveis

O playground oferece duas configurações:

1. **Ambiente simples (sem SASL)** - `docker-compose.yaml`
   - Sem autenticação SASL (usa PLAINTEXT)
   - Ideal para testes rápidos e demonstrações básicas
   - Tópicos são criados automaticamente

2. **Ambiente com SASL** - `docker-compose.sasl.yaml`
   - Com autenticação SASL/PLAIN
   - Similar a ambientes corporativos/produtivos
   - Demonstra o uso completo da biblioteca com autenticação

## Configuração dos ambientes

### Ambiente sem autenticação (PLAINTEXT)

#### 1. Iniciando os containers

```sh
docker compose -f playground/docker-compose.yaml up -d
```

Isto iniciará:
- Zookeeper
- Kafka (porta 29093)
- Schema Registry (porta 8081)
- Kafka UI (porta 8080)

O Kafka está configurado para criar tópicos automaticamente quando são acessados pela primeira vez.

#### 2. Variáveis de ambiente para este ambiente

```sh
# Para Unix/Linux/MacOS
export KAFKA_BROKERS=localhost:29093
export KAFKA_SCHEMA_REGISTRY_URL=http://localhost:8081
export KAFKA_SECURITY_PROTOCOL=plaintext
export KAFKA_SASL_MECHANISM=PLAIN
export KAFKA_SCHEMA_REGISTRY_USERNAME=schema-registry
export KAFKA_SCHEMA_REGISTRY_PASSWORD=schema-registry-password

# Para Windows (PowerShell)
$env:KAFKA_BROKERS="localhost:29093"
$env:KAFKA_SCHEMA_REGISTRY_URL="http://localhost:8081"
$env:KAFKA_SECURITY_PROTOCOL="plaintext"
$env:KAFKA_SASL_MECHANISM="PLAIN"
$env:KAFKA_SCHEMA_REGISTRY_USERNAME="schema-registry"
$env:KAFKA_SCHEMA_REGISTRY_PASSWORD="schema-registry-password"
```

### Ambiente com autenticação (SASL)

#### 1. Iniciando os containers

```sh
docker compose -f playground/docker-compose.sasl.yaml up -d
```

Isto iniciará:
- Zookeeper
- Kafka (porta 29093)
- Schema Registry (porta 8081)
- Kafka UI (porta 8080)

**Nota**: Para este ambiente, você precisará criar os tópicos manualmente através da Kafka UI.

#### 2. Variáveis de ambiente para este ambiente

```sh
# Para Unix/Linux/MacOS
export KAFKA_BROKERS=localhost:29093
export KAFKA_SCHEMA_REGISTRY_URL=http://localhost:8081
export KAFKA_SECURITY_PROTOCOL=sasl_plaintext
export KAFKA_SASL_MECHANISM=PLAIN
export KAFKA_USERNAME=admin
export KAFKA_PASSWORD=admin-secret
export KAFKA_SCHEMA_REGISTRY_USERNAME=schema-registry
export KAFKA_SCHEMA_REGISTRY_PASSWORD=schema-registry-password

# Para Windows (PowerShell)
$env:KAFKA_BROKERS="localhost:29093"
$env:KAFKA_SCHEMA_REGISTRY_URL="http://localhost:8081"
$env:KAFKA_SECURITY_PROTOCOL="sasl_plaintext"
$env:KAFKA_SASL_MECHANISM="PLAIN"
$env:KAFKA_USERNAME="admin"
$env:KAFKA_PASSWORD="admin-secret"
$env:KAFKA_SCHEMA_REGISTRY_USERNAME="schema-registry"
$env:KAFKA_SCHEMA_REGISTRY_PASSWORD="schema-registry-password"
```

## Tópicos utilizados pelos exemplos

Os seguintes tópicos são utilizados pelos exemplos:
- `avro-topic`
- `protobuf-topic`
- `json-schema-topic`
- `json-pure-topic`

Você pode verificar os tópicos acessando a interface Kafka UI em [http://localhost:8080](http://localhost:8080).

## Exemplos disponíveis

O playground contém exemplos para os seguintes formatos de serialização:

- **Avro**: `playground/avro_example/main.go` (requer Schema Registry)
- **Protobuf**: `playground/protobuf_example/main.go` (requer Schema Registry)
- **JSON Schema**: `playground/json_schema_example/main.go` (requer Schema Registry)
- **JSON Puro**: `playground/json_pure_example/main.go`

**Importante**: Os formatos Avro, Protobuf e JSON Schema **requerem** o Schema Registry em funcionamento, pois utilizam-no para registro e validação de schemas. Apenas o formato JSON Puro pode funcionar sem Schema Registry.

## Executando os exemplos

Para executar um exemplo, escolha o ambiente desejado (com ou sem SASL), configure as variáveis de ambiente conforme instruções acima e execute:

```sh
go run playground/avro_example/main.go
```

Substitua `avro_example` pelo exemplo que deseja executar.

## Interface Kafka UI

O ambiente inclui uma interface web para gerenciar o Kafka. Acesse:
[http://localhost:8080](http://localhost:8080)

## Parando o ambiente

Para parar o ambiente sem SASL:

```sh
docker compose -f playground/docker-compose.yaml down -v
```

Para o ambiente com SASL:

```sh
docker compose -f playground/docker-compose.sasl.yaml down -v
``` 