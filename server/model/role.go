package model

type Role struct {
	Id					int				`json:id,omitempty`
	UserId	    int				`json:userId,omitempty`
	Name				string		`json:name,omitempty`
	IsAdmin			bool			`json:isAdmin,omitempty`
	IsGod				bool			`json:IsGod,omitempty`
}
