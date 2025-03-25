package controller

import (
	"data-management/src/model"
	"data-management/src/service"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

//* Change this controller content with ALT+C(Case sensitive) then CTRL+D this:
//! and don't forget to register the controller in main.go
//? user and User

type UserController struct {
	router  *gin.RouterGroup
	service *service.UserService
}

func NewUserController(router *gin.RouterGroup) *UserController {
	this := &UserController{router: router, service: service.NewUserService()}

	user := this.router.Group("/user")
	users := this.router.Group("/auth")
	user.POST("/get-all", this.GetAll)
	user.GET("/get-one", this.GetOne)
	user.POST("/add", this.Add)
	user.PUT("/update", this.Update)
	user.DELETE("/delete", this.DeleteOne)
	user.PUT("/reset-password", this.Change_Password_Admin)
	users.PUT("/update-profile", this.Update_User)
	users.PUT("/change-password", this.Change_Password)

	return this
}

// @Tags user
// @Accept json
// @Param parameter body model.User_Search true "PARAM"
// @Produce json
// @Success 200 {object} object{data=[]model.User_View,meta_data=model.MetadataResponse} "OK"
// @Router /user/get-all [post]
// @Security JWT
func (this *UserController) GetAll(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.User_Search
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}
	if os.Getenv("PROD_MODE") == "true" {
		_, resp.Data, resp.Metadata = this.service.GetAll(param)
	} else {
		resp.Data, _, resp.Metadata = this.service.GetAll(param)
	}
}

// @Tags user
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{data=model.User_View,meta_data=model.MetadataResponse} "OK"
// @Router /user/get-one [get]
// @Security JWT
func (this *UserController) GetOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	if os.Getenv("PROD_MODE") == "true" {
		_, resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"))
	} else {
		resp.Data, _, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"))
	}
}

// @Tags user
// @Accept json
// @Param parameter body model.User true "PARAM"
// @Produce json
// @Success 201 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /user/add [post]
// @Security JWT
func (this UserController) Add(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.User
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.UpsertWithHashingPassword(param, false)
}

// @Tags user
// @Accept json
// @Param parameter body model.User true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /user/update [put]
// @Security JWT
func (this *UserController) Update(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.User
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.UpsertWithHashingPassword(param, true)
}

// @Tags user
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /user/delete [delete]
// @Security JWT
func (this *UserController) DeleteOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	resp.Metadata.Message = this.service.DeleteOne("_id", ctx.Query("id"))
}

// @Tags auth
// @Accept json
// @Param parameter body model.User_Profil true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /auth/update-profile [put]
// @Security JWT
func (this *UserController) Update_User(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.User_Profil
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.UpsertWithHashPassword(param, true)
}

// @Tags auth
// @Accept json
// @Param parameter body model.User_Profil true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /auth/update-profile [put]
// @Security JWT
func (this *UserController) Update_User_Admin(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.User
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.UpsertWithHashingPasswordAdmin(param, true)
}

// @Tags auth
// @Accept json
// @Param parameter body model.User_ResetPassword true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /auth/change-password [put]
// @Security JWT
func (this *UserController) Change_Password(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.User_ResetPassword
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.ResetPassword(param)
}

// @Tags user
// @Accept json
// @Param parameter body model.User_ResetPassword_Admin true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /user/reset-password [put]
// @Security JWT
func (this *UserController) Change_Password_Admin(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.User_ResetPassword_Admin
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.ResetPasswordAdmin(param)
}
