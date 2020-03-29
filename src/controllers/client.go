package controllers

import (
	"encoding/json"
	"ieliot/src/common"
	"ieliot/src/models"

	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// ClientHandler Servicio rest para la entidad Client
func ClientHandler(c *fasthttp.RequestCtx) {
	if c.IsGet() {
		getClientHandler(c)
	} else if c.IsPost() {
		postClientHandler(c)
	} else if c.IsPut() {
		putClientHandler(c)
	} else if c.IsDelete() {
		deleteClientHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// Función de obtención de clients por id y por producto de forma intrínseca
func getClientHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 9); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// obtenemos los query params
	args := c.QueryArgs()
	// creamos un client con el producto del token de origen
	client := models.Client{Product: origin.Product}

	// obtenemos el cliente si el query param _id esta presente
	if id := args.Peek("_id"); id != nil {
		client.ID = string(id)
		if err := client.FindOne(); err != nil {
			common.SendJSON(c, &bson.M{"err": err.Error()})
			return
		}
		common.SendJSON(c, &bson.M{"client": &client})
		return
	}

	// se obtiene la lista de los clientes por producto
	clients, err := client.Find()
	if err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"client": &clients})
}

// Función ingesta de un client
func postClientHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 9); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}
	// estructura para parsear la entrada, se espera un json válido
	client := models.Client{}
	if err := json.Unmarshal(c.PostBody(), &client); err != nil {
		common.BadRequest(c)
		return
	}

	// creamos un client con el producto del token de origen
	client.Product = origin.Product
	if err := client.InsertOne(); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"client": client})
	return
}

// Función de actualizacion de un client
func putClientHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 9); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}
	// estructura para parsear la entrada, se espera un json válido
	client := models.Client{}
	if err := json.Unmarshal(c.PostBody(), &client); err != nil {
		common.BadRequest(c)
		return
	}

	// actualizamos un client con el producto del token de origen
	client.Product = origin.Product
	if err := client.UpdateOne(); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"client": &client})
	return
}

// Función eliminacion de un client
func deleteClientHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 9); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// obtenemos los query params
	args := c.QueryArgs()
	// creamos un client con el producto del token de origen
	client := models.Client{Product: origin.Product}

	// solamemte eliminamos cuando el parámetro _id esta presente
	if id := args.Peek("_id"); id != nil {
		client.ID = string(id)
		if err := client.DeleteOne(); err != nil {
			common.SendJSON(c, &bson.M{"err": err.Error()})
			return
		}
		common.SendJSON(c, &bson.M{"_id": client.ID})
		return
	}

	common.BadRequest(c)
}
