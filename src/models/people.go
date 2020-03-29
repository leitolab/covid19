package models

import (
	"context"
	"ieliot/src/common"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// People persona cercana que esta emitiendo en el mapa
type People struct {
	Device   string    `json:"_id"  bson:"d"`                // id del cliente
	Accuracy float64   `json:"accuracy,omitempty"  bson:"-"` // precision de la peticion
	Coor     []float64 `json:"coor,omitempty"  bson:"coor"`  // coordenadas
	Data     bson.M    `json:"data"  bson:"data"`            // producto al cual pertenece el cliente
}

// FindNear obtenemos las personas cercanos
func (people *People) FindNear() ([]People, error) {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// verificamos que la emisión solo haya sido hasta hace un minuto
	t := time.Now()
	td := t.Add(time.Duration(-1) * time.Minute) // 1 minuto

	query := []bson.M{
		bson.M{
			"$geoNear": bson.M{
				"near": bson.M{
					"coordinates": people.Coor,
				},
				"distanceField": "data.calculated",
				"maxDistance":   people.Accuracy,
				"spherical":     true,
			}},
		bson.M{"$match": bson.M{"t": bson.M{"$gte": td}, "d": bson.M{"$ne": people.Device}}},
		bson.M{"$limit": 100},
	}

	// se ejecuta la consulta a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("covid19")
	cur, err := collection.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// obtenemos los datos de los usuarios cercanos
	var peoples []People
	for cur.Next(ctx) {
		var mPeople People
		err := cur.Decode(&mPeople)
		if err != nil {
			return nil, err
		}
		peoples = append(peoples, mPeople)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return peoples, nil
}
