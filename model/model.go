package model

import (
	"context"
	"errors"
	"fmt"

	mongodb "github.com/Wondersmasher/Referral/mongoDb"
	"github.com/Wondersmasher/Referral/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Referral struct {
	ID         bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email      string        `json:"email" bson:"email" validate:"required"`
	Username   string        `json:"username" bson:"username" validate:"required"`
	Password   string        `json:"password" bson:"password" validate:"required"`
	ReferredBy string        `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
}
type User struct {
	ID              bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName       string        `json:"firstName" bson:"firstName" validate:"required"`
	LastName        string        `json:"lastName" bson:"lastName" validate:"required"`
	Email           string        `json:"email" bson:"email" validate:"required"`
	Username        string        `json:"username" bson:"username" validate:"required"`
	Password        string        `json:"password" bson:"password" validate:"required"`
	ConfirmPassword string        `json:"confirmPassword" bson:"confirmPassword,omitempty" validate:"required"`
	ReferredBy      string        `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
	ReferralID      string        `json:"referralID,omitempty" bson:"referralID,omitempty"`
	Referrals       []Referral    `json:"referrals" bson:"referrals"`
}

type Trim struct {
	Email      string        `json:"email" bson:"email" validate:"required"`
	Username   string        `json:"username" bson:"username" validate:"required"`
	ReferredBy string        `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
	Referrals  []Referral    `json:"referrals" bson:"referrals"`
	ID         bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstNam   string        `json:"firstName" bson:"firstName" validate:"required"`
	ReferralID string        `json:"referralID" bson:"referralID"`
	LastName   string        `json:"lastName" bson:"lastName" validate:"required"`
}

func (u *User) TrimUser(isSuperTrim bool) *Trim {
	if isSuperTrim {
		return &Trim{
			Email:      u.Email,
			Username:   u.Username,
			ReferredBy: u.ReferredBy,
			ReferralID: u.ReferralID,
			ID:         u.ID,
			FirstNam:   u.FirstName,
			LastName:   u.LastName,
		}
	}
	return &Trim{
		Email:      u.Email,
		Username:   u.Username,
		ReferredBy: u.ReferredBy,
		ReferralID: u.ReferralID,
		Referrals:  u.Referrals,
		ID:         u.ID,
		FirstNam:   u.FirstName,
		LastName:   u.LastName,
	}
}

func (u *User) CreateUser(referrer string) (User, error) {
	if u.Password != u.ConfirmPassword {
		return User{}, errors.New("passwords and confirm password do not match")
	}

	u.Password, _ = utils.HashPassword(u.Password)
	u.ConfirmPassword = ""
	u.Referrals = []Referral{}
	referralID, err := utils.GenerateReferralID()
	if referrer == "" {
		u.ReferredBy = referrer
	}
	if err != nil {
		return User{}, errors.New("could not generate referral id")
	}
	u.ReferralID = referralID

	user, err := mongodb.UserCollection.InsertOne(context.TODO(), u)
	if err != nil {
		return User{}, errors.New("could not create user")
	}

	u.ID = user.InsertedID.(bson.ObjectID)

	if referrer != "" {
		fmt.Println("entered here oooo!!!!!")
		filter := bson.D{{Key: "referralID", Value: referrer}}
		referral := bson.D{{Key: "$push", Value: bson.D{{Key: "referrals", Value: u.TrimUser(true)}}}}
		_, err = mongodb.UserCollection.UpdateOne(context.TODO(), filter, referral)
		if err != nil {
			return User{}, errors.New("could not update referral for" + referrer)
		}
	}
	return *u, err
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
	filter := bson.D{{Key: "referredBy", Value: id}}
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
