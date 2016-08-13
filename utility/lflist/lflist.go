// lflist
package lflist

import (
	"errors"
	"fmt"
	"gilgamesh/protos"
	"gilgamesh/utility/config"
	"io/ioutil"
	"strings"
)

var (
	ErrUnknownCharacter  error = errors.New("unknown character set")
	ErrUnknownType       error = errors.New("unknown lflist type")
	ErrIDFormatIncorrect error = errors.New("card id format incorrect")
	ErrIndexOutOfRange   error = errors.New("index out of range")
)

type LFListContainer struct {
	option *config.ResourceOption
	lflist []protos.Public_LFList
}

func NewLFListContainer(option *config.ResourceOption) (*LFListContainer, error) {
	d, err := ioutil.ReadFile(option.LFList)
	if err != nil {
		return nil, err
	}

	lflist, err := loadLFList(d)
	if err != nil {
		return nil, err
	}

	return &LFListContainer{
		option: option,
		lflist: lflist,
	}, nil
}

func (c *LFListContainer) GetList() []*protos.Public_LFList {
	lflist := make([]*protos.Public_LFList, len(c.lflist))
	for i, v := range c.lflist {
		lflist[i].Id = v.Id
		lflist[i].Name = v.Name
	}
	return lflist
}

func (c *LFListContainer) GetData(i int) (*protos.Public_LFList, error) {
	if i >= len(c.lflist) {
		return nil, ErrIndexOutOfRange
	}
	return &c.lflist[i], nil
}

func loadLFList(d []byte) ([]protos.Public_LFList, error) {
	lines := strings.Split(string(d), "\n")
	for i, v := range lines {
		lines[i] = strings.Trim(v, "\r\n\t ")
	}

	lflist := []protos.Public_LFList{}
	row := -1
	target := 0
	for i, v := range lines {
		if len(v) < 6 || i == 0 {
			continue
		}
		switch v[0] {
		case '!':
			lflist = append(lflist, protos.Public_LFList{
				Id:        uint32(i),
				Name:      v[1:],
				Forbidden: []uint32{},
				Limit:     []uint32{},
				SemiLimit: []uint32{},
			})
			row++
		case '#':
			switch v[1:] {
			case "forbidden":
				target = 1
			case "limit":
				target = 2
			case "semi limit":
				target = 3
			default:
				return nil, ErrUnknownType
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			id := 0
			_, err := fmt.Sscan(strings.TrimLeft(v, "0"), &id)
			if err != nil {
				return nil, err
			}
			if id == 0 {
				return nil, ErrIDFormatIncorrect
			}
			switch target {
			case 1:
				lflist[row].Forbidden = append(lflist[row].Forbidden, uint32(id))
			case 2:
				lflist[row].Limit = append(lflist[row].Limit, uint32(id))
			case 3:
				lflist[row].SemiLimit = append(lflist[row].SemiLimit, uint32(id))
			default:
				return nil, ErrUnknownType
			}
		default:
			return nil, ErrUnknownCharacter
		}
	}

	return lflist, nil
}
