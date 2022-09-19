package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"pintap/model"
	"pintap/utils"
	"strings"
	"time"
	"github.com/dgrijalva/jwt-go"
)
type UsersController struct {
	DB *gorm.DB
}

func (controller *UsersController) Register(ctx *gin.Context) {
	var (
		req  model.User
		user model.User
		userView model.UserViews
		empty model.Empty
	)
	req.LastLogin = time.Now()
	// get data from request parameter
	if err := ctx.BindJSON(&req); err != nil {
		res := utils.Response(http.StatusBadRequest, "can't get data, err :", empty)
		ctx.JSON(http.StatusOK, res)
		return
	}
	// validation request
	// request validation
	validate, trans := utils.InitValidate()
	if err := validate.Struct(req); err != nil {
		var errStrings []string
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			// can translate each error one at a time.
			errStrings = append(errStrings, e.Translate(trans))
		}
		res := utils.Response(http.StatusBadRequest, strings.Join(errStrings, ", "), empty)
		ctx.JSON(http.StatusOK, res)
		return
	}
	// validation for password format
	if err := utils.VerifyPassword(req.Password); err != nil {
		res := utils.Response(http.StatusBadRequest, err.Error(), empty)
		ctx.JSON(http.StatusOK, res)
		return
	}
	// check if confirmation password is not same with password
	if req.Password != req.PasswordConfirmation {

		res := utils.Response(http.StatusBadRequest, "confirmation password is doesn't match", empty)
		ctx.JSON(http.StatusOK, res)
		return
	}
	// check one is exist
	if err := controller.DB.Where("phone LIKE ?", req.Phone).First(&user).Error; !(gorm.IsRecordNotFoundError(err)) {
		res := utils.Response(http.StatusBadRequest, "Phone is already taken ", empty)
		ctx.JSON(http.StatusOK, res)
		return
	}

	// check name is exist
	if err := controller.DB.Where("nama LIKE ?", req.Name).First(&user).Error; !(gorm.IsRecordNotFoundError(err)) {
		res := utils.Response(http.StatusBadRequest, "Name is already taken ", empty)
		ctx.JSON(http.StatusOK, res)
		return
	}

	// assign hash password to user and create user
	pass, _ := utils.HashAndSalt([]byte(req.Password))
	req.Password = pass
	req.PasswordConfirmation = pass
	if err := controller.DB.Create(&req).Error; err != nil {
		res := utils.Response(http.StatusInternalServerError, "can't create user, err: "+err.Error(), empty)
		ctx.JSON(http.StatusOK, res)
		return
	}
	userView.ID = req.ID
	userView.Name = req.Name
	res := utils.Response(http.StatusOK, "Success Register", userView)
	ctx.JSON(http.StatusOK, res)
}
func (controller *UsersController) Login(ctx *gin.Context)  {
	var (
		req  model.User
		user model.User
		userView model.UserView
		empty model.Empty
	)
	if err := ctx.BindJSON(&req); err != nil {
		res := utils.Response(http.StatusBadRequest, "can't bind struct, err:"+err.Error(), empty)
		ctx.JSON(http.StatusOK, res)
		return
	}

	if err := controller.DB.Set("gorm:auto_preload", true).Where("nama LIKE ?", req.Name).Debug().Find(&user).Error; gorm.IsRecordNotFoundError(err) {
		res := utils.Response(http.StatusBadRequest, "Name is incorrect", empty)
		ctx.JSON(http.StatusOK, res)
		return
	}

	// create expired time for auth token
	expiresAt := time.Now().Add(time.Minute * 1000).Unix()
	// check if password not same
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		res := utils.Response(http.StatusBadRequest, "password is incorrect", empty)
		ctx.JSON(http.StatusOK, res)
		return
	}

	tk := &model.Token{
		ID: user.ID,
		Name: user.Name,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	// generate jwt auth
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "can't signed token"})
		return
	}
	// assign jwt token to auth token user
	user.AuthToken = tokenString
	user.LastLogin = time.Now()
	// update user
	if err := controller.DB.Save(&user).Error; err != nil {
		res := utils.Response(http.StatusInternalServerError, "can't update user, err: "+err.Error(), empty)
		ctx.JSON(http.StatusOK, res)
		return
	}
	userView.ID = user.ID
	userView.Name = user.Name
	userView.AuthToken = user.AuthToken
	res := utils.Response(http.StatusOK, "Success Login", userView)
	ctx.JSON(http.StatusOK, res)
	return
}
func (controller *UsersController) GetProfil(ctx *gin.Context)  {
	id := ctx.Params.ByName("id")
	var (
		user model.User
		empty model.Empty
		view model.UserViews
	)

	qry_customer := controller.DB.Set("gorm:auto_preload", true).Where("id =?", id).Find(&user)

	if qry_customer.RecordNotFound() {
		res := utils.Response(http.StatusNotFound, "Record User Not Found", empty)
		ctx.JSON(http.StatusOK, res)
		return

	}
	if qry_customer.Error != nil {
		errMsg := qry_customer.Error.Error()
		ctx.JSON(http.StatusInternalServerError,errMsg)
		return

	}
	view.ID = user.ID
	view.Name = user.Name
	res := utils.Response(http.StatusOK, "Success", view)
	ctx.JSON(http.StatusOK, res)
	return
}

func (controller *UsersController) GetAllUser(ctx *gin.Context)  {
	var (
		user []model.User
		empty model.Empty
	)

	qry_faq := controller.DB.Set("gorm:auto_preload", true).Debug().Find(&user)

	if qry_faq.RecordNotFound() {
		res := utils.Response(http.StatusNotFound, "Record Info Not Found", empty)
		ctx.JSON(http.StatusOK, res)
		return

	}
	if qry_faq.Error != nil {
		errMsg := qry_faq.Error.Error()
		ctx.JSON(http.StatusInternalServerError,errMsg)
		return

	}

	res := utils.Response(http.StatusOK, "Success", user)
	ctx.JSON(http.StatusOK, res)
	return

}

func (controller *UsersController) UpdateUser(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	var (
		req model.User
		user model.User
		empty model.Empty
		change model.UserViews
	)
	updt_user := controller.DB.Debug().Where("id = ?", id).Find(&user)
	if updt_user.RecordNotFound() {
		res := utils.Response(http.StatusNotFound, "Record User Not Found", empty)
		ctx.JSON(http.StatusOK, res)
		return
	}

	if err := ctx.BindJSON(&req); err != nil {
		res := utils.Response(http.StatusBadRequest, "can't bind struct, err:"+err.Error(), empty)
		ctx.JSON(http.StatusOK, res)
		return
	}

	if req.Password != req.PasswordConfirmation {

		res := utils.Response(http.StatusBadRequest, "confirmation password is doesn't match", empty)
		ctx.JSON(http.StatusOK, res)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(hash)
	user.Name = req.Name
	req.LastLogin = time.Now()
	if err := controller.DB.Save(&user).Error; err != nil {
		res := utils.Response(http.StatusInternalServerError, "can't update user, err: "+err.Error(), empty)
		ctx.JSON(http.StatusOK, res)
		return
	}

	change.ID = user.ID
	change.Name = user.Name
	res := utils.Response(http.StatusOK, "Success Change User", change)
	ctx.JSON(http.StatusOK, res)

}

func (controller *UsersController) Deleteuser(ctx *gin.Context)  {
	id := ctx.Params.ByName("id")
	var (
		user model.User
		empty model.Empty
	)

	qry_customer := controller.DB.Set("gorm:auto_preload", true).Where("id =?", id).Delete(&user)

	if qry_customer.RecordNotFound() {
		res := utils.Response(http.StatusNotFound, "Record User Not Found", empty)
		ctx.JSON(http.StatusOK, res)
		return

	}
	if qry_customer.Error != nil {
		errMsg := qry_customer.Error.Error()
		ctx.JSON(http.StatusInternalServerError,errMsg)
		return

	}
	res := utils.Response(http.StatusOK, "Delete Success", empty)
	ctx.JSON(http.StatusOK, res)
	return
}