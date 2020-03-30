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

// Función encargada de manejar la emision de datos de los usuarios
func postIotHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	var err error
	origin := &models.Device{}
	if err = origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// estructura para parsear la entrada, se espera un json válido
	var data interface{}
	// paseamos la informacion del dispositivo, se espera un json valido. En caso de que no se consiga se responde con un bad request
	if err = json.Unmarshal(c.PostBody(), &data); err != nil {
		common.BadRequest(c)
		return
	}

	// adicionamos información de origen de la data en base al token proporcionado
	iot := models.Iot{}
	iot.Device = origin.ID         // id del dispositivo
	iot.Client = origin.Client     // id del cliente del cual pertenece el dispositivo
	iot.Created = time.Now().UTC() // hora de registro en el sistema
	iot.Data = data

	// actualizamos el mapa de tics para la ubicación
	if err = iot.Upsert(origin.Product); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	// buscamos las personas en rango
	var iots []models.Iot
	if iots, err = iot.Near(origin.Product); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	if len(iots) > 0 {
		// intentamos guardar los contactos en caso de presentarse
		iot.InsertContact(&iots)
	}

	// buscamos las localizaciones en rango
	var places []models.Place
	if places, err = iot.NearPlaces(origin.Product); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	if len(places) > 0 {
		// intentamos guardar los contactos en caso de presentarse
		iot.InsertContactPlaces(&places)
		iot.UpdateRiskPlaces(&places)
	}

	// entregamos el resultado de la transacción al usuario
	response := bson.M{}
	response["success"] = true
	common.SendJSON(c, &response)
}
