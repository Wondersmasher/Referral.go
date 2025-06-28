package model

import (
	"context"
	"errors"
	"fmt"

	mongodb "github.com/Wondersmasher/Referral/mongoDb"
	"github.com/Wondersmasher/Referral/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Referral struct {
	ID         bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email      string        `json:"email" bson:"email" validate:"required,email"`
	Username   string        `json:"username" bson:"username" validate:"required"`
	Password   string        `json:"password" bson:"password" validate:"required"`
	ReferredBy string        `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
	IPAddress  string        `json:"ipAddress" bson:"ipAddress"`
}
type User struct {
	ID              bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	IPAddress       string        `json:"ipAddress" bson:"ipAddress"`
	FirstName       string        `json:"firstName" bson:"firstName" validate:"required"`
	LastName        string        `json:"lastName" bson:"lastName" validate:"required"`
	Email           string        `json:"email" bson:"email" validate:"required,email"`
	Username        string        `json:"username" bson:"username" validate:"required"`
	Password        string        `json:"password" bson:"password" validate:"required"`
	ConfirmPassword string        `json:"confirmPassword" bson:"confirmPassword,omitempty" validate:"required,eqfield=Password"`
	ReferredBy      string        `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
	ReferralID      string        `json:"referralID,omitempty" bson:"referralID,omitempty"`
	// Referrals       []Referral    `json:"referrals" bson:"referrals"`
}

type Trim struct {
	Email      string `json:"email" bson:"email" validate:"required,email"`
	Username   string `json:"username" bson:"username" validate:"required"`
	ReferredBy string `json:"referredBy,omitempty" bson:"referredBy,omitempty"`
	IPAddress  string `json:"ipAddress" bson:"ipAddress"`
	// Referrals  []Referral    `json:"referrals" bson:"referrals"`
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
			IPAddress:  u.IPAddress,
		}
	}
	return &Trim{
		Email:      u.Email,
		Username:   u.Username,
		ReferredBy: u.ReferredBy,
		ReferralID: u.ReferralID,
		// Referrals:  u.Referrals,
		ID:        u.ID,
		FirstNam:  u.FirstName,
		LastName:  u.LastName,
		IPAddress: u.IPAddress,
	}
}

func (u *User) CreateUser(referrer string) (User, error) {
	if u.Password != u.ConfirmPassword {
		return User{}, errors.New("passwords and confirm password do not match")
	}

	u.Password, _ = utils.HashPassword(u.Password)
	u.ConfirmPassword = ""
	// u.Referrals = []Referral{}
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

	// this is TODO:
	// 	u.IPAddress = c.ClientIP()
	// isValid := model.IsIPValid(u.IPAddress, u.ReferralID)
	// if !isValid {
	// 	c.JSON(400, utils.ApiErrorResponse("referer and referee cannot use same IP address"))
	// 	return
	// }
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

func GetReferrals(id string) ([]Trim, error) {
	filter := bson.D{{Key: "referredBy", Value: id}}
	projection := bson.D{
		{Key: "email", Value: 1},
		{Key: "username", Value: 1},
		{Key: "referredBy", Value: 1},
		{Key: "referralID", Value: 1},
		// {Key: "referrals", Value: 1},
		{Key: "_id", Value: 1},
		{Key: "firstName", Value: 1},
		{Key: "lastName", Value: 1},
	}
	opts := options.Find().SetProjection(projection)
	cursor, err := mongodb.UserCollection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var trims []Trim
	for cursor.Next(context.TODO()) {
		var t Trim
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		trims = append(trims, t)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return trims, nil
}

func IsIPUsed(IP string) bool {
	u := User{}
	filter := bson.D{{Key: "ipAddress", Value: IP}}

	err := mongodb.UserCollection.FindOne(context.TODO(), filter).Decode(&u)
	return err == nil
}
