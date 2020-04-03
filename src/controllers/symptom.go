package controllers

import (
	"encoding/json"
	"fmt"
	"ieliot/src/common"
	"ieliot/src/models"
	"time"

	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// SymptomHandler Servicio rest para la entidad Device
func SymptomHandler(c *fasthttp.RequestCtx) {
	if c.IsPost() {
		postSymptomHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

func postSymptomHandler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	// estructura para parsear la entrada, se espera un json válido
	symptom := models.Symptom{}
	if err := json.Unmarshal(c.PostBody(), &symptom); err != nil {
		common.BadRequest(c)
		return
	}
	t := time.Now().UTC()
	symptom.Created = &t

	fmt.Println(symptom)

	if err := symptom.UpdateOne(origin.ID); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	common.SendJSON(c, &bson.M{"success": true})
}
