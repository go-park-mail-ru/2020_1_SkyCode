package models

type User struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	Avatar    string `json:"profile_photo"`
	Role      string `json:"role"`
}

func (u *User) IsAdmin() bool {
	return u.Role == "Admin"
}

func (u *User) IsManager() bool {
	return u.Role == "Moderator"
}

func (u *User) IsUser() bool {
	return u.Role == "User"
}

func (u *User) IsSupport() bool {
	return u.Role == "Support"
}

func (u *User) SetAdmin() {
	u.Role = "Admin"
}

func (u *User) SetManager() {
	u.Role = "Moderator"
}

func (u *User) SetUser() {
	u.Role = "User"
}
