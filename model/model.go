package model

import (
	"context"
	"errors"

	mongodb "github.com/Wondersmasher/Referral/mongoDb"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Email      string `json:"email" bson:"email" validate:"required"`
	Username   string `json:"username" bson:"username" validate:"required"`
	Password   string `json:"password" bson:"password" validate:"required"`
	ReferredBy string `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
	Referrals  []User `json:"referrals,omitempty" bson:"referrals,omitempty"`
}

func GetUserByEmail(email, password string) (*User, error) {
	u := User{}
	filter := bson.D{{Key: "email", Value: email}}

	err := mongodb.UserCollection.FindOne(context.TODO(), filter).Decode(&u)

	if err != nil {
		return &User{}, errors.New("invalid credentials")
	}

	return &u, nil
}
