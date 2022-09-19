package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"pintap/controllers"
	"pintap/utils"
)

type Routes struct {
	DB *gorm.DB
}

func (r *Routes) Setup(port string) {

	app := gin.Default()
	//app.MaxMultipartMemory = 2048
	fmt.Println("check memory an : ",app.MaxMultipartMemory)
	user := app.Group("users")
	{
		usersCtrl := controllers.UsersController{DB: r.DB}
		user.POST("/register", usersCtrl.Register)
		user.POST("/login", usersCtrl.Login)
		user.GET("/profil/:id", usersCtrl.GetProfil)
		user.GET("/",utils.Middleware, usersCtrl.GetAllUser)
		user.PATCH("/update/:id",usersCtrl.UpdateUser)
		user.DELETE("/delete/:id", usersCtrl.Deleteuser)


	}

	runningOnPort := fmt.Sprintf(":%s", port)
	app.Run(runningOnPort)
}