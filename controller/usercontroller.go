package controller

import (
	"context"
	customerrors "jwt_use/customErrors"
	"jwt_use/database"
	"jwt_use/models"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var getUser models.User
		var err = c.BindJSON(&getUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Here we are adding the implicit data
		getUser.Id = primitive.NewObjectID()
		getUser.UserId = uuid.NewString()
		getUser.UpdatedAt = time.Now()
		getUser.CreatedAt = time.Now()
		getUser.Password = HashPassword(getUser.Password)

		// Now checking if any document is present in the collection
		count, err := UserCollection.CountDocuments(ctx, bson.M{"username": getUser.Username})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// If yes, then return
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrUserExits.Error()})
			return
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"email": getUser.Email})

		// checking exixts with email
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrUserEmailExsits.Error()})
			return
		}

		// checking exixts with phone number
		count, err = UserCollection.CountDocuments(ctx, bson.M{"email": getUser.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrUserEmailExsits.Error()})
			return
		}

		// If not, then insert
		_, err = UserCollection.InsertOne(ctx, getUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrSomethingWentWrong.Error()})
			return
		}

		c.JSON(http.StatusCreated, getUser)
	}
}

func GetUserById() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Query("id")

		log.Print(">>>>>>>>>>" + userId)
		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": customerrors.ErrEmptyUserId.Error()})
			return
		}

		user_id, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": customerrors.ErrUserIDNotValid.Error()})
			return
		}

		var userDetails models.User
		err = UserCollection.FindOne(ctx, bson.M{"_id": user_id}).Decode(&userDetails)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": customerrors.ErrUserNotFound.Error()})
			return
		}

		c.JSON(http.StatusOK, userDetails)
	}
}

func GetAllUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var AllUser []models.User

		userpointcusror, err := UserCollection.Find(ctx, bson.M{})

		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		err = userpointcusror.All(ctx, &AllUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrListParsing.Error()})
		}

		c.JSON(http.StatusOK, AllUser)
	}
}

func DeleteUserBuyId() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userId := c.Query("id")

		if userId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": customerrors.ErrEmptyUserId.Error()})
		}

		user_id, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": customerrors.ErrUserIDNotValid.Error()})
		}

		_, err = UserCollection.DeleteOne(ctx, bson.M{"_id": user_id})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User Not deteletd"})
		}
		c.JSON(http.StatusOK, gin.H{"error": "user deleted success fully"})
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var editUser models.User
		err := c.BindJSON(&editUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad requrst"})
			return
		}

		if editUser.Id.String() == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": customerrors.ErrEmptyUserId.Error()})
			return
		}

		// assigning filter
		filter := bson.D{{Key: "_id", Value: editUser.Id}}

		// now we set teh update value
		update := bson.M{}

		//  will use reflector module to get and iterate over all the data which comes in the user requesy body dynamically '

		ref := reflect.ValueOf(&editUser).Elem()
		typeOfS := ref.Type()

		for i := 0; i < ref.NumField(); i++ {

			fieldName := typeOfS.Field(i).Tag.Get("bson")
			fieldValue := ref.Field(i).Interface()
			if update["$set"] == nil {
				update["$set"] = bson.M{}
			}

			if fieldName == "_id" || fieldName == "userid" || fieldName == "password" || fieldValue == "" {
				continue
			} else if fieldName == "email" || fieldName == "phoneno" || fieldName == "username" {

				if fieldName == "username" {
					if fieldValue == "" {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "username can not be empty"})
						return
					}
					count, err := UserCollection.CountDocuments(ctx, bson.M{"username": editUser.Username})
					if err != nil || count > 0 {
						c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrUpdateUsernameExits.Error()})
						return
					}
					update["$set"].(bson.M)[fieldName] = fieldValue
				}

				if fieldName == "email" {
					if fieldValue == "" {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "email can not be empty"})
						return
					}
					count, err := UserCollection.CountDocuments(ctx, bson.M{"email": editUser.Email})
					if err != nil || count > 0 {
						c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrUpdateEmailExits.Error()})
						return
					}
					update["$set"].(bson.M)[fieldName] = fieldValue
				}

				if fieldName == "phoneno" {
					if fieldValue == "" {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "phoneno can not be empty"})
						return
					}
					count, err := UserCollection.CountDocuments(ctx, bson.M{"phoneno": editUser.PhoneNo})
					if err != nil || count > 0 {
						c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrUpdatePhoneNoeExits})
						return
					}
					update["$set"].(bson.M)[fieldName] = fieldValue
				}
			} else {
				update["$set"].(bson.M)[fieldName] = fieldValue
			}
		}
		update["$set"].(bson.M)["updatedat"] = time.Now()
		updateResult, err := UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrCanNotUpdateUser})
		}
		c.JSON(http.StatusOK, gin.H{"message": "user updated sucessfully", "updatedCount": updateResult.ModifiedCount})
	}

}

func HashPassword(password string) (hashed string) {

	haseedByte, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err.Error())
	}
	return string(haseedByte)
}

func VerifyPassword(dbpassword string, getPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(getPassword))
}
