package service

import (
	"context"
	"data-management/src/model"
	"data-management/src/model/enum"
	secret "data-management/src/util/db"
	db "data-management/src/util/db/mongo"
	"data-management/src/util/encryption/uaes"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

//* Change this service content with ALT+C(Case sensitive) then CTRL+D this:
//? user and User

type UserService struct {
	collectionName string
	ctx            context.Context
	dbUtil         *db.MongoDbUtil
	passwordCost   int
}

func NewUserService() *UserService {
	this := &UserService{
		collectionName: "user",
		ctx:            context.Background(),
		passwordCost:   15,
	}
	this.dbUtil = db.NewMongoDbUtilUseEnv(this.collectionName)

	return this
}

func (this *UserService) ChangeCollectionName(collName string) {
	this.dbUtil = db.NewMongoDbUtilUseEnv(collName)
}

func (this *UserService) BaseGetAll(param model.User_Search, collection *db.MongoDbUtil) (data []model.User_View, metadata model.MetadataResponse) {
	filter := bson.M{}
	var err error
	var pipeline []bson.M
	listFilterAnd := []bson.M{}
	param.HandleFilter(&listFilterAnd)

	if len(listFilterAnd) > 0 {
		filter["$and"] = listFilterAnd
	}

	var ordered int
	if param.Request_Pagination.Order == "ASC" {
		ordered = 1
	} else {
		ordered = -1
	}

	pipeline = append(pipeline, bson.M{"$lookup": bson.M{"from": "role", "localField": "roleId", "foreignField": "_id", "as": "role"}})
	pipeline = append(pipeline, bson.M{"$project": bson.M{"about": 1, "address": 1, "avatar": 1, "createdAt": 1, "email": 1, "emailVerifiedAt": 1, "fullname": 1, "nik": 1, "phone": 1, "roleId": 1, "rt": 1, "rw": 1, "status": 1, "updatedAt": 1, "username": 1, "workspaceId": 1, "workspace": 1, "role": bson.M{"$arrayElemAt": []interface{}{"$role", 0}}}})
	pipeline = append(pipeline, bson.M{"$match": filter})
	pipeline = append(pipeline, bson.M{"$sort": bson.M{param.Request_Pagination.OrderBy: ordered}})
	pip, _ := json.Marshal(filter)
	log.Println(string(pip))

	metadata.Pagination, err = collection.AggsPagination(pipeline, param.Request_Pagination, &data)
	if err != nil {
		metadata.Message = err.Error()
	}

	return
}

func (this *UserService) GetAll(param model.User_Search) (data []model.User_View, resEn string, metadata model.MetadataResponse) {
	// this.ChangeCollectionName("v_detail_user")
	if os.Getenv("PROD_MODE") == "true" {
		SECRET := secret.GenerateRandomString(7)
		log.Println("SECRET", SECRET)
		var Uaes = uaes.NewAES(SECRET)
		resP, _ := this.BaseGetAll(param, this.dbUtil)
		resEn, _ := Uaes.Encrypt_Any(resP)
		return data, resEn + SECRET, metadata
	} else {
		res, data := this.BaseGetAll(param, this.dbUtil)
		return res, "", data
	}
	// return this.BaseGetAll(param, this.dbUtil)
}

func (this *UserService) GetOne(key, value string) (res model.User_View, resEn string, errMessage string) {
	this.ChangeCollectionName("user")
	if os.Getenv("PROD_MODE") == "true" {
		SECRET := secret.GenerateRandomString(7)
		log.Println("SECRET", SECRET)
		var Uaes = uaes.NewAES(SECRET)
		this.dbUtil.FindOne(key, value, &res)
		resEn, _ = Uaes.Encrypt_Any(res)
		return res, resEn + SECRET, errMessage
	} else {
		this.dbUtil.FindOne(key, value, &res)
		return
	}
}

func (this *UserService) getValidUsername(username string) string {
	this.ChangeCollectionName("user")
	return strings.ToLower(username)
}

func (this *UserService) preparePassword(param *model.User) (err error) {
	if param.Password == "" {
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), this.passwordCost)
	if err != nil {
		log.Println(err)
		return errors.New("Fail to hash password")
	}

	param.Password = string(hashPassword)
	return
}

// func (this *UserService) preparePasswords(param *model.User_Profil) (err error) {
// 	if param.Password == "" {
// 		return
// 	}

// 	hashPassword, err := bcrypt.GenerateFromPassword([]byte(param.Password), this.passwordCost)
// 	if err != nil {
// 		log.Println(err)
// 		return errors.New("Fail to hash password")
// 	}

// 	param.Password = string(hashPassword)
// 	return
// }

func (this *UserService) UpsertWithHashingPassword(param model.User, isUpdate bool) (resp model.Response) {
	param.Username = this.getValidUsername(param.Username)

	if err := this.preparePassword(&param); err != nil {
		resp.Metadata.Message = err.Error()
		return
	}

	if err := this.dbUtil.CheckDuplicate(param.ID, []bson.M{
		{"username": param.Username, "status": bson.M{"$ne": "archive"}},
	}); err != nil {
		resp.Metadata.Message = err.Error()
		return
	}

	upsertErr, upsertId := this.dbUtil.UpsertAndGetId(isUpdate, &param)
	resp.Metadata.Message = upsertErr
	resp.Data = model.Response_Data_Upsert{
		ID: upsertId,
	}

	return
}

func (this *UserService) UpsertWithHashingPasswordAdmin(param model.User, isUpdate bool) (resp model.Response) {

	upsertErr, upsertId := this.dbUtil.UpsertAndGetId(isUpdate, &param)
	resp.Metadata.Message = upsertErr
	resp.Data = model.Response_Data_Upsert{
		ID: upsertId,
	}

	return
}

func (this *UserService) UpsertWithHashPassword(param model.User_Profil, isUpdate bool) (resp model.Response) {
	this.ChangeCollectionName("user")
	param.Username = this.getValidUsername(param.Username)

	if err := this.dbUtil.CheckDuplicate(param.ID, []bson.M{
		{"username": param.Username},
	}); err != nil {
		resp.Metadata.Message = err.Error()
		return
	}

	upsertErr, upsertId := this.dbUtil.UpsertAndGetId(isUpdate, &param)
	resp.Metadata.Message = upsertErr
	resp.Data = model.Response_Data_Upsert{
		ID: upsertId,
	}

	return
}

func (this *UserService) DeleteOne(key, value string) (errMessage string) {
	errMessage = this.dbUtil.DeleteOne(key, value)
	return
}

func (this *UserService) GetUserWithValidatePassword(key, value, passwordToTest string, ptrUser *model.User) (resp model.Response, ok bool) {
	if err := this.dbUtil.BaseFindOne(bson.M{key: value, "status": bson.M{"$ne": "archive"}}, &ptrUser); err != nil {
		log.Println(err)
		resp.Metadata.Message = "Data not found"
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(ptrUser.Password), []byte(passwordToTest)); err != nil {
		log.Println(err)
		resp.Metadata.Message = "Password doesn't match"
		return
	}

	//? Status validation
	if ptrUser.Status != enum.UserStatus_active.String() {
		err := errors.New("User is not active")
		log.Println(err)
		resp.Metadata.Message = err.Error()
		return
	}

	//? Don't show user password
	ptrUser.Password = ""
	ok = true
	return
}

func (this *UserService) ResetPassword(param model.User_ResetPassword) (resp model.Response) {
	findRes := model.User{}
	resp, ok := this.GetUserWithValidatePassword("_id", param.ID, param.OldPassword, &findRes)
	if !ok {
		return
	}

	findRes.Password = param.NewPassword
	if upsertRes := this.UpsertWithHashingPassword(findRes, true); upsertRes.Metadata.Message != "" {
		resp.Metadata.Message = upsertRes.Metadata.Message
		return
	}

	return
}

func (this *UserService) ResetPasswordAdmin(param model.User_ResetPassword_Admin) (resp model.Response) {
	findRes := model.User{}

	findRes.Password = param.NewPassword
	if upsertRes := this.UpsertWithHashingPasswordAdmin(findRes, true); upsertRes.Metadata.Message != "" {
		resp.Metadata.Message = upsertRes.Metadata.Message
		return
	}

	return
}
