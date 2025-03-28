package service

import (
	"context"
	"data-management/src/model"
	secret "data-management/src/util/db"
	db "data-management/src/util/db/mongo"
	"data-management/src/util/encryption/uaes"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

//* Change this service content with ALT+C(Case sensitive) then CTRL+D this:
//? dataSet and DataSet

type DataSetService struct {
	collectionName string
	ctx            context.Context
	dbUtil         *db.MongoDbUtil
}

func NewDataSetService() *DataSetService {
	this := &DataSetService{
		collectionName: "dataSet",
		ctx:            context.Background(),
	}
	this.dbUtil = db.NewMongoDbUtilUseEnv(this.collectionName)

	return this
}

func (this *DataSetService) BaseGetAll(param model.DataSet_Search, collection *db.MongoDbUtil) (data []model.DataSet_View,
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

func (this *DataSetService) GetAll(param model.DataSet_Search) (data []model.DataSet_View, resEn string, metadata model.MetadataResponse) {
	if os.Getenv("PROD_MODE") == "true" {
		SECRET := secret.GenerateRandomString(7)
		log.Println("SECRET", SECRET)
		var Uaes = uaes.NewAES(SECRET)
		resP, _ := this.BaseGetAll(param, this.dbUtil)
		resEn, _ := Uaes.Encrypt_Any(resP)
		return data, SECRET + resEn, metadata
	} else {
		res, data := this.BaseGetAll(param, this.dbUtil)
		return res, "", data
	}
	// return this.BaseGetAll(param, this.dbUtil)
}

func (this *DataSetService) GetOne(key, value string) (res model.DataSet_View, resEn string, errMessage string) {
	// this.dbUtil.FindOne(key, value, &res)
	if os.Getenv("PROD_MODE") == "true" {
		SECRET := secret.GenerateRandomString(7)
		log.Println("SECRET", SECRET)
		var Uaes = uaes.NewAES(SECRET)
		this.dbUtil.FindOne(key, value, &res)
		resEn, _ = Uaes.Encrypt_Any(res)
		return res, SECRET + resEn, errMessage
	} else {
		this.dbUtil.FindOne(key, value, &res)
		return
	}
	// return
}

func (this *DataSetService) Upsert(param model.DataSet, isUpdate bool) (resp model.Response) {

	upsertErr, upsertId := this.dbUtil.UpsertAndGetId(isUpdate, &param)
	resp.Metadata.Message = upsertErr
	resp.Data = model.Response_Data_Upsert{
		ID: upsertId,
	}
	return
}

func (this *DataSetService) DeleteOne(key, value string) (errMessage string) {
	errMessage = this.dbUtil.DeleteOne(key, value)
	return
}
