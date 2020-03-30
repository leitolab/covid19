package models

import (
	"context"
	"ieliot/src/common"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Contact estructura para para almacenar los contactos
type Contact struct {
	Coor []float64 `json:"coor,omitempty" bson:"coor"`
	T    time.Time `json:"t,omitempty" bson:"t"` // fecha de creacion de la inserción
}

// estructura para decodificar la info extraida de la base de datos
type contactB struct {
	ID string `bson:"b"`
}

// GetContactIds obtenemos los ids de las personas que han estado en contacto conmigo
func (contact *Contact) GetContactIds(device string) ([]primitive.ObjectID, error) {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// traemos los parametros de control
	config := Config{}
	config.GetConfig()

	// dias atrás segun la configuración
	t := time.Now().UTC()
	td := t.AddDate(0, 0, config.Delta)

	// se ejecuta la consulta a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("contacts")
	cur, err := collection.Find(ctx, bson.M{"a": device, "t": bson.M{"$gte": td}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var ids []primitive.ObjectID
	for cur.Next(ctx) {
		// extraemos los ids con los cuales he tenido contacto
		var mContact contactB
		err := cur.Decode(&mContact)
		if err != nil {
			return nil, err
		}

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

// GetInfected dados los ids con los cuales he interactuado buscamos cuales han resultado positivos
func (contact *Contact) GetInfected(devices *[]primitive.ObjectID) (int64, error) {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": bson.M{"$in": (*devices)}, "data.infected": 1}
	// se ejecuta el conteo de usuarios infectados y que esten en mis ids de contacto
	collection := common.Client.Database(common.DATABASE).Collection("devices")
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return count, nil
}
