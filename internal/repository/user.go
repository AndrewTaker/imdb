package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type PaginationOptions struct {
	Limit, Offset int64
}

type UserSortOptions struct {
	SortBy SortableVar
	Order  OrderVar
}

type SortableVar string

var CreatedAt SortableVar = "created_at"
var ID SortableVar = "_id"

type OrderVar int8

var ASC OrderVar = 1
var DESC OrderVar = -1

func CreateOptions(pag PaginationOptions, sort []UserSortOptions) *options.FindOptions {
	opts := options.Find()
	opts.SetLimit(pag.Limit)
	opts.SetSkip(pag.Offset)
	for _, s := range sort {
		opts.SetSort(bson.D{{Key: string(s.SortBy), Value: s.Order}})
	}

	return opts
}

func (r *UserRepository) GetAll(ctx context.Context, pag PaginationOptions, sort []UserSortOptions) ([]User, error) {
	opts := CreateOptions(pag, sort)
	rows, err := r.c.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer rows.Close(ctx)

	users := make([]User, 0)
	if err := rows.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.c.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	return &u, err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %w", err)
	}

	var u User
	err = r.c.FindOne(ctx, bson.M{"_id": objID}).Decode(&u)

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
