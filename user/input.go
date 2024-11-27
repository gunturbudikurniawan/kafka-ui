package user

type RegisterUsersInput struct {
	Users []RegisterUserInput `json:"users"`
}

type RegisterUserInput struct {
	Name       string `json:"name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}
