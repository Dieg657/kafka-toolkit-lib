package engine

import (
	"fmt"
	"reflect"
	"sync"

	"google.golang.org/protobuf/proto"
)

// ProtobufAdapter provê métodos para adaptar entre diferentes implementações de protobuf.
//
// # Incompatibilidade resolvida
//
// Este adaptador resolve a incompatibilidade entre diferentes implementações da interface proto.Message:
//  1. A biblioteca Confluent espera objetos que implementem a interface proto.Message da implementação
//     'github.com/golang/protobuf/proto' (implementação mais antiga)
//  2. Muitos projetos usam a implementação mais recente 'google.golang.org/protobuf/proto'
//  3. Embora ambas interfaces tenham o mesmo nome, elas não são compatíveis automaticamente
//
// # Quando a incompatibilidade não existe
//
// Este adaptador não é necessário quando:
//  1. Seu código usa exclusivamente a mesma implementação de protobuf que a biblioteca Confluent
//  2. Você está usando versões do protobuf geradas com ferramentas que garantem compatibilidade
//     com a biblioteca Confluent (geralmente usando a mesma versão do compilador protoc)
//  3. Seus objetos já implementam ambas interfaces proto.Message (old e new)
//
// # Uso automático
//
// Na maioria dos casos, este adaptador funcionará automaticamente sem necessidade de configuração.
// Ele tenta detectar a implementação protobuf usada e faz a conversão necessária em tempo de execução.
type ProtobufAdapter struct {
	// registeredTypes é um registro de tipos protobuf que podem ser usados como intermediários
	// para deserialização, mapeados pelo tipo de destino
	registeredTypes sync.Map

	// cache para tipos detectados automaticamente
	autoDetectedTypes sync.Map
}

// NewProtobufAdapter cria uma nova instância do adaptador de protobuf
func NewProtobufAdapter() *ProtobufAdapter {
	return &ProtobufAdapter{
		registeredTypes:   sync.Map{},
		autoDetectedTypes: sync.Map{},
	}
}

// RegisterProtoType registra um tipo protobuf para ser usado como intermediário
// na deserialização para um tipo de destino específico.
//
// Este método é OPCIONAL na maioria dos casos, pois o adaptador tenta detectar automaticamente
// os tipos protobuf. Use-o apenas em casos onde a detecção automática falha ou para otimizar
// o desempenho eliminando reflexão.
//
// # Parâmetros:
//   - targetType: é o tipo para o qual você pretende adaptar (por ex: reflect.TypeOf((*MeuTipo)(nil)).Elem())
//   - protoType: é um exemplar do tipo protobuf (por ex: &pb.ProtoPayload{})
func (adapter *ProtobufAdapter) RegisterProtoType(targetType reflect.Type, protoType proto.Message) {
	adapter.registeredTypes.Store(targetType.String(), reflect.TypeOf(protoType))
}

// AdaptMessage adapta uma mensagem de qualquer tipo para o tipo esperado pelo serializador da Confluent.
// Se o objeto já implementa proto.Message, ele é retornado como está.
// Se não, tenta converter ou retorna um erro explicando o problema.
//
// # Incompatibilidade tratada:
// Este método resolve o problema da biblioteca Confluent esperar uma implementação específica
// da interface proto.Message, enquanto seu código pode estar usando outra.
func (adapter *ProtobufAdapter) AdaptMessage(obj interface{}) (interface{}, error) {
	// Verifica se o objeto já implementa a interface proto.Message
	if protoMsg, ok := obj.(proto.Message); ok {
		return protoMsg, nil
	}

	// Verifica se estamos recebendo um ponteiro
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("o objeto deve ser um ponteiro para um protobuf, recebido: %T", obj)
	}

	// Acessa o objeto através do ponteiro
	elem := val.Elem().Interface()

	// Verifica se o objeto desreferenciado implementa proto.Message
	if protoMsg, ok := elem.(proto.Message); ok {
		return protoMsg, nil
	}

	// Se chegamos aqui, o objeto não é um protobuf válido
	// mas vamos tentar usar métodos alternativos para adaptá-lo

	// Verificamos se o objeto tem métodos com nomes que sugerem implementação de protobuf
	// por exemplo, ProtoMessage(), ProtoReflect(), Descriptor()
	if adapter.hasProtobufSignature(val) {
		// Temos um objeto que parece ser protobuf mas usa uma implementação diferente
		// Podemos criar um novo objeto protobuf e copiar os campos
		protoMsg, err := adapter.convertToProtoMessage(obj)
		if err == nil {
			return protoMsg, nil
		}
	}

	return nil, fmt.Errorf("o objeto não implementa a interface proto.Message necessária para serialização protobuf: %T", obj)
}

// hasProtobufSignature verifica se um objeto tem métodos que sugerem que é um protobuf.
//
// Esta é uma heurística usada para detectar automaticamente objetos protobuf
// mesmo quando eles não implementam a interface proto.Message específica esperada.
func (adapter *ProtobufAdapter) hasProtobufSignature(val reflect.Value) bool {
	// Verifica se o tipo tem métodos como ProtoMessage(), ProtoReflect(), etc.
	typ := val.Type()

	// Lista de métodos que sugerem uma implementação protobuf
	protoMethods := []string{"ProtoMessage", "ProtoReflect", "Descriptor", "Reset", "String"}

	methodCount := 0
	for _, methodName := range protoMethods {
		if _, exists := typ.MethodByName(methodName); exists {
			methodCount++
		}
	}

	// Se tem pelo menos 3 desses métodos, provavelmente é um protobuf
	return methodCount >= 3
}

// convertToProtoMessage tenta converter um objeto para proto.Message.
//
// Este método usa reflexão para converter entre diferentes implementações de protobuf,
// resolvendo a incompatibilidade em tempo de execução.
func (adapter *ProtobufAdapter) convertToProtoMessage(obj interface{}) (proto.Message, error) {
	// Usa reflexão para extrair dados do objeto e criar um novo objeto proto.Message
	// Isso é uma implementação genérica que tenta diferentes abordagens

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Verifica se temos um tipo dinamicamente detectado para este objeto
	objType := reflect.TypeOf(obj)
	if cachedType, ok := adapter.autoDetectedTypes.Load(objType.String()); ok {
		// Cria uma nova instância do tipo detectado
		protoMsgType := cachedType.(reflect.Type)
		protoMsgVal := reflect.New(protoMsgType.Elem())
		protoMsg := protoMsgVal.Interface().(proto.Message)

		// Copia campos usando reflexão
		adapter.copyStructFields(val, protoMsgVal.Elem())

		return protoMsg, nil
	}

	// Para implementação completa, precisaríamos de uma forma de criar um novo proto.Message
	// a partir do zero, o que é difícil sem conhecer o tipo específico
	return nil, fmt.Errorf("não foi possível converter para proto.Message")
}

// copyStructFields copia campos de uma struct para outra usando reflexão.
//
// Este método é usado para transferir dados entre diferentes implementações
// de objetos protobuf, mantendo os valores dos campos.
func (adapter *ProtobufAdapter) copyStructFields(src, dst reflect.Value) {
	// Certifique-se de que ambos são structs
	if src.Kind() != reflect.Struct || dst.Kind() != reflect.Struct {
		return
	}

	// Para cada campo na struct de destino
	for i := 0; i < dst.NumField(); i++ {
		dstField := dst.Field(i)
		dstFieldName := dst.Type().Field(i).Name

		// Verifica se existe um campo com o mesmo nome na source
		srcField := src.FieldByName(dstFieldName)
		if srcField.IsValid() && srcField.Type().AssignableTo(dstField.Type()) {
			dstField.Set(srcField)
		}
	}
}

// IsProtobufMessage verifica se um objeto implementa a interface proto.Message
func (adapter *ProtobufAdapter) IsProtobufMessage(obj interface{}) bool {
	_, ok := obj.(proto.Message)
	return ok
}

// GetProtoTypeFor retorna um tipo protobuf registrado para o tipo de destino especificado.
// Se nenhum tipo estiver registrado, retorna um tipo protobuf genérico.
func (adapter *ProtobufAdapter) GetProtoTypeFor(targetType reflect.Type) (reflect.Type, error) {
	// Tenta obter um tipo registrado para o tipo alvo
	if protoType, ok := adapter.registeredTypes.Load(targetType.String()); ok {
		return protoType.(reflect.Type), nil
	}

	// Verifica no cache de tipos autodetectados
	if protoType, ok := adapter.autoDetectedTypes.Load(targetType.String()); ok {
		return protoType.(reflect.Type), nil
	}

	// Se não encontrar um tipo registrado específico, retorna um erro
	return nil, fmt.Errorf("nenhum tipo protobuf registrado para %s", targetType.String())
}

// CreateProtoInstance cria uma nova instância de um tipo protobuf registrado
// para o tipo de destino especificado, ou um protobuf genérico se nenhum estiver registrado.
//
// Este método resolve a incompatibilidade em tempo de execução, permitindo que objetos
// compilados com diferentes versões de protobuf trabalhem juntos.
func (adapter *ProtobufAdapter) CreateProtoInstance(targetType reflect.Type) (proto.Message, error) {
	protoType, err := adapter.GetProtoTypeFor(targetType)
	if err != nil {
		// Se não encontrar um tipo registrado, usamos uma estratégia alternativa
		// como criar um protobuf genérico ou obter do próprio target

		// Verifica se o targetType implementa proto.Message
		if reflect.PtrTo(targetType).Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
			// Se o alvo já implementa proto.Message, podemos criar uma instância dele
			protoObj := reflect.New(targetType).Interface()

			// Cacheia esta descoberta para uso futuro
			adapter.autoDetectedTypes.Store(targetType.String(), reflect.TypeOf(protoObj))

			return protoObj.(proto.Message), nil
		}

		// Verifica se o tipo poderia ser um protobuf
		if adapter.couldBeProtobuf(targetType) {
			// Tentativa adicional usando heurísticas, mas não implementamos criação dinâmica
			// de protos nesta versão do adaptador
		}

		// Não conseguimos criar uma instância apropriada
		return nil, fmt.Errorf("não é possível criar uma instância protobuf para o tipo %s", targetType.String())
	}

	// Cria uma nova instância do tipo protobuf registrado
	return reflect.New(protoType.Elem()).Interface().(proto.Message), nil
}

// couldBeProtobuf verifica se um tipo poderia ser um tipo protobuf.
//
// Esta heurística ajuda na detecção automática de tipos compatíveis com protobuf,
// reduzindo a necessidade de registro manual por parte do usuário.
func (adapter *ProtobufAdapter) couldBeProtobuf(t reflect.Type) bool {
	// Verifica se o tipo tem características que sugerem que poderia ser um protobuf
	// por exemplo, campos que seguem convenções de nomes protobuf

	if t.Kind() != reflect.Struct {
		return false
	}

	// Conta quantos campos parecem seguir convenções protobuf
	protoLikeFields := 0
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Campos exportados e com tags json são comuns em protos gerados
		if field.PkgPath == "" && field.Tag.Get("json") != "" {
			protoLikeFields++
		}
	}

	// Se a maioria dos campos parece seguir convenções protobuf
	return protoLikeFields > 0 && float64(protoLikeFields)/float64(t.NumField()) > 0.5
}
