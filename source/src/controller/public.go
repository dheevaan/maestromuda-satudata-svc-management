package controller

import (
	"data-management/src/model"
	"data-management/src/service"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

//* Change this controller content with ALT+C(Case sensitive) then CTRL+D this:
//! and don't forget to register the controller in main.go
//? public and Public

type PublicController struct {
	router *gin.RouterGroup
}

func NewPublicController(router *gin.RouterGroup, jwtMiddleware *jwt.GinJWTMiddleware) *PublicController {
	this := &PublicController{router: router}

	public := this.router.Group("/public")
	auth := public.Group("auth")
	auth.POST("/login", jwtMiddleware.LoginHandler)
	auth.GET("/refresh", jwtMiddleware.RefreshHandler)
	auth.GET("/refresh2", this.Refresh2)

	user := public.Group("/user")
	user.POST("/reset-password", this.ResetPassword)
	user.POST("/register", this.Register)

	return this
}

// @Tags public
// @Accept json
// @Param parameter body model.User_Login true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /public/auth/login [post]
func (this *PublicController) Login(ctx *gin.Context) {}

// @Tags public
// @Accept json
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /public/auth/refresh [get]
// @Security JWT
func (this *PublicController) Refresh(ctx *gin.Context) {}

// @Tags public
// @Accept json
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /public/auth/refresh2 [get]
// @Security JWT
func (this *PublicController) Refresh2(ctx *gin.Context) {

}

// @Tags public
// @Accept json
// @Param parameter body model.User_ResetPassword true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /public/user/reset-password [put]
func (this *PublicController) ResetPassword(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.User_ResetPassword
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = service.NewUserService().ResetPassword(param)
}

// @Tags public
// @Accept json
// @Param parameter body model.User true "PARAM"
// @Produce json
// @Success 201 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /public/user/register [post]
func (this *PublicController) Register(ctx *gin.Context) {
	UserController{service: service.NewUserService()}.Add(ctx)
}
