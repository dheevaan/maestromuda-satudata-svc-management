package db

import (
	"data-management/src/model"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (this MongoDbUtil) AggsPagination(pipeline []bson.M, requestPagination model.Request_Pagination, pointerDecodeTo interface{}) (
	paginationResp *model.PaginationResponse, err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	//* -------------------------------- PIPELINE -------------------------------- */
	skip, limit := GetSkipAndLimit(requestPagination)

	// `[{"$group":{"_id":"$user.divisionId","totalHour":{"$sum":"$totalHour"}}},{"$sort":{"totalHour":-1}},{"$skip":4},{"$limit":2}]`
	pipeline = append(pipeline, []bson.M{
		// {"$sort": bson.M{"totalHour": GetSortValue(requestPagination)}}, //? Just example
	}...)

	pipelineWithoutPagination := pipeline
	// //* ---------------------------------- COUNT --------------------------------- */
	count, err := col.Aggregate(this.ctx, pipelineWithoutPagination, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		log.Println(err)
		asJson, _ := json.Marshal(pipelineWithoutPagination)
		log.Println(string(asJson))
		return
	}
	countValue := 0

	countRes := []bson.M{}
	if err = count.All(this.ctx, &countRes); err != nil {
		if err == io.EOF {
			countValue = 0
		} else {
			log.Println(err)
			query, _ := json.Marshal(pipelineWithoutPagination)
			fmt.Printf("query: \n %s\n", query)
			return
		}
	} else {
		countValue = len(countRes)
	}

	//* --------------------------------- EXECUTE PAGINATION -------------------------------- */
	pipelineWithPagination := append(pipeline, []primitive.M{
		{"$skip": skip},
		{"$limit": limit},
	}...)
	res, err := col.Aggregate(this.ctx, pipelineWithPagination, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		log.Println(err)
		return
	}

	if err = res.All(this.ctx, pointerDecodeTo); err != nil {
		log.Println(err)
		return
	}

	//* --------------------------- SET PAGINATION RESP -------------------------- */
	paginationResp = &model.PaginationResponse{
		Size:          int(limit),
		TotalElements: int64(countValue),
		TotalPages:    int64(math.Ceil(float64(countValue) / float64(limit))),
	}

	if countValue == 0 {
		err = errors.New("No data found.")
	}
	return
}

func (this MongoDbUtil) AggsCount(widget string, pipeline []bson.M, pointerDecodeTo interface{}) (
	paginationResp interface{}, err error) {
	client, err := this.Connect()
	if err != nil {
		return
	}
	defer this.Disconnect(client)
	col := client.Database(this.dbName).Collection(this.collectionName)

	//* -------------------------------- PIPELINE -------------------------------- */

	// `[{"$group":{"_id":"$user.divisionId","totalHour":{"$sum":"$totalHour"}}},{"$sort":{"totalHour":-1}},{"$skip":4},{"$limit":2}]`
	pipeline = append(pipeline, []bson.M{
		// {"$sort": bson.M{"totalHour": GetSortValue(requestPagination)}}, //? Just example
	}...)

	// //* ---------------------------------- COUNT --------------------------------- */
	count, err := col.Aggregate(this.ctx, pipeline, &options.AggregateOptions{})
	if err != nil {
		log.Println(err)
		asJson, _ := json.Marshal(pipeline)
		log.Println(string(asJson))
		return
	}
	countValue := 0

	countRes := []bson.M{}
	if err = count.All(this.ctx, &countRes); err != nil {
		if err == io.EOF {
			countValue = 0
		} else {
			log.Println(err)
			query, _ := json.Marshal(pipeline)
			fmt.Printf("query: \n %s\n", query)
			return
		}
	} else {
		countValue = len(countRes)
	}

	//* --------------------------- SET PAGINATION RESP -------------------------- */
	paginationResp = map[string]interface{}{
		"key":   widget,
		"value": countValue,
	}

	if countValue == 0 {
		err = errors.New("No data found.")
	}
	return
}
