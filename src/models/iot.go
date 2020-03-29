package models

import (
	"context"
	"errors"
	"ieliot/src/common"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Iot estructura para manejar los envios de datos de los dispositivos
type Iot struct {
	ID      string      `json:"_id" bson:"_id,omitempty"` // _id de mongo
	Data    interface{} `json:"data"  bson:"data"`        // data del cliente
	Device  string      `json:"d"  bson:"d"`              // producto al cual pertenece el cliente
	Client  string      `json:"c"  bson:"c,omitempty"`    // producto al cual pertenece el cliente
	Created time.Time   `json:"t,omitempty" bson:"t"`     // fecha de creacion de la inserción
}

// Upsert actualización de la data de iot en mongo
func (iot *Iot) Upsert(product string) error {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// buscamos el dispositivo que pertenece al grupo cliente
	filter := bson.M{"d": iot.Device, "c": iot.Client}
	// generamos la petición para actualizar, la nueva posición del dispositivo
	update := bson.M{"$set": iot}

	// se ejecuta el upsert para el mapa de localización
	collection := common.Client.Database(common.DATABASE).Collection(product)
	res, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	// se elimina el cliente para optimizar espacio pero se almacena para restituir en caso de
	// posteriormente se requiera esta petición
	client := iot.Client
	iot.Client = ""
	// se ejecuta la inserción a las localizaciones
	collection = common.Client.Database(common.DATABASE).Collection("locations")
	if _, err := collection.InsertOne(ctx, iot); err != nil {
		return err
	}
	iot.Client = client

	// si no se insertó, modifico o actualizó nada en la base de datos generamos el error
	if res.MatchedCount == 0 && res.ModifiedCount == 0 && res.UpsertedCount == 0 {
		return errors.New("upsert fail")
	}

	return nil
}

// Near obtenemos los devices cercanos y actualizamos los contactos
// db.covid19.createIndex({ "data.coor" : "2dsphere" })
// db.locations.createIndex({ "data.coor" : "2dsphere" })
func (iot *Iot) Near(product string) ([]Iot, error) {
	// data empaquetada que envió el usuario necesitamos las coordenadas para determinar
	// las personas cercanas de la tabla de usuarios
	data := iot.Data.(map[string]interface{})
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config := Config{}
	config.GetConfig()

	// la última posición emitida del usuario es tenida en cuenta hasta por 72 horas
	t := time.Now()
	td := t.AddDate(0, 0, -3) // 72 horas de vida del virus

	// filtro a la base de datos para buscar las personas cercanas
	query := []bson.M{
		bson.M{
			"$geoNear": bson.M{
				"near": bson.M{
					"coordinates": data["coor"],
				},
				"distanceField": "data.calculated",
				"maxDistance":   config.Accuracy,
				"spherical":     true,
			}},
		bson.M{
			"$match": bson.M{
				"d": bson.M{"$ne": iot.Device},
				"t": bson.M{"$gte": td},
			}},
	}

	// se ejecuta la solicitud de los datos de las personas cercanas
	collection := common.Client.Database(common.DATABASE).Collection(product)
	cur, err := collection.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// generamos un array con todos os contactos del tic recorriendo un cursor
	var iots []Iot
	for cur.Next(ctx) {
		var mIot Iot
		err := cur.Decode(&mIot)
		if err != nil {
			return nil, err
		}
		iots = append(iots, mIot)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	// regresamos las personas que hicieron contacto
	return iots, nil
}

// InsertContact actualizamos la última fecha de contacto con un usuario dado el tic
func (iot *Iot) InsertContact(product string, iots *[]Iot) {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se intenta ejecutar el upsert a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection(product)
	for _, iotb := range *iots {
		filter := bson.M{"a": iot.Device, "b": iotb.Device}
		update := bson.M{"a": iot.Device, "b": iotb.Device, "t": time.Now().UTC(), "coor_a": iot.Data, "coor_b": iotb.Data}
		collection.UpdateOne(ctx, filter, bson.M{"$set": update}, options.Update().SetUpsert(true))
	}
}
