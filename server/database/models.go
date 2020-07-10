package database

// All database models
var Models = []interface{}{
	&User{},
	&Role{},
}

type User struct {
	Id					int					`sql:",pk"`
	Name				string
	Email				string
	Password		[]uint8
	Role				*Role
	IsFacebook	bool
	IsGoogle		bool
}

type Role struct {
	Id					int
	UserId	    int
	Name				string
	IsAdmin			bool
	IsGod				bool
}
