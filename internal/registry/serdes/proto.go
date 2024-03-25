package serdes

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"

	"github.com/irononet/mallbots/internal/registry"
)

type ProtoSerde struct {
	r registry.Registry
}

var _ registry.Serde = (*ProtoSerde)(nil)

var protoT = reflect.TypeOf((*proto.Message)(nil)).Elem()

func (c ProtoSerde) Register(v registry.Registrable, options ...registry.BuildOption) error {
	if !reflect.TypeOf(v).Implements(protoT) {
		return fmt.Errorf("%T does not implement proto.Message")
	}
	return registry.Register(c.r, v, c.serialize, c.deserialize, options)
}

func (c ProtoSerde) RegisterKey(key string, v any, options ...registry.BuildOption) error {
	if !reflect.TypeOf(v).Implements(protoT) {
		return fmt.Errorf("%T does not implement proto.Message", v)
	}
	return registry.RegisterKey(c.r, key, v, c.serialize, c.deserialize, options)
}

func (c ProtoSerde) RegisterFactory(key string, fn func() any, options ...registry.BuildOption) error {
	if v := fn(); v == nil {
		return fmt.Errorf("%s factory returns a nil value", key)
	} else if _, ok := v.(proto.Message); !ok {
		return fmt.Errorf("%s does not implement proto.Message", key)
	}
	return registry.RegisterFactory(c.r, key, fn, c.serialize, c.deserialize, options)
}


func (c ProtoSerde) serialize(v any) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (c ProtoSerde) deserialize(data []byte, value any) error {
	return proto.Unmarshal(data, value.(proto.Message))
}
