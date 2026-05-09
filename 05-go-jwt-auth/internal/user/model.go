package user

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	
	Name string `bson:"name" json:"name"`

	Email string `bson:"email" json:"email"`

	Password string `bson:"password" json:"-"`

	Role string `bson:"role" json:"role"`

	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`

	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

type PublicUser struct {
	ID string `json:"_id"`
	
	Name string `json:"name"`

	Email string `json:"email"`

	Role string `json:"role"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`
}

func ToPublic(u User) PublicUser {
	return PublicUser{
		ID: u.ID.Hex(),
		Name: u.Name,
		Email: u.Email,
		Role: u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}


