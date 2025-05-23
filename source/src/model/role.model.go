package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Generated by https://quicktype.io

type Role struct {
	MetadataWithID `bson:",inline"`

	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Privileges  []map[string]interface{} `json:"privileges"`
}

type Role_View struct {
	Role `bson:",inline"`
	UserCount int `json:"userCount,omitempty" bson:"userCount,omitempty"`
}

type Role_Search struct {
	//? Regex
	Search   string                   `json:"search"`
	SearchBy []string                 `json:"searchBy"`
	Filter   []map[string]interface{} `json:"filter"`

	Request
}

func (this *Role_Search) HandleFilter(listFilterAnd *[]bson.M) {
	if search := this.Search; search != "" {
		filterOr := bson.M{}
		listFilterOr := []bson.M{}
		for _, searchField := range this.SearchBy {
			listFilterOr = append(listFilterOr, bson.M{searchField: primitive.Regex{Pattern: search, Options: "i"}})
		}

		if len(listFilterOr) > 0 {
			filterOr["$or"] = listFilterOr
		}
		*listFilterAnd = append(*listFilterAnd, filterOr)
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
