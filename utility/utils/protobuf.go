// protobuf
package utils

import (
	"errors"
	"gilgamesh/protos"
	"reflect"

	"github.com/golang/protobuf/proto"
)

var (
	ErrUnknownProtobufType error = errors.New("unknown protobuf type")
)

func Marshal(m proto.Message) []byte {
	d, err := proto.Marshal(m)
	if err != nil {
		panic(err)
	}

	f := protos.Gilgamesh{
		Type: proto.MessageName(m),
		Data: d,
	}
	d, err = proto.Marshal(&f)
	if err != nil {
		panic(err)
	}

	return d
}

func Unmarshal(d []byte) (proto.Message, string, error) {
	f := protos.Gilgamesh{}
	err := proto.Unmarshal(d, &f)
	if err != nil {
		return nil, "", err
	}

	t := proto.MessageType(f.Type)
	if t == nil {
		return nil, f.Type, ErrUnknownProtobufType
	}
	val := reflect.New(t.Elem())
	obj, ok := val.Interface().(proto.Message)
	if !ok {
		return nil, f.Type, ErrUnknownProtobufType
	}

	err = proto.Unmarshal(f.Data, obj)
	if err != nil {
		return nil, f.Type, err
	}

	return obj, f.Type, nil
}

func UnmarshalWithType(d []byte, m proto.Message) (proto.Message, string, error) {
	f := protos.Gilgamesh{}
	err := proto.Unmarshal(d, &f)
	if err != nil {
		return nil, "", err
	}

	if f.Type != proto.MessageName(m) {
		return nil, f.Type, ErrUnknownProtobufType
	}

	err = proto.Unmarshal(f.Data, m)
	if err != nil {
		return nil, f.Type, err
	}

	return m, f.Type, nil
}
