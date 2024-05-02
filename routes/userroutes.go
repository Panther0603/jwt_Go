package routes

import (
	usercontroller "jwt_use/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {

	usergroup := incomingRoutes.Group("/users")

	usergroup.POST("/", usercontroller.CreateUser())
	usergroup.GET("/:id", usercontroller.GetUserById())
	usergroup.GET("", usercontroller.GetAllUser())
	usergroup.DELETE("", usercontroller.DeleteUserBuyId())
	usergroup.PUT("/update", usercontroller.UpdateUser())
}
