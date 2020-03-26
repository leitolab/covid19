package models

import (
	"context"
	"errors"
	"ieliot/src/common"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

const clientsCollection = "clients"

// Client ...
type Client struct {
	ID      string     `json:"_id,omitempty" bson:"_id,omitempty"` // _id de mongo
	Data    bson.M     `json:"data,required"`                      // data del cliente
	Product string     `json:"product,required"`                   // producto al cual pertenece el cliente
	Created *time.Time `json:"created,omitempty"`                  // fecha de creacion del usuario
	Updated *time.Time `json:"updated,omitempty"`                  // fecha de actualizacion del usuario
}

// FindOne Obtenemos de la base de dato un producto dado su id
func (client *Client) FindOne() error {
	// Parseamos el id a una estructura de id de mongo
	clientID, err := primitive.ObjectIDFromHex(client.ID)
	if err != nil {
		return err
	}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("clients")
	if err := collection.FindOne(ctx, bson.M{"_id": clientID, "product": client.Product}).Decode(client); err != nil {
		return err
	}
	return nil
}

// Find Obtenemos de la base de datos todos los clientes de un producto
func (client *Client) Find() (*[]Client, error) {
	var clients []Client

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("clients")
	cur, err := collection.Find(ctx, bson.M{"product": client.Product})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// una variable para cliente encontrado
		var mClient Client
		err := cur.Decode(&mClient)
		if err != nil {
			return nil, err
		}
		// agregamos los clientes al array
		clients = append(clients, mClient)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return &clients, nil
}

// InsertOne Insertamos en la base de datos un cliente
func (client *Client) InsertOne() error {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if client.Data == nil {
		return errors.New("params are required: data")
	}

	now := time.Now().UTC()
	client.ID = ""
	client.Created = &now
	client.Updated = &now

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("clients")
	res, err := collection.InsertOne(ctx, client)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		client.ID = oid.Hex()
	}

	return nil
}

// UpdateOne Actualizamos en la base de datos un cliente
func (client *Client) UpdateOne() error {
	if client.ID == "" {
		return errors.New("params are required: _id")
	}
	if client.Data == nil {
		return errors.New("params are required: data")
	}

	// Parseamos el id a una estructura de id de mongo
	clientID, err := primitive.ObjectIDFromHex(client.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": clientID, "product": client.Product}
	update := bson.M{"$set": bson.M{"data": client.Data, "updated": time.Now().UTC()}}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("clients")
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("The document to update was not found")
	}

	return nil
}

// DeleteOne Eliminamos de la base de datos un elemento
func (client *Client) DeleteOne() error {
	// Parseamos el id a una estructura de id de mongo
	clientID, err := primitive.ObjectIDFromHex(client.ID)
	if err != nil {
		return err
	}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("clients")
	res, err := collection.DeleteOne(ctx, bson.M{"_id": clientID})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("The document to delete was not found")
	}

	return nil
}
