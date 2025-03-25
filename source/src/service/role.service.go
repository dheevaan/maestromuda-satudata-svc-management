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
//? role and Role

type RoleService struct {
	collectionName string
	ctx            context.Context
	dbUtil         *db.MongoDbUtil
}

func NewRoleService() *RoleService {
	this := &RoleService{
		collectionName: "role",
		ctx:            context.Background(),
	}
	this.dbUtil = db.NewMongoDbUtilUseEnv(this.collectionName)

	return this
}

func (this *RoleService) BaseGetAll(param model.Role_Search, collection *db.MongoDbUtil) (data []model.Role_View,
	metadata model.MetadataResponse) {
	var err error
	var pipeline []bson.M
	filter := bson.M{}
	listFilterAnd := []bson.M{}
	param.HandleFilter(&listFilterAnd)

	if len(listFilterAnd) > 0 {
		filter["$and"] = listFilterAnd
	}

	// pipeline = append(pipeline, bson.M{"$lookup": bson.M{"from": "user", "localField": "_id", "foreignField": "id", "as": "role"}})
	pipeline = append(pipeline, bson.M{"$match": bson.D{{"status", bson.D{{"$ne", "archive"}}}}})
	pipeline = append(pipeline, bson.M{"$lookup": bson.M{"from": "user", "localField": "_id", "foreignField": "roleId", "as": "users"}})
	pipeline = append(pipeline, bson.M{"$project": bson.M{
		"_id":         1,
		"createdAt":   1,
		"updatedAt":   1,
		"name":        1,
		"description": 1,
		"privileges":  1,
		"alias":       1,
		"level":       1,
		"users":       "$users",
		"userCount":   bson.M{"$size": bson.A{"$users"}}}})
	log.Println(pipeline)

	metadata.Pagination, err = collection.AggsPagination(pipeline, param.Request_Pagination, &data)
	if err != nil {
		metadata.Message = err.Error()
	}
	return
}

func (this *RoleService) GetAll(param model.Role_Search) (data []model.Role_View, resEn string, metadata model.MetadataResponse) {
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

func (this *RoleService) GetOne(key, value string) (res model.Role_View, resEn string, errMessage string) {
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

func (this *RoleService) Upsert(param model.Role, isUpdate bool) (resp model.Response) {

	upsertErr, upsertId := this.dbUtil.UpsertAndGetId(isUpdate, &param)
	resp.Metadata.Message = upsertErr
	resp.Data = model.Response_Data_Upsert{
		ID: upsertId,
	}
	return
}

func (this *RoleService) DeleteOne(key, value string) (errMessage string) {
	errMessage = this.dbUtil.DeleteOne(key, value)
	return
}
