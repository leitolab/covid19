package models

import (
	"context"
	"errors"
	"ieliot/src/common"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Iot ...
type Iot struct {
	ID      string      `json:"_id" bson:"_id,omitempty"` // _id de mongo
	Data    interface{} `json:"data"  bson:"data"`        // data del cliente
	Device  string      `json:"d"  bson:"d"`              // producto al cual pertenece el cliente
	Client  string      `json:"c"  bson:"c,omitempty"`    // producto al cual pertenece el cliente
	Created time.Time   `json:"t,omitempty" bson:"t"`     // fecha de creacion de la inserción
}

// Upsert actualizacion de la data de iot en mongo con la estructura prefijada
func (iot *Iot) Upsert(product string) error {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"d": iot.Device, "c": iot.Client}
	update := bson.M{"$set": iot}

	// se ejecuta el upsert para el mapa de localizacion
	collection := common.Client.Database(common.DATABASE).Collection(product)
	res, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	// se elimina el cliente para optimizar espacio pero se almacena para restituir en caso de
	// posteriormente se requiera
	client := iot.Client
	iot.Client = ""
	// se ejecuta la insercion a las localizaciones
	collection = common.Client.Database(common.DATABASE).Collection("locations")
	if _, err := collection.InsertOne(ctx, iot); err != nil {
		return err
	}
	iot.Client = client

	if res.MatchedCount == 0 && res.ModifiedCount == 0 && res.UpsertedCount == 0 {
		return errors.New("upsert fail")
	}

	return nil
}

type data struct {
	coor []float64 `json:"coor"  bson:"coor"`
}

// Near obtenemos los devices cercanos y actualizamos los contactos
// db.covid19.createIndex({ "data.coor" : "2dsphere" })
// db.locations.createIndex({ "data.coor" : "2dsphere" })
func (iot *Iot) Near(product string) ([]string, error) {
	data := iot.Data.(map[string]interface{})
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := Config{}
	config.GetConfig()

	query := []bson.M{
		bson.M{
			"$geoNear": bson.M{
				"near": bson.M{
					"type":        "Point",
					"coordinates": data["coor"],
				},
				"distanceField": "data.coor",
				"maxDistance":   config.Accuracy,
				"spherical":     true,
			}},
		bson.M{"$match": bson.M{"d": bson.M{"$ne": iot.Device}}},
	}

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection(product)
	cur, err := collection.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// obtenemos los ids de contacto
	var ids []string
	for cur.Next(ctx) {
		var mIot Iot
		err := cur.Decode(&mIot)
		if err != nil {
			return nil, err
		}
		ids = append(ids, mIot.Device)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

// Contact ...
func (iot *Iot) Contact(product string, ids *[]string) error {
	var err error
	dataIot := iot.Data.(map[string]interface{})
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection(product)
	for _, id := range *ids {
		filter := bson.M{"a": iot.Device, "b": id}
		update := bson.M{"a": iot.Device, "b": id, "t": time.Now().UTC(), "coor": dataIot["coor"]}
		_, err = collection.UpdateOne(ctx, filter, bson.M{"$set": update}, options.Update().SetUpsert(true))
	}
	if err != nil {
		return err
	}
	return nil
}
