package models

import (
	"context"
	"errors"
	"ieliot/src/common"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Symptom estructura del dispositivo a nivel de base de datos pero tambien a nivel de JWT
// Estado de infección en data.infected 0: no infectado, 1: infectado, 2 recuperado, 3: muerto
type Symptom struct {
	Data    []int64    `json:"symptom",bson:"symptom"` // sintomas del usuario
	Created *time.Time `json:"t",bson:"t"`             // fecha de reporte de los mismos
}

// UpdateOne ...
func (symptom *Symptom) UpdateOne(id string) error {
	// Parseamos el id a una estructura de id de mongo
	deviceID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": deviceID}
	update := bson.M{"$set": bson.M{"symptom_last": symptom.Data, "symptom_t": symptom.Created}, "$addToSet": bson.M{"symptom": symptom}}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("devices")
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("The document to update was not found")
	}

	return nil
}
