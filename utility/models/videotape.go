// videotape
package models

type Videotape struct {
	Id      int64
	Account string `xorm:"varchar(32) notnull"`
	Name    string `xorm:"varchar(32) notnull"`
	Hash    string `xorm:"varchar(32) notnull"`
	Created int64  `xorm:"created notnull"`
}

func CreateVideotape(account, name, hashName string) error {
	vt := Videotape{
		Account: account,
		Name:    name,
		Hash:    hashName,
	}
	_, err := engine.InsertOne(&vt)
	if err != nil {
		return err
	}

	return nil
}

func QueryVideotapeList(account string) ([]*Videotape, error) {
	vtList := []*Videotape{}
	vt := Videotape{
		Account: account,
	}

	err := engine.Find(&vtList, &vt)
	if err != nil {
		return nil, err
	}

	return vtList, nil
}
