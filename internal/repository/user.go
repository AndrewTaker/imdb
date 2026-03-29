package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
}

type UserRepository struct {
	c *mongo.Collection
}

func NewUserRepository(c *mongo.Collection) *UserRepository {
	return &UserRepository{c: c}
}

func (r *UserRepository) GetAll(ctx context.Context, pag PaginationOptions, sort []SortOptions) ([]User, error) {
	opts := CreateQueryOptions(pag, sort, nil)
	rows, err := r.c.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}
	defer rows.Close(ctx)

	users := make([]User, 0)
	for rows.Next(ctx) {
		var user User
		if err := rows.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.c.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	return &u, err
}

func (r *UserRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*User, error) {
	var u User
	err := r.c.FindOne(ctx, bson.M{"_id": id}).Decode(&u)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string) error {
	if _, err := r.c.InsertOne(ctx, User{Email: email, PasswordHash: passwordHash}); err != nil {
		return err
	}
	return nil
}
