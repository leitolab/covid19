package models

import (
	"context"
	"errors"
	"ieliot/src/common"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Place Lugar almacenado en el mapa
type Place struct {
	ID       string    `json:"_id" bson:"_id,omitempty"`     // _id de mongo
	Device   string    `json:"d"  bson:"d"`                  // producto al cual pertenece el cliente
	Accuracy float64   `json:"accuracy,omitempty"  bson:"-"` // producto al cual pertenece el cliente
	Coor     []float64 `json:"coor"  bson:"coor"`            // producto al cual pertenece el cliente
	Data     bson.M    `json:"data"  bson:"data"`            // data del cliente
	Created  time.Time `json:"t,omitempty" bson:"t"`         // fecha de creacion de la inserción
}

// FindNear obtenemos los lugares cercanos
// db.places.createIndex({ "coor" : "2dsphere" })
func (place *Place) FindNear() ([]Place, error) {

	if len(place.Coor) != 2 {
		return nil, errors.New("params are required: coor")
	}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := Config{}
	config.GetConfig()

	t := time.Now().UTC()
	td := t.AddDate(0, 0, -3) // 72 horas de vida del virus

	query := []bson.M{
		bson.M{
			"$geoNear": bson.M{
				"near": bson.M{
					"coordinates": place.Coor,
				},
				"distanceField": "data.calculated",
				"maxDistance":   place.Accuracy,
				"spherical":     true,
			}},
		bson.M{"$match": bson.M{"t": bson.M{"$gte": td}}},
		bson.M{"$limit": 100},
	}

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("places")
	cur, err := collection.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// obtenemos los ids de contacto
	var places []Place
	for cur.Next(ctx) {
		var mPlace Place
		err := cur.Decode(&mPlace)
		if err != nil {
			return nil, err
		}
		places = append(places, mPlace)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return places, nil
}

// InsertOne Insertamos en la base de datos un place
func (place *Place) InsertOne() error {

	if len(place.Coor) != 2 {
		return errors.New("params are required: coor")
	}
	if place.Data == nil {
		return errors.New("params are required: data")
	}

	place.Created = time.Now().UTC()

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("places")
	res, err := collection.InsertOne(ctx, place)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		place.ID = oid.Hex()
	}
	return nil
}

// DeleteOne Eliminamos de la base de datos un elemento
func (place *Place) DeleteOne() error {
	// Parseamos el id a una estructura de id de mongo
	placeID, err := primitive.ObjectIDFromHex(place.ID)
	if err != nil {
		return err
	}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("places")
	res, err := collection.DeleteOne(ctx, bson.M{"_id": placeID, "d": place.Device})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("The document to delete was not found")
	}

	return nil
}
