package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo struct {
	coll *mongo.Collection
}


func NewRepo(db *mongo.Database) *Repo {
	return &Repo {
		coll: db.Collection("user"),
	}
}

func (r *Repo) FindByEmail(ctx context.Context, email string) (User, error) {
	// clean email
	email = strings.ToLower(strings.TrimSpace(email))

	filter := bson.M{
		"email": email,
	}

	var u User;
	err := r.coll.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return User{}, mongo.ErrNoDocuments
		}

		return User{}, fmt.Errorf("Find user by email Failed: %w", err)
	}

	return u, nil
}

func (r *Repo) Create(ctx context.Context, u User ) (User, error) {

	res, err := r.coll.InsertOne(ctx, u);
	if err != nil {
		return User{}, fmt.Errorf("Failed to insert user to database: %w", err);
	}

	id, ok := res.InsertedID.(bson.ObjectID);
	if !ok {
		return User{}, fmt.Errorf("Error from InsertOne ObjectID");
	}

	u.ID = id

	return u, nil
}