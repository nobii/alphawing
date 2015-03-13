package permission

type Permission interface {
	CreateGroup(name string) (groupId string, err error)
	DeleteGroup(groupId string) (groupId string, err error)
	AddUser(groupId string, email string) (permId string, err error)
	DeleteUser(groupId string, permId string) error
	GetUserList(groupId string, name string) (emails []string, err error)
}