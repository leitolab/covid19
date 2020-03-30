package controllers

import (
	"ieliot/src/common"
	"ieliot/src/models"
	"strconv"

	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// PeopleHandler Servicio rest para la entidad People
func PeopleHandler(c *fasthttp.RequestCtx) {
	if c.IsGet() {
		getPeopleHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// convertimos el parámetro a flotante
func float64Frombytes(bytes []byte) float64 {
	f, err := strconv.ParseFloat(string(bytes), 64)
	if err != nil {
		return 0.0
	}
	return f
}

// Función de obtención de personas dada un área máxima
func getPeopleHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	var err error
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// generamos la variable persona
	people := models.People{Device: origin.ID, Coor: []float64{0, 0}}

	//extraemos los argumentos de los params
	args := c.QueryArgs()
	var arg []byte

	// extracción de la latitud
	if arg = args.Peek("lat"); arg == nil {
		common.SendJSON(c, &bson.M{"err": "params are required: lat"})
		return
	}
	people.Coor[0] = float64Frombytes(arg)

	// extracción de la longitud
	if arg = args.Peek("lon"); arg == nil {
		common.SendJSON(c, &bson.M{"err": "params are required: lon"})
		return
	}
	people.Coor[1] = float64Frombytes(arg)

	// extracción de la precisión
	if arg = args.Peek("accuracy"); arg == nil {
		common.SendJSON(c, &bson.M{"err": "params are required: accuracy"})
		return
	}
	people.Accuracy = float64Frombytes(arg)

	// buscamos las personas en rango
	var peoples []models.People
	if peoples, err = people.FindNear(); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"people": peoples})
}
