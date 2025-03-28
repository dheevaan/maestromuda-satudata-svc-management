package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Generated by https://quicktype.io

type Catalog struct {
	MetadataWithID `bson:",inline"`

	Name     string `json:"name" bson:"name"`
	Category string `json:"category" bson:"category"`
}

type Catalog_Search struct {
	//? Regex
	Search string                   `json:"search"`
	Filter []map[string]interface{} `json:"filter"`

	Request
}

type Catalog_View struct {
	Catalog `bson:",inline"`
}

func (this *Catalog_Search) HandleFilter(listFilterAnd *[]bson.M) {
	if search := this.Search; search != "" {
		*listFilterAnd = append(*listFilterAnd, bson.M{"name": primitive.Regex{Pattern: search, Options: "i"}})
	}
	if filter := this.Filter; len(filter) > 0 {
		filterAnd := bson.M{}
		filtersAnd := []bson.M{}
		for _, filterOpt := range filter {
			filtersAnd = append(filtersAnd, bson.M{filterOpt["field"].(string): primitive.Regex{Pattern: filterOpt["value"].(string), Options: "i"}})
		}

		if len(filtersAnd) > 0 {
			filterAnd["$and"] = filtersAnd
		}
		*listFilterAnd = append(*listFilterAnd, filterAnd)
	}
}
