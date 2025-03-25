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
//? dataSet and DataSet

type DataSetController struct {
	router  *gin.RouterGroup
	service *service.DataSetService
}

func NewDataSetController(router *gin.RouterGroup) *DataSetController {
	this := &DataSetController{router: router, service: service.NewDataSetService()}

	dataSet := this.router.Group("/dataSet")
	dataSet.POST("/get-all", this.GetAll)
	dataSet.GET("/get-one", this.GetOne)
	dataSet.POST("/add", this.Add)
	dataSet.PUT("/update", this.Update)
	dataSet.DELETE("/delete", this.DeleteOne)

	return this
}

// @Tags dataSet
// @Accept json
// @Param parameter body model.DataSet_Search true "PARAM"
// @Produce json
// @Success 200 {object} object{data=[]model.DataSet_View,meta_data=model.MetadataResponse} "OK"
// @Router /dataSet/get-all [post]
// @Security JWT
func (this *DataSetController) GetAll(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.DataSet_Search
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

// @Tags dataSet
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{data=model.DataSet_View,meta_data=model.MetadataResponse} "OK"
// @Router /dataSet/get-one [get]
// @Security JWT
func (this *DataSetController) GetOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	if os.Getenv("PROD_MODE") == "true" {
		_, resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"))
	} else {
		resp.Data, _, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"))
	}
	// resp.Data, resp.Metadata.Message = this.service.GetOne("_id", ctx.Query("id"))
}

// @Tags dataSet
// @Accept json
// @Param parameter body model.DataSet true "PARAM"
// @Produce json
// @Success 201 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /dataSet/add [post]
// @Security JWT
func (this *DataSetController) Add(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.DataSet
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, false)
}

// @Tags dataSet
// @Accept json
// @Param parameter body model.DataSet true "PARAM"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /dataSet/update [put]
// @Security JWT
func (this *DataSetController) Update(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	var param model.DataSet
	if err := ctx.BindJSON(&param); err != nil {
		log.Println(err)
		return
	}

	resp = this.service.Upsert(param, true)
}

// @Tags dataSet
// @Accept json
// @Param id query string true "ID"
// @Produce json
// @Success 200 {object} object{meta_data=model.MetadataResponse} "OK"
// @Router /dataSet/delete [delete]
// @Security JWT
func (this *DataSetController) DeleteOne(ctx *gin.Context) {
	resp := model.Response{}
	defer SetMetadataResponse(ctx, time.Now(), &resp)

	resp.Metadata.Message = this.service.DeleteOne("_id", ctx.Query("id"))
}
