package models

import (
	"context"
	"ieliot/src/common"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Config ...
type Config struct {
	Accuracy float64 `json:"accuracy,required" bson:"accuracy"` // data del cliente
	Delta    int     `json:"delta,required" bson:"delta"`       // producto al cual pertenece el cliente
}

// GetConfig ...
func (config *Config) GetConfig() error {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("config")
	if err := collection.FindOne(ctx, bson.M{"id": 1}).Decode(config); err != nil {
		return err
	}

	return nil
}
