package controllers

import (
	"encoding/json"
	"ieliot/src/common"
	"ieliot/src/models"
	"strconv"

	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// TaskHandler Servicio rest para la entidad People
func TaskHandler(c *fasthttp.RequestCtx) {
	if c.IsPost() {
		postTaskHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

func postTaskHandler(c *fasthttp.RequestCtx) {

	// estructura para parsear la entrada, se espera un json válido
	var body map[string]string
	if err := json.Unmarshal(c.PostBody(), &body); err != nil {
		common.BadRequest(c)
		return
	}

	origin := &models.Device{}
	if err := origin.ValidateToken(body["token"], 9); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	var err error
	switch body["task"] {
	case "decreasing_risk":
		err = decreasingRisk(body)
	}

	if err != nil {
		common.SendJSON(c, &bson.M{"error": err.Error()})
	}
	common.SendJSON(c, &bson.M{"success": true})
}

func decreasingRisk(body map[string]string) error {
	var f float64
	var err error
	if f, err = strconv.ParseFloat(body["risk"], 64); err != nil {
		return err
	}

	if err = models.DecreasingRisk(f); err != nil {
		return err
	}

	return nil
}
