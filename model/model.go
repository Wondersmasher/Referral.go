package model

import (
	"context"
	"errors"

	mongodb "github.com/Wondersmasher/Referral/mongoDb"
	"github.com/Wondersmasher/Referral/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type Referral struct {
	Email      string `json:"email" bson:"email" validate:"required"`
	Username   string `json:"username" bson:"username" validate:"required"`
	Password   string `json:"password" bson:"password" validate:"required"`
	ReferredBy string `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
}
type User struct {
	Email      string     `json:"email" bson:"email" validate:"required"`
	Username   string     `json:"username" bson:"username" validate:"required"`
	Password   string     `json:"password" bson:"password" validate:"required"`
	ReferredBy string     `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
	Referrals  []Referral `json:"referrals,omitempty" bson:"referrals,omitempty"`
}

type Trim struct {
	Email      string     `json:"email" bson:"email" validate:"required"`
	Username   string     `json:"username" bson:"username" validate:"required"`
	ReferredBy string     `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
	Referrals  []Referral `json:"referrals,omitempty" bson:"referrals,omitempty"`
}

func (u *User) TrimUser() *Trim {
	return &Trim{
		Email:      u.Email,
		Username:   u.Username,
		ReferredBy: u.ReferredBy,
		Referrals:  u.Referrals,
	}
}

func (u *User) CreateUser() (any, error) {
	user, err := mongodb.UserCollection.InsertOne(context.TODO(), u)
	if err != nil {
		return User{}, err
	}
	return user.InsertedID, err
}

func GetUserByEmail(email, password string) (*User, error) {
	u := User{}
	filter := bson.D{{Key: "email", Value: email}}

	err := mongodb.UserCollection.FindOne(context.TODO(), filter).Decode(&u)
	if err != nil {
		return &User{}, errors.New("invalid credentials")
	}

	isValid := utils.ValidateHashedPassword(u.Password, password)
	if !isValid {
		return &User{}, errors.New("invalid credentials")
	}

	return &u, nil
}

func GetReferrals(id string) ([]User, error) {
	filter := bson.D{{Key: "id", Value: id}}
	// Retrieves documents that match the query filter
	cursor, err := mongodb.UserCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	u := []User{}
	if err = cursor.All(context.TODO(), &u); err != nil {
		return nil, err
	}
	return u, nil
}
