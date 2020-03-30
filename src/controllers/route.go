package controllers

import (
	"ieliot/src/common"
	"ieliot/src/models"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// RouteHandler Servicio rest para la entidad People
func RouteHandler(c *fasthttp.RequestCtx) {
	if c.IsGet() {
		getRouteHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// convertimos el parámetro a time
func timeFrombytes(bytes []byte) time.Time {
	i, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		return time.Now().UTC()
	}
	return time.Unix(i, 0)
}

// Función de obtención de personas dada un área máxima
func getRouteHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	var err error
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	//extraemos los argumentos de los params
	args := c.QueryArgs()
	var arg []byte

	// extracción de la latitud
	if arg = args.Peek("t0"); arg == nil {
		common.SendJSON(c, &bson.M{"err": "params are required: t0"})
		return
	}
	t0 := timeFrombytes(arg)

	// extracción de la longitud
	if arg = args.Peek("t1"); arg == nil {
		common.SendJSON(c, &bson.M{"err": "params are required: t1"})
		return
	}
	t1 := timeFrombytes(arg)

	// buscamos los lugares con los que tuvo contacto en esta ventana de tiempo
	var places []interface{}
	if places, err = origin.FindContactByTime(t0, t1); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"places": places})
}
