package entities

//UserInfoAtomicService .
type UserInfoAtomicService struct{}

//UserInfoService .
var UserInfoService = UserInfoAtomicService{}

// Save .
func (*UserInfoAtomicService) Save(u *UserInfo) error {
	_, err := engine.Insert(u)
	return err
}

// FindAll .
func (*UserInfoAtomicService) FindAll() []UserInfo {
	everyone := make([]UserInfo, 0)
	err := engine.Find(&everyone)
	checkErr(err)
	return everyone
}

// FindByID .
func (*UserInfoAtomicService) FindByID(id int) *UserInfo {
	u := new(UserInfo)
	u.UID = id
	has, err := engine.Get(u)
	checkErr(err)
	if !has {
		u.UserName = "unknown"
	}
	return u
}
