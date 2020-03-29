package controllers

import (
	"encoding/json"
	"ieliot/src/common"
	"ieliot/src/models"

	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// PlaceHandler Servicio rest para la entidad Place
func PlaceHandler(c *fasthttp.RequestCtx) {
	if c.IsGet() {
		getPlaceHandler(c)
	} else if c.IsPost() {
		postPlaceHandler(c)
	} else if c.IsDelete() {
		deletePlaceHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// Función de obtención de marcadores dado un área máxima
func getPlaceHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	var err error
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	//extraemos los argumentos de los params
	args := c.QueryArgs()

	var place models.Place
	var arg []byte
	place.Coor = []float64{0, 0}

	// extracción de la latitud
	if arg = args.Peek("lat"); arg == nil {
		common.SendJSON(c, &bson.M{"err": "params are required: lat"})
		return
	}
	place.Coor[0] = float64frombytes(arg)

	// extracción de la longitud
	if arg = args.Peek("lon"); arg == nil {
		common.SendJSON(c, &bson.M{"err": "params are required: lon"})
		return
	}
	place.Coor[1] = float64frombytes(arg)

	// extracción de la precisión
	if arg = args.Peek("accuracy"); arg == nil {
		common.SendJSON(c, &bson.M{"err": "params are required: accuracy"})
		return
	}
	place.Accuracy = float64frombytes(arg)

	// buscamos los lugares en rango
	var places []models.Place
	if places, err = place.FindNear(); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"places": places})
}

// Función de almacenamiento de marcadores
func postPlaceHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	var err error
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// estructura para parsear la entrada, se espera un json válido
	place := models.Place{Device: origin.ID}
	// paseamos la informacion del dispositivo, se espera un json valido. En caso de que no se consiga se responde con un bad request
	if err = json.Unmarshal(c.PostBody(), &place); err != nil {
		common.BadRequest(c)
		return
	}

	// Insertamos el marcador
	if err = place.InsertOne(); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"_id": place.ID})
}

// Función de eliminación de un device
func deletePlaceHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	args := c.QueryArgs()

	// solamente eliminamos un device si el parámetro _id esta presente
	if id := args.Peek("_id"); id != nil {
		place := models.Place{ID: string(id), Device: origin.ID}
		if err := place.DeleteOne(); err != nil {
			common.SendJSON(c, &bson.M{"err": err.Error()})
			return
		}
		common.SendJSON(c, &bson.M{"success": true})
		return
	}

	common.BadRequest(c)
}
