package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type MetadataResponse struct {
	Status        bool   `json:"status"`
	Message       string `json:"message"`
	TimeExecution string `json:"timeExecution"`

	Pagination *PaginationResponse `json:"pagination" bson:"-"`
}

type Response struct {
	Metadata MetadataResponse `json:"metaData"`
	Data     interface{}      `json:"data"`
}

type ResponseUser struct {
	Metadata MetadataResponse `json:"metaData"`
	Data     map[string]interface{}      `json:"data"`
}

type Response_Data_Upsert struct {
	ID string `json:"id"`
}

type Response_Data_Encrypt struct {
	Data string `json:"data"`
}

type Request struct {
	Request_Pagination
	Request_Search
}

type Range struct {
	Field string `json:"field" example:"updatedAt"`
	Start int64  `json:"start" example:"1646792565000"`
	End   int64  `json:"end" example:"1646792565000"`
}

type Request_Search struct {
	Range *Range `json:"range"`
}

func (this Request_Search) BaseHandle(filter bson.M, rangeField string) (res bson.M) {
	if requestRange := this.Range; requestRange != nil && requestRange.Field != "" {
		rangeField = this.Range.Field
	}
	if rangeField == "" {
		rangeField = "updatedAt"
	}
	res = filter

	if this.Range != nil {
		if this.Range.Start == 0 && this.Range.End == 0 {
			timeNow := time.Now()
			this.Range.End = timeNow.UnixMilli()
			this.Range.Start = timeNow.AddDate(0, 0, -7).UnixMilli()
		}
		filter[rangeField] = bson.M{
			"$gte": this.Range.Start, "$lt": this.Range.End,
		}
	}

	return
}
func (this Request_Search) Handle_RequestSearch(filter bson.M) (res bson.M) {
	return this.BaseHandle(filter, "")
}

type Metadata struct {
	CreatedAt int64 `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt int64 `json:"updatedAt" bson:"updatedAt,omitempty"`
}
type MetadataWithID struct {
	ID       string `json:"_id" bson:"_id"`
	Metadata `bson:",inline"`
}

type MetadataLower struct {
	CreatedAt int64 `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at,omitempty"`
}

type MetadataWithIDLower struct {
	ID       string `json:"id" bson:"id"`
	MetadataLower `bson:",inline"`
}

type Request_Pagination struct {
	OrderBy string `json:"orderBy" example:"createdAt"`
	Order   string `json:"order" example:"DESC"`

	Page int64 `example:"1" json:"page"`
	Size int64 `example:"11" json:"size"`
}

type PaginationResponse struct {
	Size          int   `json:"size"`
	TotalPages    int64 `json:"totalPages"`
	TotalElements int64 `json:"totalElements"`
}
