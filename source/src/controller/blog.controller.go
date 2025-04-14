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
//? blog and Blog

type BlogController struct {
	router  *gin.RouterGroup
	service *service.BlogService
}

func NewBlogController(router *gin.RouterGroup) *BlogController {
	this := &BlogController{router: router, service: service.NewBlogService()}

	blog := this.router.Group("/blog")
	blog.POST("/get-all", this.GetAll)
	blog.GET("/get-one", this.GetOne)
	blog.POST("/add", this.Add)
	blog.PUT("/update", this.Update)
	blog.DELETE("/delete", this.DeleteOne)

	return this
}

// @Tags blog
// @Accept json
// @Param parameter body model.Blog_Search true "PARAM"
// @Produce json
// @Success 200 {object} object{data=[]model.Blog_View,meta_data=model.MetadataResponse} "OK"
// @Router /blog/get-all [post]
// @Security JWT
func (this *BlogController) GetAll(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Blog_Search
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

// @Tags blog
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{data=model.Blog_View,meta_data=model.MetadataResponse} "OK"
// @Router /blog/get-one [get]
// @Security JWT
func (this *BlogController) GetOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	if os.Getenv("PROD_MODE") == "true" {
		_, resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"), ctx.Param("blog"))
	} else {
		resp.Data, _, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"), ctx.Param("blog"))
	}
	// resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"), ctx.Param("blog"))
}

// @Tags blog
// @Accept json
// @Param parameter body model.Blog true "PARAM"
// @Produce json
// @Success 201 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /blog/add [post]
// @Security JWT
func (this *BlogController) Add(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Blog
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, false)
}

// @Tags blog
// @Accept json
// @Param parameter body model.Blog true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /blog/update [put]
// @Security JWT
func (this *BlogController) Update(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Blog
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, true)
}

// @Tags blog
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /blog/delete [delete]
// @Security JWT
func (this *BlogController) DeleteOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	resp.Metadata.Message = this.service.DeleteOne("_id", ctx.Query("id"))
}
