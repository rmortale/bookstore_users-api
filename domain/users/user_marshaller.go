package users

type PublicUser struct {
	Id          int64  `json:"id"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

type PrivateUser struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	DateCreated string `json:"date_created"`
	Status      string `json:"status"`
}

func (users Users) Marshall(isPublic bool) []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Marshall(isPublic)
	}
	return result
}

func (user *User) Marshall(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:          user.Id,
			DateCreated: user.DateCreated,
			Status:      user.Status,
		}
	}
	return PrivateUser{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		DateCreated: user.DateCreated,
		Status:      user.Status,
	}
}
