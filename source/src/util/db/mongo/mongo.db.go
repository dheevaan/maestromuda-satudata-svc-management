package db

import (
	"context"
	"data-management/src/config"
	"data-management/src/model"
	"data-management/src/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"time"

	"github.com/logrusorgru/aurora"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
)

var Mongo_LiveConnection int

type MongoDbUtil struct {
	srv            string
	dbName         string
	collectionName string
	ctx            context.Context
}

func NewMongoDbUtil(srv, dbName, collectionName string) *MongoDbUtil {
	return &MongoDbUtil{
		srv:            srv,
		dbName:         dbName,
		collectionName: collectionName,
		ctx:            context.Background(),
	}
}
func NewMongoDbUtilLocal(collectionName string) *MongoDbUtil {
	return &MongoDbUtil{
		srv:            "mongodb://localhost:27017",
		dbName:         os.Getenv(config.ENV_KEY_MONGO_DB),
		collectionName: collectionName,
		ctx:            context.Background(),
	}
}
func NewMongoDbUtilUseEnv(collectionName string) *MongoDbUtil {
	return &MongoDbUtil{
		srv:            os.Getenv(config.ENV_KEY_MONGO_SRV),
		dbName:         os.Getenv(config.ENV_KEY_MONGO_DB),
		collectionName: collectionName,
		ctx:            context.Background(),
	}
}

func (this MongoDbUtil) Connect() (client *mongo.Client, err error) {
	clientOptions := options.Client()
	clientOptions.ApplyURI(this.srv)

	client, err = mongo.Connect(this.ctx, clientOptions)
	if err != nil {
		log.Println(err)
		return
	}

	Mongo_LiveConnection += 1
	return
}
func (this MongoDbUtil) Disconnect(client *mongo.Client) {
	if client == nil {
		return
	}

	if err := client.Disconnect(this.ctx); err != nil {
		log.Println(err)
		return
	}

	Mongo_LiveConnection -= 1
}

func (this MongoDbUtil) Upsert(isUpdate bool, ptrParam interface{}) (errMessage string) {
	errMessage, _ = this.UpsertAndGetId(isUpdate, ptrParam)
	return
}

func (this MongoDbUtil) UpdateMany(filter, data interface{}) (updateRes *mongo.UpdateResult, err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)
	updateRes, err = col.UpdateMany(context.Background(), filter, data)
	return
}

func (this MongoDbUtil) UpdateOne(filter, data interface{}) (updateRes *mongo.UpdateResult, err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)
	updateRes, err = col.UpdateOne(context.Background(), filter, data)
	return
}

func (this MongoDbUtil) UpsertAndGetId(isUpdate bool, ptrParam interface{}) (errMessage string, newDataId string) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	paramAsReflect := reflect.ValueOf(ptrParam).Elem()
	if updatedAt := paramAsReflect.FieldByName("UpdatedAt"); updatedAt.IsValid() {
		updatedAt.SetInt(time.Now().UnixMilli())
	}
	idField := paramAsReflect.FieldByName("ID")

	if !isUpdate {
		if createdAt := paramAsReflect.FieldByName("CreatedAt"); createdAt.IsValid() {
			createdAt.SetInt(time.Now().UnixMilli())
		}
		if idField.IsValid() && idField.String() == "" {
			idField.SetString(util.GenerateID())
		}

		if insertRes, err := col.InsertOne(this.ctx, ptrParam); err != nil {
			log.Println(err)
			return
		} else {
			newDataId = fmt.Sprint(insertRes.InsertedID)
		}
	} else {
		if !idField.IsValid() {
			log.Println("Fail to get field ID")
			return
		}

		id := fmt.Sprint(paramAsReflect.FieldByName("ID"))
		if updateRes, err := col.UpdateByID(this.ctx, id, bson.M{"$set": ptrParam}); err != nil {
			log.Println(err)
			return
		} else {
			if updateRes.MatchedCount == 0 {
				errMessage = "Data not found."
				return
			}
			newDataId = id
		}
	}

	return
}

func (this MongoDbUtil) UpsertMap(isUpdate bool, ptrParam interface{}) (errMessage string) {
	errMessage, _ = this.UpsertAndGetIdMap(isUpdate, ptrParam)
	return
}

func (this MongoDbUtil) UpsertAndGetIdMap(isUpdate bool, ptrParam interface{}) (errMessage string, newDataId string) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	paramAsReflect := reflect.ValueOf(ptrParam).Elem()
	if updatedAt := paramAsReflect.MapIndex(reflect.ValueOf("updated_at")); updatedAt.IsValid() {
		paramAsReflect.SetMapIndex(reflect.ValueOf("updated_at"), reflect.ValueOf(time.Now().UnixMilli()))
	}
	idField := paramAsReflect.MapIndex(reflect.ValueOf("_id"))
	if !isUpdate {
		paramAsReflect.SetMapIndex(reflect.ValueOf("created_at"), reflect.ValueOf(time.Now().UnixMilli()))
		if !idField.IsValid() {
			paramAsReflect.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(util.GenerateID()))
		}

		if insertRes, err := col.InsertOne(this.ctx, ptrParam); err != nil {
			log.Println(err)
			return
		} else {
			newDataId = fmt.Sprint(insertRes.InsertedID)
		}
	} else {
		if !idField.IsValid() {
			log.Println("Fail to get field ID")
			return
		}

		id := fmt.Sprint(paramAsReflect.MapIndex(reflect.ValueOf("_id")))
		if updateRes, err := col.UpdateByID(this.ctx, id, bson.M{"$set": ptrParam}); err != nil {
			log.Println(err)
			return
		} else {
			if updateRes.MatchedCount == 0 {
				errMessage = "Data not found."
				log.Println(errMessage)
				return
			}
			newDataId = id
		}
	}

	return
}

func (this MongoDbUtil) UpsertAndGetIdForm(isUpdate bool, ptrParam primitive.D) (errMessage string, newDataId string) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)
	idField := false
	idIndex := -1

	for i, k := range ptrParam {
		if k.Key == "_id" {
			idField = true
			idIndex = i
		}
	}

	if !isUpdate {
		ptrParam = append(ptrParam, primitive.E{"created_at", time.Now().UnixMilli()})

		if !idField {
			ptrParam = append(ptrParam, primitive.E{"_id", util.GenerateID()})
			// paramAsReflect.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(util.GenerateID()))
		}

		if insertRes, err := col.InsertOne(this.ctx, ptrParam); err != nil {
			log.Println(err)
			return
		} else {
			newDataId = fmt.Sprint(insertRes.InsertedID)
		}
	} else {
		if !idField {
			log.Println("Fail to get field ID")
			return
		}

		ptrParam = append(ptrParam, primitive.E{"updated_at", time.Now().UnixMilli()})

		id := ptrParam[idIndex].Value.(string)
		if updateRes, err := col.ReplaceOne(this.ctx, bson.M{"_id": id}, ptrParam); err != nil {
			log.Println(err)
			return
		} else {
			if updateRes.MatchedCount == 0 {
				errMessage = "Data not found."
				log.Println(errMessage)
				return
			}
			newDataId = id
		}
	}
	return
}

func (this MongoDbUtil) CountData(widget string, filter bson.M, pointerDecodeTo interface{}) (res int64, err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	// countValue := 0

	res, err = col.CountDocuments(this.ctx, filter)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (this MongoDbUtil) BaseFindOne(filter bson.M, pointerDecodeTo interface{}) (err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	log.Printf("this: %+v\n", this)

	res := col.FindOne(this.ctx, filter)

	log.Println("++++++++", res.Err())

	if res.Err() != nil {
		fmt.Printf("this: %+v\n", this)
		filterAsJson, _ := json.Marshal(filter)
		log.Println(err, this.collectionName, string(filterAsJson))
		err = errors.New("Data not found.")
		return
	}

	if err := res.Decode(pointerDecodeTo); err != nil {
		log.Println(err)
	}
	return
}

func (this MongoDbUtil) BaseFindOneUser(filter bson.M, pointerDecodeTo interface{}) (err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	res := col.FindOne(this.ctx, filter, options.FindOne().SetProjection(bson.D{{"password", 0}}))
	if res.Err() != nil {
		fmt.Printf("this: %+v\n", this)
		filterAsJson, _ := json.Marshal(filter)
		log.Println(err, this.collectionName, string(filterAsJson))
		err = errors.New("Data not found.")
		return
	}

	if err := res.Decode(pointerDecodeTo); err != nil {
		log.Println(err)
	}
	return
}

func (this MongoDbUtil) BaseFindOnePenduduk(filter bson.M, pointerDecodeTo interface{}) (err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	res := col.FindOne(this.ctx, filter, options.FindOne().SetProjection(bson.D{{"workspaceId", 0}, {"updated_at", 0}, {"status", 0}, {"edited_by", 0}, {"act_by_username", 0}, {"act_by", 0}, {"workspace", 0}, {"umum.jenis_perubahan", 0}, {"lain_lain.keterangan", 0}, {"umum.penambahan", 0}, {"umum.pengurangan", 0}, {"user_id", 0}, {"umum.kampung", 0}}))
	if res.Err() != nil {
		fmt.Printf("this: %+v\n", this)
		filterAsJson, _ := json.Marshal(filter)
		log.Println(err, this.collectionName, string(filterAsJson))
		err = errors.New("Data not found.")
		return
	}

	if err := res.Decode(pointerDecodeTo); err != nil {
		log.Println(err)
	}
	return
}

func (this MongoDbUtil) FindOne(key, value string, pointerDecodeTo interface{}) (err error) {
	return this.BaseFindOne(bson.M{key: value}, pointerDecodeTo)
}

func (this MongoDbUtil) FindOneSdgsNik(key, value string, pointerDecodeTo interface{}) (err error) {
	return this.BaseFindOne(bson.M{key: value, "status": bson.M{"$ne": "archive"}}, pointerDecodeTo)
}

func (this MongoDbUtil) FindOneTTD(key1, value1, key2, value2 string, pointerDecodeTo interface{}) (err error) {
	return this.BaseFindOne(bson.M{key1: value1, key2: value2}, pointerDecodeTo)
}

func (this MongoDbUtil) FindOnePenduduk(key, value string, pointerDecodeTo interface{}) (err error) {
	return this.BaseFindOnePenduduk(bson.M{key: value}, pointerDecodeTo)
}
func (this MongoDbUtil) FindOneUser(key, value string, pointerDecodeTo interface{}) (err error) {
	return this.BaseFindOneUser(bson.M{key: value}, pointerDecodeTo)
}

func (this MongoDbUtil) BaseFindOneMap(filter bson.M) (result interface{}, err error) {
	client, err := this.Connect()
	var ress bson.D
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	res := col.FindOne(this.ctx, filter)
	if res.Err() != nil {
		fmt.Printf("this: %+v\n", this)
		filterAsJson, _ := json.Marshal(filter)
		log.Println(err, this.collectionName, string(filterAsJson))
		err = errors.New("Data not found.")
		return
	}

	err = res.Decode(&ress)

	result = ress
	if err != nil {
		log.Println(err)
	}

	return
}

func (this MongoDbUtil) BaseFindOneMapSUrat(filter bson.M) (result interface{}, err error) {
	client, err := this.Connect()
	var ress bson.D
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	res := col.FindOne(this.ctx, filter, options.FindOne().SetProjection(bson.D{{"workspaceId", 0}, {"status", 0}}))
	if res.Err() != nil {
		fmt.Printf("this: %+v\n", this)
		filterAsJson, _ := json.Marshal(filter)
		log.Println(err, this.collectionName, string(filterAsJson))
		err = errors.New("Data not found.")
		return
	}

	err = res.Decode(&ress)

	result = ress
	if err != nil {
		log.Println(err)
	}

	return
}

func (this MongoDbUtil) BaseFindOneMapInduk(filter bson.M) (result interface{}, err error) {
	client, err := this.Connect()
	var ress bson.D
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	res := col.FindOne(this.ctx, filter, options.FindOne().SetProjection(bson.D{{"workspaceId", 0}, {"status", 0}, {"kategori_umur", 0}, {"edited_by", 0}, {"act_by_username", 0}, {"act_by", 0}, {"workspace", 0}, {"umum.jenis_perubahan", 0}, {"lain_lain.keterangan", 0}, {"umum.penambahan", 0}, {"umum.pengurangan", 0}, {"user_id", 0}, {"umum.kampung", 0}}))

	if res.Err() != nil {
		fmt.Printf("this: %+v\n", this)
		filterAsJson, _ := json.Marshal(filter)
		log.Println(err, this.collectionName, string(filterAsJson))
		err = errors.New("Data not found.")
		return
	}

	err = res.Decode(&ress)

	result = ress
	if err != nil {
		log.Println(err)
	}

	return
}

func (this MongoDbUtil) FindOneMap(key, value string) (result interface{}, err error) {
	return this.BaseFindOneMap(bson.M{key: value})
}

func (this MongoDbUtil) FindOneMapSurat(key, value string) (result interface{}, err error) {
	return this.BaseFindOneMapSUrat(bson.M{key: value})
}

func (this MongoDbUtil) FindOneMapInduk(key, value string) (result interface{}, err error) {
	return this.BaseFindOneMapInduk(bson.M{key: value})
}

func (this MongoDbUtil) FindOneMapnik(key, value, key1, value1 string) (result interface{}, err error) {
	return this.BaseFindOneMap(bson.M{key: value, key1: value1})
}

func (this MongoDbUtil) FindOneMapWorkspace(key, value, key1 string, value1 interface{}) (result interface{}, err error) {
	return this.BaseFindOneMap(bson.M{key: value, key1: value1})
}

func (this MongoDbUtil) baseFindByCol(col *mongo.Collection, filter bson.M, opts *options.FindOptions, ptrDecodeTo interface{}) {
	findRes, err := col.Find(this.ctx, filter, opts)
	if err != nil {
		filterJson, _ := json.Marshal(filter)
		log.Println(string(filterJson))

		log.Println(err)
		return
	}

	if err := findRes.All(this.ctx, ptrDecodeTo); err != nil {
		log.Println(err)
		return
	}
	log.Println("ptorotgg")
}
func (this MongoDbUtil) BaseFind(filter bson.M, opts *options.FindOptions, ptrDecodeTo interface{}) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)
	this.baseFindByCol(col, filter, opts, ptrDecodeTo)
}
func (this MongoDbUtil) BaseFindPagination(filter bson.M,
	fieldRequestParam, pointerDecodeTo interface{}, rangeField string) (paginationResp *model.PaginationResponse, errMessage string) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	var requestPagination model.Request_Pagination
	//* ----------------------------- SET FILTER REQUEST ---------------------------- */
	switch requestAsType := fieldRequestParam.(type) {
	case model.Request:
		filter = requestAsType.BaseHandle(filter, rangeField)
		requestPagination = requestAsType.Request_Pagination
	case model.Request_Pagination:
		requestPagination = requestAsType
	}

	//* ------------------------- SET PAGINATION OPTIONS ------------------------- */
	//? Sorting
	findOptions := options.FindOptions{}
	if orderBy := requestPagination.OrderBy; orderBy != "" {
		findOptions.Sort = bson.M{requestPagination.OrderBy: GetSortValue(requestPagination)}
	}
	skip, limit := GetSkipAndLimit(requestPagination)
	findOptions.Skip = &skip
	findOptions.Limit = &limit

	//* ------------------------------- DEEP FILTER ------------------------------ */
	filter["status"] = bson.M{"$ne": "archive"} //? Delete operation will set data status to archive instead removing the data.
	//* ---------------------------------- FIND ---------------------------------- */
	this.baseFindByCol(col, filter, &findOptions, pointerDecodeTo)
	log.Println(pointerDecodeTo)

	//* ------------------------- SET RESPONSE PAGINATION ------------------------ */
	totalElements, err := col.CountDocuments(this.ctx, filter)
	if err != nil {
		log.Println(err)
		return
	}
	paginationResp = &model.PaginationResponse{
		Size:          int(limit),
		TotalElements: totalElements,
		TotalPages:    int64(math.Ceil(float64(totalElements) / float64(limit))),
	}

	if totalElements == 0 {
		errMessage = "No data found."
		filterAsJson, err := json.Marshal(filter)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("Filter[%s]\n%s\n", this.collectionName, aurora.Yellow(string(filterAsJson)))
	}
	return
}
func (this MongoDbUtil) Find(filter bson.M, fieldRequestParam, pointerDecodeTo interface{}) (paginationResp *model.PaginationResponse, errMessage string) {
	return this.BaseFindPagination(filter, fieldRequestParam, pointerDecodeTo, "")
}

type IfindWithRangeCustomField interface {
	GetRangeField() string
}

func (this MongoDbUtil) FindWithRangeCustomField(filter bson.M,
	fieldRequestParam, pointerDecodeTo interface{}, modelImplementTheInterface IfindWithRangeCustomField) (paginationResp *model.PaginationResponse, errMessage string) {
	return this.BaseFindPagination(filter, fieldRequestParam, pointerDecodeTo, modelImplementTheInterface.GetRangeField())
}

func (this MongoDbUtil) CheckDuplicate(id string, listFilterOr []bson.M) (err error) {
	var checkDuplicate bson.M
	_ = this.BaseFindOne(bson.M{"$or": listFilterOr}, &checkDuplicate)

	var oldData bson.M
	if id != "" {
		_ = this.FindOne("_id", id, &oldData)
	}

	for _, filterOr := range listFilterOr {
		for field, newDataValue := range filterOr {
			if newDataValue == checkDuplicate[field] {
				//* To allow update, check if the old data == new data
				//? but the new data not duplicate in another data
				if oldData[field] != newDataValue {
					return errors.New(field + " is already exists, and need to be unique")
				}
			}
		}
	}

	return
}

func (this MongoDbUtil) getSoftDeleteFlag() bson.M {
	return bson.M{"$set": bson.M{"status": "archive"}}
}

// ? Test command: curl -X 'DELETE' 'http://localhost:52306/api/v1/taskComment/delete?id=tws83946fdbfeeb4ef0a879fa561577f197' -H 'accept: application/json'
func (this MongoDbUtil) DeleteOne(key, value string) (errMessage string) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	filter := bson.M{key: value}
	res, err := col.UpdateOne(this.ctx, filter, this.getSoftDeleteFlag())
	if err != nil {
		log.Println(err)
		return
	}
	if res.ModifiedCount == 0 {
		errMessage = "Data not found."
		fmt.Printf("[%s] filter: %v\n", this.collectionName, filter)
	}

	return
}

func (this MongoDbUtil) Delete(filter bson.M) (errMessage string) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	res, err := col.UpdateMany(this.ctx, filter, this.getSoftDeleteFlag())
	if err != nil {
		log.Println(err)
		return
	}
	if res.ModifiedCount == 0 {
		errMessage = "Data not found."
		fmt.Printf("[%s] filter: %v\n", this.collectionName, filter)
	}

	return
}

func (this MongoDbUtil) CreateViewIfNotExists(viewName string, pipeline []bson.M) (err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	db := client.Database(this.dbName)
	if listCollectionName, err := db.ListCollectionNames(this.ctx, bson.M{}, &options.ListCollectionsOptions{}); err != nil {
		log.Println(err)
		return err
	} else {
		if slices.Contains(listCollectionName, viewName) {
			return nil
		}
	}

	if err := db.CreateView(this.ctx, viewName, this.collectionName, pipeline, &options.CreateViewOptions{}); err != nil {
		log.Println(viewName, err)
		return err
	}

	log.Printf("%s->%s created\n", viewName, this.collectionName)
	return
}
