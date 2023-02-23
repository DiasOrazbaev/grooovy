package repo

import (
	"context"
	"errors"
	"grovo/internal/common"
	"grovo/internal/entity"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongo struct {
	database   *mongo.Database
	collection string
}

func NewMongoUserRepository(database *mongo.Database) *UserMongo {
	return &UserMongo{
		database:   database,
		collection: "users",
	}
}

func (r *UserMongo) Login(ctx context.Context, username string) (*entity.User, error) {
	user := new(entity.User)
	w := r.database.Collection(r.collection).FindOne(ctx, bson.M{"username": username})
	if errors.Is(w.Err(), mongo.ErrNoDocuments) {
		return nil, common.ErrInvalidCredentials
	}
	err := w.Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserMongo) Register(ctx context.Context, username, password string) (*entity.User, error) {
	user := &entity.User{
		Id:       uuid.New().URN()[len("urn:uuid:"):],
		Username: username,
		Password: password,
	}

	_, err := r.database.Collection(r.collection).InsertOne(ctx, bson.M{"_id": user.Id, "username": user.Username, "password": user.Password})
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserMongo) IsUnique(ctx context.Context, username string) (bool, error) {
	w := r.database.Collection(r.collection).FindOne(ctx, bson.M{"username": username})
	if errors.Is(w.Err(), mongo.ErrNoDocuments) {
		return true, nil
	} else if w.Err() != nil {
		return false, w.Err()
	}
	return false, nil
}
