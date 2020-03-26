package controllers

import (
	"encoding/json"
	"ieliot/src/common"
	"ieliot/src/models"
	"time"

	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// IotHandler Manejador de los metodos de entrada de la peticion
func IotHandler(c *fasthttp.RequestCtx) {
	if c.IsPost() {
		postIotHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// Funcion de alimento de datos de los dispositivos iot
func postIotHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido
	var err error
	origin := &models.Device{}
	if err = origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// estructura para parsear la entrada, se espera un json válido
	var data interface{}
	// paseamos la informacion del dispositivo, se espera un json valido en caso de que no se consiga se responder con un bad request
	if err = json.Unmarshal(c.PostBody(), &data); err != nil {
		common.BadRequest(c)
		return
	}

	iot := models.Iot{}
	// adicionamos información de origen de la data en base al token proporcionado
	iot.Device = origin.ID         // id del dispositivo
	iot.Client = origin.Client     // id del cliente del cual pertenece el dispositivo
	iot.Created = time.Now().UTC() // hora de registro en el sistema
	iot.Data = data

	if err = iot.Upsert(origin.Product); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	var ids []string
	if ids, err = iot.Near(origin.Product); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}
	if len(ids) > 0 {
		if err := iot.Contact("contacts", &ids); err != nil {
			common.SendJSON(c, &bson.M{"err": err.Error()})
			return
		}
	}

	// entregamos al usuario el id resultado de la inserción, esto puede servir para hacer una traza de inserciones
	response := bson.M{}
	response["success"] = true
	// enviamos la respuesta al usuario
	common.SendJSON(c, &response)
}
