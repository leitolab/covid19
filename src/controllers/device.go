package controllers

import (
	"encoding/json"
	"ieliot/src/common"
	"ieliot/src/models"

	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// DeviceHandler Servicio rest para la entidad Device
func DeviceHandler(c *fasthttp.RequestCtx) {
	if c.IsGet() {
		getDeviceHandler(c)
	} else if c.IsPost() {
		postDeviceHandler(c)
	} else if c.IsPut() {
		putDeviceHandler(c)
	} else if c.IsDelete() {
		deleteDeviceHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// Función obtencion de un device
func getDeviceHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// intentamos obtener un device por _id
	if err := origin.FindOne(); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}
	common.SendJSON(c, &bson.M{"device": &origin})
	return

	common.BadRequest(c)
}

// Función de creación de un device
func postDeviceHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 0); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// estructura para parsear la entrada, se espera un json válido
	device := models.Device{}
	if err := json.Unmarshal(c.PostBody(), &device); err != nil {
		common.BadRequest(c)
		return
	}

	// creamos un device con el producto del token de origen
	device.Product = "covid19"
	device.Client = "5e7b4796f5bd74c4162edf1e"
	device.Scope = 1
	device.Status = 1
	if err := device.InsertOne(); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"client": device})
}

// Función de actualizacion de device
func putDeviceHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// estructura para parsear la entrada, se espera un json válido
	device := models.Device{}
	if err := json.Unmarshal(c.PostBody(), &device); err != nil {
		common.BadRequest(c)
		return
	}

	// actualizamos un device con el producto del token de origen
	device.ID = origin.ID
	device.Product = "covid19"
	device.Client = "5e7b4796f5bd74c4162edf1e"
	device.Scope = 1
	device.Status = 1
	if err := device.UpdateOne(); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"device": device})
}

// Función de eliminación de un device
func deleteDeviceHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 9); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// obtenemos los query params
	args := c.QueryArgs()
	// creamos un client con el producto del token de origen
	device := models.Device{Product: origin.Product}

	// solamente eliminamos un device si el parámetro _id esta presente
	if id := args.Peek("_id"); id != nil {
		device.ID = string(id)
		if err := device.DeleteOne(); err != nil {
			common.SendJSON(c, &bson.M{"err": err.Error()})
			return
		}
		common.SendJSON(c, &bson.M{"_id": device.ID})
		return
	}

	common.BadRequest(c)
}
