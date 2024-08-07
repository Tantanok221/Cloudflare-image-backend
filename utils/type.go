package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

type Path struct {
	ImageLink string      `json:"image-link" bson:"image_link"`
	Name      string      `json:"author_name" bson:"author_name"`
	Paths     primitive.M `json:"path" bson:"path"`
}
