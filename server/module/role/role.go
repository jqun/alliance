package role

type Role struct {
	Name         string
	AllianceName string // 公会Id
}

type User struct {
	Role *Role
}

func NewUser(roleName string) *User {
	m := &User{}
	m.Role = &Role{
		Name: roleName,
	}
	return m
}

func (u *User) SetAllianceName(allianceName string) {
	u.Role.AllianceName = allianceName
}