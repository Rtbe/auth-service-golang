package mongo

import (
	"context"
	"errors"
	"log"

	"example.com/auth-service-go/config"
	"example.com/auth-service-go/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//TokenRepository is an token entity related abstraction for interacting with mongoDB.
type TokenRepository struct {
	cl         *mongo.Client
	collection string
}

//NewTokenRepository returns a new TokenRepository.
func NewTokenRepository(cl *mongo.Client, coll string) *TokenRepository {
	return &TokenRepository{
		cl:         cl,
		collection: coll,
	}
}

//Insert inserts pair of tokens into mongoDB.
func (t *TokenRepository) Insert(ctx context.Context, tokenPair *entity.TokenPair) error {
	cfg := config.New()
	log.Printf("Inserting tokens into mongoDB. Database name: %s, Collection: %s", cfg.DbName, t.collection)

	//Convert refresh token into bcrypt hash before inserting it in mongoDB.
	refreshTokenHash, err := entity.GenerateHash(tokenPair.RefreshToken.Token)
	if err != nil {
		return errors.New("Error generating hash for refresh token")
	}

	//Insert refresh token into mongoDB.
	refreshToken := entity.RefreshToken{
		UUID:      tokenPair.RefreshToken.UUID,
		UserID:    tokenPair.RefreshToken.UserID,
		Token:     refreshTokenHash,
		ExpiresAt: tokenPair.RefreshToken.ExpiresAt,
		Used:      tokenPair.RefreshToken.Used,
	}
	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		if _, err := t.cl.Database(cfg.DbName).Collection(t.collection).InsertOne(sessCtx, &refreshToken); err != nil {
			return nil, err
		}
		return nil, nil
	}

	session, err := t.cl.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("Tokens were successfully stored in mongoDB")
	return nil
}

//DeleteRefreshToken deletes particular refresh token from mongoDB.
func (t *TokenRepository) DeleteRefreshToken(ctx context.Context, userID, refreshTokenUUID string) error {
	cfg := config.New()
	log.Printf("Deleting refresh token: %s from MongoDB. Database name: %s, Collection: %s", refreshTokenUUID, cfg.DbName, t.collection)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		filter := bson.M{"_id": refreshTokenUUID, "user_id": userID}
		result, err := t.cl.Database(cfg.DbName).Collection(t.collection).DeleteOne(sessCtx, filter)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	session, err := t.cl.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Println("Refresh token was successfully deleted from mongoDB")
	return nil
}

//RefreshTokenSetIsUsed sets field used to true for particular refresh token.
func (t *TokenRepository) RefreshTokenSetIsUsed(ctx context.Context, refreshTokenUUID string) error {
	cfg := config.New()
	log.Printf("Updating refresh token: %s in MongoDB. Database name: %s, Collection: %s", refreshTokenUUID, cfg.DbName, t.collection)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		refreshTokenFilter := bson.M{"_id": refreshTokenUUID}
		result, err := t.cl.Database(cfg.DbName).Collection(t.collection).UpdateOne(sessCtx, refreshTokenFilter, bson.M{"$set": bson.M{"used": true}})
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	session, err := t.cl.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	updated := int(result.(*mongo.UpdateResult).ModifiedCount)
	if updated == 0 {
		log.Println("There was no such refresh token in mongoDB")
		return errors.New("There was no such refresh token")
	}
	log.Println("Refresh token has been successfully updated")
	return nil
}

//DeleteUserRefreshTokens deletes all tokens from mongoDB that relates to particular user.
func (t *TokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID string) error {
	cfg := config.New()
	log.Printf("Deleting all tokens from MongoDB related to user with id=%v. Database name: %s, Collection: %s.", userID, cfg.DbName, t.collection)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		filter := bson.M{"user_id": userID}
		result, err := t.cl.Database(cfg.DbName).Collection(t.collection).DeleteMany(sessCtx, filter)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	session, err := t.cl.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	result, err := session.WithTransaction(ctx, callback)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	deletedCount := int(result.(*mongo.DeleteResult).DeletedCount)
	log.Printf("%v records was deleted from mongoDB", deletedCount)
	return nil
}

//IsUserInDB check existence of particular user by given id in mongoDB.
func (t *TokenRepository) IsUserInDB(ctx context.Context, userID string) bool {
	if userID == "" {
		return false
	}
	cfg := config.New()
	log.Printf("Searching user with id=%v in MongoDB. Database name: %s, Collection: %s", userID, cfg.DbName, t.collection)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		filter := bson.M{"user_id": userID}
		result := t.cl.Database(cfg.DbName).Collection(t.collection).FindOne(sessCtx, filter)
		//Check in case of no documents was found.
		if err := result.Err(); err != nil {
			return nil, err
		}
		return result, nil
	}

	session, err := t.cl.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, callback)
	if err == mongo.ErrNoDocuments {
		log.Println("There is no such user in MongoDB")
		return false
	}
	log.Println("There is such user in mongoDB")
	return true
}

//IsRefreshTokenInDB checks existence of particular unused refresh token in mongoDB.
func (t *TokenRepository) IsRefreshTokenInDB(ctx context.Context, refreshTokenUUID string) bool {
	if refreshTokenUUID == "" {
		return false
	}
	cfg := config.New()
	log.Printf("Searching refresh token with id=%v in MongoDB. Database name: %s, Collection: %s", refreshTokenUUID, cfg.DbName, t.collection)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		//This filter makes it seems like only unused refresh tokens are in mongoDB.
		filter := bson.M{"_id": refreshTokenUUID, "used": false}
		result := t.cl.Database(cfg.DbName).Collection(t.collection).FindOne(sessCtx, filter)
		//Check in case of no documents was found.
		if err := result.Err(); err != nil {
			return nil, err
		}
		return result, nil
	}

	session, err := t.cl.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, callback)
	if err == mongo.ErrNoDocuments {
		log.Println("There is no such refresh token in MongoDB")
		return false
	}
	log.Println("There is such refresh token in mongoDB")
	return true
}
