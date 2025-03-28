package controller

import (
	"log"
	"os"
	"data-management/src/model"
	"data-management/src/service"
	"time"

	"github.com/gin-gonic/gin"
)

//* Change this controller content with ALT+C(Case sensitive) then CTRL+D this:
//! and don't forget to register the controller in main.go
//? catalog and Catalog

type CatalogController struct {
	router  *gin.RouterGroup
	service *service.CatalogService
}

func NewCatalogController(router *gin.RouterGroup) *CatalogController {
	this := &CatalogController{router: router, service: service.NewCatalogService()}

	catalog := this.router.Group("/catalog")
	catalog.POST("/get-all", this.GetAll)
	catalog.GET("/get-one", this.GetOne)
	catalog.POST("/add", this.Add)
	catalog.PUT("/update", this.Update)
	catalog.DELETE("/delete", this.DeleteOne)

	return this
}

// @Tags catalog
// @Accept json
// @Param parameter body model.Catalog_Search true "PARAM"
// @Produce json
// @Success 200 {object} object{data=[]model.Catalog_View,meta_data=model.MetadataResponse} "OK"
// @Router /catalog/get-all [post]
// @Security JWT
func (this *CatalogController) GetAll(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Catalog_Search
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

// @Tags catalog
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{data=model.Catalog_View,meta_data=model.MetadataResponse} "OK"
// @Router /catalog/get-one [get]
// @Security JWT
func (this *CatalogController) GetOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	if os.Getenv("PROD_MODE") == "true"{
		_, resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"), ctx.Param("catalog"))	
	}else{
		resp.Data, _, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"), ctx.Param("catalog"))
	}
	// resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"), ctx.Param("catalog"))
}
// @Tags catalog
// @Accept json
// @Param parameter body model.Catalog true "PARAM"
// @Produce json
// @Success 201 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /catalog/add [post]
// @Security JWT
func (this *CatalogController) Add(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Catalog
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, false)
}

// @Tags catalog
// @Accept json
// @Param parameter body model.Catalog true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /catalog/update [put]
// @Security JWT
func (this *CatalogController) Update(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.Catalog
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, true)
}

// @Tags catalog
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /catalog/delete [delete]
// @Security JWT
func (this *CatalogController) DeleteOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	resp.Metadata.Message = this.service.DeleteOne("_id", ctx.Query("id"))
}
