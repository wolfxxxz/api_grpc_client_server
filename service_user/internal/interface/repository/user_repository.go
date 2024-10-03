package repository

import (
	"context"
	"service_user/internal/apperrors"
	"service_user/internal/domain/models"
	"service_user/internal/usecase/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
)

const userCollection = "users"

type userRepository struct {
	log        *logrus.Logger
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, log *logrus.Logger) repository.UserRepository {
	return &userRepository{collection: db.Collection(userCollection), log: log}
}

func (ur *userRepository) CreateUser(ctx context.Context, u *models.User) (string, error) {
	res, err := ur.collection.InsertOne(ctx, u)
	if err != nil {
		ur.log.Errorf("Mongo Create User Err %s", err)
		return "", apperrors.MongoCreateUserFailedError.AppendMessage(err)
	}

	objectId := res.InsertedID.(primitive.Binary).Data
	objectUUID, err := uuid.FromBytes(objectId)
	if err != nil {
		ur.log.Errorf("mongo Create User conversion to uuid er %s", err)
		return "", apperrors.MongoCreateUserFailedError.AppendMessage(err)
	}

	ur.log.Info("ID result SaveUser ", objectUUID)
	return objectUUID.String(), nil
}

func (ur *userRepository) GetUsersByPageAndPerPage(ctx context.Context, page, perPage int) ([]*models.User, error) {
	offset := (page - 1) * perPage
	filter := bson.D{}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetSkip(int64(offset)).
		SetLimit(int64(perPage))

	cursor, err := ur.collection.Find(ctx, filter, opts)
	if err != nil {
		ur.log.Errorf("Mongo get users by page err %s", err)
		return nil, apperrors.MongoGetFailedError.AppendMessage(err)
	}

	defer cursor.Close(ctx)
	ur.log.Info("Get users by page and per page MONGO_DB")
	return decodeUsers(ctx, cursor)
}

func (ur *userRepository) GetUserByID(ctx context.Context, userUUID *uuid.UUID) (*models.User, error) {
	filter := bson.M{"_id": userUUID}
	res := ur.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		ur.log.Errorf("mongo get user by id, cann't find user by _id. Err: %+v", res.Err())
		return nil, apperrors.MongoGetUserByIDErr.AppendMessage(res.Err())
	}

	user := models.User{}
	err := res.Decode(&user)
	if err != nil {
		ur.log.Errorf("mongo get user by id err %s", err)
		return nil, apperrors.MongoGetUserByIDErr.AppendMessage(res.Err())
	}

	ur.log.Info("The user has been finded, successfully.")
	return &user, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	filter := bson.M{"user_email": email}
	res := ur.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		ur.log.Errorf("Cannot find user by email. Err: %+v", res.Err())
		return nil, apperrors.MongoGetUserByEmailErr.AppendMessage(res.Err())
	}

	user := models.User{}
	err := res.Decode(&user)
	if err != nil {
		ur.log.Errorf("mongo get user by email %s", err)
		return nil, apperrors.MongoGetUserByEmailErr.AppendMessage(err)
	}

	ur.log.Info("The user has been finded, successfully.")
	return &user, nil
}

func (ur *userRepository) UpdateUserByID(ctx context.Context, user *models.User) (string, error) {
	filter := bson.M{"_id": user.ID}
	ur.log.Infof("USER %v", user)
	update := bson.M{
		"$set": bson.M{
			"user_email": user.Email,
			"first_name": user.FirstName,
			"user_name":  user.UserName,
			"last_name":  user.LastName,
			"updated_at": user.Updated,
		},
	}

	res, err := ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		ur.log.Errorf("mongo update user by id updateOne %s", err)
		return "", apperrors.MongoUpdateUserByIDErr.AppendMessage(err)
	}

	if res.MatchedCount != 1 {
		ur.log.Error("mongo update user by id updateOne MatchedCount")
		return "", apperrors.MongoUpdateUserByIDErr.AppendMessage("not found %v", user.ID)
	}

	if res.ModifiedCount != 1 {
		ur.log.Errorf("mongo update user by id ModifiedCount ")
		return "", apperrors.MongoUpdateUserByIDErr.AppendMessage("nothing is modified %v", user.ID)
	}

	ur.log.Infof(" result  UpdateUserById %+v", res)
	return user.Email, nil
}

func (ur *userRepository) DropUserByID(ctx context.Context, userUUID *uuid.UUID) error {
	filter := bson.M{"_id": userUUID}
	res, err := ur.collection.DeleteOne(ctx, filter)
	if err != nil {
		appErr := apperrors.MongoDropUserByIDErr.AppendMessage(err)
		ur.log.Errorf("Cannot find user by _id. Err: %+v", err)
		return appErr
	}

	if res.DeletedCount != 1 {
		ur.log.Errorf("mongo DeletedCount == 0")
		return apperrors.MongoDropUserByIDErr.AppendMessage("nothing was deleted %v", userUUID)
	}

	ur.log.Info("The user has been deleted, successfully.")
	return nil
}

func decodeUsers(ctx context.Context, cursor *mongo.Cursor) ([]*models.User, error) {
	defer cursor.Close(ctx)
	var users []*models.User
	for cursor.Next(ctx) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, apperrors.MongoGetFailedError.AppendMessage(err)
		}

		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, apperrors.MongoGetFailedError.AppendMessage(err)
	}

	return users, nil
}
