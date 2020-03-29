package controllers

import (
	"ieliot/src/common"
	"ieliot/src/models"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// ContactHandler Servicio rest para la entidad Contact
func ContactHandler(c *fasthttp.RequestCtx) {
	if c.IsGet() {
		getContactHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// Función de obtención de clients por id y por producto de forma intrínseca
func getContactHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	var err error
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// obtenemos los ids de los usuarios con los cuales he tenido contacto en los últimos X dias
	var ids []primitive.ObjectID
	contact := models.Contact{}
	if ids, err = contact.GetContactIds(origin.ID); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	// si no hay contactos regresamos
	if len(ids) == 0 {
		common.SendJSON(c, &bson.M{"count": 0})
		return
	}

	// si hay contactos contamos cuantos estan con marcador positivo
	var count int64
	if count, err = contact.GetInfected(&ids); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"count": count})
}
