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
//? catalog and Catalog

type CatalogService struct {
	collectionName string
	ctx            context.Context
	dbUtil         *db.MongoDbUtil
}

func NewCatalogService() *CatalogService {
	this := &CatalogService{
		collectionName: "catalog",
		ctx:            context.Background(),
	}
	this.dbUtil = db.NewMongoDbUtilUseEnv(this.collectionName)

	return this
}

func (this *CatalogService) BaseGetAll(param model.Catalog_Search, collection *db.MongoDbUtil) (data []model.Catalog_View,
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

func (this *CatalogService) ChangeCollectionName(collName string) {
	this.dbUtil = db.NewMongoDbUtilUseEnv(collName)
}

func (this *CatalogService) GetAll(param model.Catalog_Search) (data []model.Catalog_View, resEn string, metadata model.MetadataResponse) {
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

func (this *CatalogService) GetOne(key, value string, collectionName string) (res model.Catalog_View, resEn string, errMessage string) {
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

func (this *CatalogService) Upsert(param model.Catalog, isUpdate bool) (resp model.Response) {
	upsertErr, upsertId := this.dbUtil.UpsertAndGetId(isUpdate, &param)
	resp.Metadata.Message = upsertErr
	resp.Data = model.Response_Data_Upsert{
		ID: upsertId,
	}

	return
}

func (this *CatalogService) DeleteOne(key, value string) (errMessage string) {
	errMessage = this.dbUtil.DeleteOne(key, value)
	return
}

func (this *CatalogService) GetAllCatalog(param model.Catalog_Search) (data []model.Catalog_View, resEn string, metadata model.MetadataResponse) {
	this.ChangeCollectionName("sdgs_kuesioner_individu")
	if os.Getenv("PROD_MODE") == "true" {
		SECRET := secret.GenerateRandomString(7)
		log.Println("SECRET", SECRET)
		var Uaes = uaes.NewAES(SECRET)
		resP, _ := this.BaseGetAll(param, this.dbUtil)
		resEn, _ := Uaes.Encrypt_Any(resP)
		return data, SECRET+resEn, metadata
	} else {
		res, data := this.BaseGetAll(param, this.dbUtil)
		return res, "", data
	}
	// return this.BaseGetAll(param, this.dbUtil)
}

func (this *CatalogService) GetOneCatalog(key, value string, collectionName string) (res model.Catalog_View, resEn string, errMessage string) {
	// this.dbUtil.FindOne(key, value, &res)
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
