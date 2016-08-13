// account
package models

type Account struct {
	Id       int64
	Account  string `xorm:"varchar(32) notnull unique"`
	Password string `xorm:"varchar(32) notnull"`
	Money    int64  `xorm:"notnull"`
	Created  int64  `xorm:"created notnull"`
}

func AccountVerifyPassword(account, password string) (string, bool, error) {
	acc := Account{
		Account: account,
	}
	ok, err := engine.Get(&acc)
	if err != nil {
		return "服务器错误", false, err
	}
	if !ok {
		return "账户不存在", false, nil
	}
	if acc.Password != password {
		return "密码错误", false, nil
	}
	return "", true, nil
}

func Register(account, password string) (string, bool, error) {
	acc := Account{
		Account: account,
	}
	ok, err := engine.Get(&acc)
	if err != nil {
		return "服务器错误", false, err
	}
	if ok {
		return "账户已存在", false, nil
	}

	acc = Account{
		Account:  account,
		Password: password,
		Money:    0,
	}
	_, err = engine.InsertOne(&acc)
	if err != nil {
		return "服务器错误", false, err
	}

	return "", true, nil
}

func ModifyPassword(account, password, npassword string) (string, bool, error) {
	acc := Account{
		Account: account,
	}
	ok, err := engine.Get(&acc)
	if err != nil {
		return "服务器错误", false, err
	}
	if !ok {
		return "账户不存在", false, nil
	}

	if acc.Password != password {
		return "原密码错误", false, nil
	}

	_, err = engine.Exec("update account set password = ? where id = ?", npassword, acc.Id)
	if err != nil {
		return "服务器错误", false, err
	}

	return "", true, nil
}
