package repositories

import (
	"context"
	"errors"
	"mindbridge/domain/entities"
	domainRepo "mindbridge/domain/repositories"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
}

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func isValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func NewUserRepository(collection *mongo.Collection) domainRepo.IUserRepository {
	return &UserRepository{
		collection: collection,
	}
}

func (r *UserRepository) Create(user *entities.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = ""
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	if !isValidEmail(email) {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(username string) (*entities.User, error) {
	if !isValidUsername(username) {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user entities.User
	opts := options.FindOne().SetProjection(bson.M{"password": 0})
	err := r.collection.FindOne(ctx, bson.M{"username": username}, opts).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		filter["_id"] = objID // most common case
	} else {
		filter["_id"] = id // fallback for non‑ObjectID ids
	}

	type mongoUser struct {
		ID        interface{}        `bson:"_id"` // could be ObjectID or string
		Name      string             `bson:"name,omitempty"`
		Username  string             `bson:"username,omitempty"`
		Email     string             `bson:"email,omitempty"`
		CreatedAt primitive.DateTime `bson:"created_at,omitempty"`
		UpdatedAt primitive.DateTime `bson:"updated_at,omitempty"`
	}
	var mu mongoUser

	opts := options.FindOne().SetProjection(bson.M{"password": 0}) // never return the hash
	if err := r.collection.FindOne(ctx, filter, opts).Decode(&mu); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	user := &entities.User{
		Name:      mu.Name,
		Username:  mu.Username,
		Email:     mu.Email,
		CreatedAt: mu.CreatedAt.Time(),
		UpdatedAt: mu.UpdatedAt.Time(),
	}

	switch v := mu.ID.(type) {
	case primitive.ObjectID:
		user.ID = v.Hex()
	case string:
		user.ID = v
	}

	return user, nil
}
