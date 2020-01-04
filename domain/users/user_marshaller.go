package users

import "encoding/json"

type PublicUser struct {
	Id        int64  `json:"user_id"`
	CreatedAt string `json:"date_registered"`
	Status    string `json:"user_status"`
}

type PrivateUser struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
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
		// First Way --> if we have custom JSON field (e.g: Id int64 `json:"user_id"`)
		return PublicUser{
			Id:        user.Id,
			CreatedAt: user.CreatedAt,
			Status:    user.Status,
		}
	}

	// Second Way --> if we have same JSON field (e.g: Id int64 `json:"id"`)
	userJson, _ := json.Marshal(user)
	var privateUser PrivateUser
	json.Unmarshal(userJson, &privateUser)
	return privateUser
}
