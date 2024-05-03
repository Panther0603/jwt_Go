package controller

import (
	"context"
	customerrors "jwt_use/customErrors"
	"jwt_use/models"
	"jwt_use/tokens"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthReponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

func Login() gin.HandlerFunc {

	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		var authrequestdata *AuthRequest
		err := c.BindJSON(&authrequestdata)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": customerrors.ErrInvalidAuthReq})
			return
		}

		if authrequestdata.Username == "" || authrequestdata.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": customerrors.ErrInvalidAuthReq})
			return
		}

		var user models.User

		err = UserCollection.FindOne(ctx, bson.M{"username": authrequestdata.Username}).Decode(&user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrUserNotFound.Error()})
			return
		}

		err = VerifyPassword(user.Password, authrequestdata.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrUsernamePasswordMismatch})
			return
		}

		token, err := tokens.GenerateToken(user.Username, user.Email, user.UserId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		filter := bson.D{{Key: "_id", Value: user.Id}}
		update := bson.M{"$set": bson.M{"token": token}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": customerrors.ErrUpdateAuthToken})
			c.Abort()
			return
		}

		response := &AuthReponse{
			Username: user.Username,
			Token:    token,
		}

		c.JSON(http.StatusOK, response)
	}
}
