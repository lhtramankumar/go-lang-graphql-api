package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type BookListingDB struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Bookname    string             `json:"bookname"`
	Description string             `bson:"description"`
	Author      string             `bson:"author"`
	AddedOn     float64            `json:"addedOn"`
}

type BookWithID struct {
	ID          primitive.ObjectID `json:"_id"`
	Title       string             `json:"title"`
	Bookname    string             `json:"bookname"`
	Description string             `json:"description"`
	Author      string             `json:"author"`
	AddedOn     float64            `json:"addedOn"`
	// Other fields...
}
