package service

import (
	"context"
	"log"
	"os"
	"data-management/src/model"
	secret "data-management/src/util/db"
	db "data-management/src/util/db/mongo"
	"data-management/src/util/encryption/uaes"

	"go.mongodb.org/mongo-driver/bson"
)

//* Change this service content with ALT+C(Case sensitive) then CTRL+D this:
//? blog and Blog

type BlogService struct {
	collectionName string
	ctx            context.Context
	dbUtil         *db.MongoDbUtil
}

func NewBlogService() *BlogService {
	this := &BlogService{
		collectionName: "blog",
		ctx:            context.Background(),
	}
	this.dbUtil = db.NewMongoDbUtilUseEnv(this.collectionName)

	return this
}

func (this *BlogService) BaseGetAll(param model.Blog_Search, collection *db.MongoDbUtil) (data []model.Blog_View,
	metadata model.MetadataResponse) {
	filter := bson.M{}
	listFilterAnd := []bson.M{}
	param.HandleFilter(&listFilterAnd)

	if len(listFilterAnd) > 0 {
		filter["$and"] = listFilterAnd
	}
	metadata.Pagination, metadata.Message = collection.Find(filter,
		param.Request, &data)
	return
}

func (this *BlogService) BaseGetAllMap(param model.Blog_Search, collection *db.MongoDbUtil) (data []map[string]interface{}, metadata model.MetadataResponse) {
	filter := bson.M{}
	listFilterAnd := []bson.M{}
	param.HandleFilter(&listFilterAnd)

	if len(listFilterAnd) > 0 {
		filter["$and"] = listFilterAnd
	}

	metadata.Pagination, metadata.Message = collection.Find(filter, param.Request, &data)
	return
}

func (this *BlogService) ChangeCollectionName(collName string) {
	this.dbUtil = db.NewMongoDbUtilUseEnv(collName)
}

func (this *BlogService) GetAll(param model.Blog_Search) (data []model.Blog_View, resEn string, metadata model.MetadataResponse) {
	this.ChangeCollectionName("blog")
	if os.Getenv("PROD_MODE") == "true" {
		SECRET := secret.GenerateRandomString(7)
		log.Println("SECRET", SECRET)
		var Uaes = uaes.NewAES(SECRET)
		resP, _ := this.BaseGetAll(param, this.dbUtil)
		resEn, _ := Uaes.Encrypt_Any(resP)
		return data, resEn+SECRET, metadata
	} else {
		res, data := this.BaseGetAll(param, this.dbUtil)
		return res, "", data
	}
	// return this.BaseGetAll(param, this.dbUtil)
}

func (this *BlogService) GetOne(key, value string, collectionName string) (res model.Blog_View, resEn string, errMessage string) {
	this.ChangeCollectionName("blog")
	// this.dbUtil.FindOne(key, value, &res)
	if os.Getenv("PROD_MODE") == "true" {
		SECRET := secret.GenerateRandomString(7)
		log.Println("SECRET", SECRET)
		var Uaes = uaes.NewAES(SECRET)
		this.dbUtil.FindOne(key, value, &res)
		resEn, _ = Uaes.Encrypt_Any(res)
		return res, resEn+SECRET, errMessage
	} else {
		this.dbUtil.FindOne(key, value, &res)
		return
	}
	// return
}

func (this *BlogService) Upsert(param model.Blog, isUpdate bool) (resp model.Response) {
	this.ChangeCollectionName("blog")
	upsertErr, upsertId := this.dbUtil.UpsertAndGetId(isUpdate, &param)
	resp.Metadata.Message = upsertErr
	resp.Data = model.Response_Data_Upsert{
		ID: upsertId,
	}

	return
}

func (this *BlogService) DeleteOne(key, value string) (errMessage string) {
	this.ChangeCollectionName("blog")
	errMessage = this.dbUtil.DeleteOne(key, value)
	return
}

func (this *BlogService) GetAllSource(param model.Blog_Search) (data []map[string]interface{}, resEn string, metadata model.MetadataResponse) {
	this.ChangeCollectionName("source")
	if os.Getenv("PROD_MODE") == "true" {
		SECRET := secret.GenerateRandomString(7)
		log.Println("SECRET", SECRET)
		var Uaes = uaes.NewAES(SECRET)
		resP, _ := this.BaseGetAllMap(param, this.dbUtil)
		resEn, _ := Uaes.Encrypt_Any(resP)
		return data, SECRET+resEn, metadata
	} else {
		res, data := this.BaseGetAllMap(param, this.dbUtil)
		return res, "", data
	}
}

func (this *BlogService) GetOneBlog(key, value string, collectionName string) (res model.Blog_View, resEn string, errMessage string) {
	// this.dbUtil.FindOne(key, value, &res)
	this.ChangeCollectionName("blog")
	if os.Getenv("PROD_MODE") == "true" {
		SECRET := secret.GenerateRandomString(7)
		log.Println("SECRET", SECRET)
		var Uaes = uaes.NewAES(SECRET)
		this.dbUtil.FindOne(key, value, &res)
		resEn, _ = Uaes.Encrypt_Any(res)
		return res, SECRET+resEn, errMessage
	} else {
		this.dbUtil.FindOne(key, value, &res)
		return
	}
	// return
}
