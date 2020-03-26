package models

import (
	"context"
	"ieliot/src/common"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Contact ...
type Contact struct {
	Coor []float64 `json:"coor,omitempty" bson:"coor"`
	T    time.Time `json:"t,omitempty" bson:"t"` // fecha de creacion de la inserción
}

type contactB struct {
	ID string `bson:"b"`
}

type contactID struct {
	ID bson.ObjectId `bson:"_id"`
}

// GetContactIds ...
func (contact *Contact) GetContactIds(device string) ([]primitive.ObjectID, error) {
	var err error
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la consulta a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("contacts")

	config := Config{}
	config.GetConfig()

	t := time.Now()
	td := t.AddDate(0, 0, config.Delta)

	cur, err := collection.Find(ctx, bson.M{"a": device, "t": bson.M{"$gte": td}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var ids []primitive.ObjectID
	for cur.Next(ctx) {
		// una variable para cada contacto encontrado
		var mContact contactB
		err := cur.Decode(&mContact)
		if err != nil {
			return nil, err
		}
		// almacenamos los ids
		deviceID, err := primitive.ObjectIDFromHex(mContact.ID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, deviceID)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

// GetInfected ...
func (contact *Contact) GetInfected(devices *[]primitive.ObjectID) (int64, error) {

	var err error
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": (*devices)}, "data.infected": 1}
	// se ejecuta el conteo a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("devices")
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
