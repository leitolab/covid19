package models

import (
	"context"
	"errors"
	"ieliot/src/common"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// DecreasingRisk ...
func DecreasingRisk(risk float64) error {

	t := time.Now().UTC()
	td := t.AddDate(0, 0, -3) // 72 horas de vida del virus

	filter := bson.M{"data.risk": bson.M{"$gte": (-1 * risk)}, "t": bson.M{"$gte": td}}
	update := bson.M{"$inc": bson.M{"data.risk": risk}}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("places")
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("The document to update was not found")
	}

	return nil
}
