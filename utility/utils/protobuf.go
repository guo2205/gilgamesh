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
		return nil, "", ErrUnknownProtobufType
	}
	val := reflect.New(t)
	obj, ok := val.Interface().(proto.Message)
	if !ok {
		return nil, "", ErrUnknownProtobufType
	}

	err = proto.Unmarshal(f.Data, obj)
	if err != nil {
		return nil, "", err
	}

	return obj, f.Type, nil
}
