package controllers

import (
	"encoding/json"
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
	} else if c.IsPut() {
		putContactHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// Función de obtención de clients por id y por producto de forma intrínseca
func getContactHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido
	var err error
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	var ids []primitive.ObjectID
	contact := models.Contact{}
	if ids, err = contact.GetContactIds(origin.ID); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	if len(ids) == 0 {
		common.SendJSON(c, &bson.M{"count": 0})
		return
	}

	var count int64
	if count, err = contact.GetInfected(&ids); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}
	common.SendJSON(c, &bson.M{"count": count})
}

// Función de actualizacion de un client
func putContactHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
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
}
