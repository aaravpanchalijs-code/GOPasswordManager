package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Vault struct {
	ID         bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     bson.ObjectID `json:"user_id" bson:"user_id"`
	Website    string        `json:"website" bson:"website"`
	LoginEmail string        `json:"login_email" bson:"login_email"`
	Password   string        `json:"password" bson:"password"`
	Notes      string        `json:"notes" bson:"notes"`
	CreatedAt  time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at" bson:"updated_at"`
}
