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
//? role and Role

type RoleController struct {
	router  *gin.RouterGroup
	service *service.RoleService
}

func NewRoleController(router *gin.RouterGroup) *RoleController {
	this := &RoleController{router: router, service: service.NewRoleService()}

	role := this.router.Group("/role")
	role.POST("/get-all", this.GetAll)
	role.GET("/get-one", this.GetOne)
	role.POST("/add", this.Add)
	role.PUT("/update", this.Update)
	role.DELETE("/delete", this.DeleteOne)

	return this
}

// @Tags role
// @Accept json
// @Param parameter body model.Role_Search true "PARAM"
// @Produce json
// @Success 200 {object} object{data=[]model.Role_View,meta_data=model.MetadataResponse} "OK"
// @Router /role/get-all [post]
// @Security JWT
func (this *RoleController) GetAll(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Role_Search
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	if os.Getenv("PROD_MODE") == "true" {
		_, resp.Data, resp.Metadata = this.service.GetAll(param)
	} else {
		resp.Data, _, resp.Metadata = this.service.GetAll(param)
	}
	// resp.Data, resp.Metadata = this.service.GetAll(param)
}

// @Tags role
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{data=model.Role_View,meta_data=model.MetadataResponse} "OK"
// @Router /role/get-one [get]
// @Security JWT
func (this *RoleController) GetOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	if os.Getenv("PROD_MODE") == "true" {
		_, resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"))
	} else {
		resp.Data, _, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"))
	}
	// resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"))
}

// @Tags role
// @Accept json
// @Param parameter body model.Role true "PARAM"
// @Produce json
// @Success 201 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /role/add [post]
// @Security JWT
func (this *RoleController) Add(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Role
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, false)
}

// @Tags role
// @Accept json
// @Param parameter body model.Role true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /role/update [put]
// @Security JWT
func (this *RoleController) Update(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Role
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, true)
}

// @Tags role
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /role/delete [delete]
// @Security JWT
func (this *RoleController) DeleteOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	resp.Metadata.Message = this.service.DeleteOne("_id", ctx.Query("id"))
}
