package controllers

import (
	"encoding/json"
	"ieliot/src/common"
	"ieliot/src/models"

	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// LoginHandler Manejador de los metodos de entrada de la peticion
func LoginHandler(c *fasthttp.RequestCtx) {
	if c.IsPost() {
		postLoginHandler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// Funcion de login recibe username y password y genera un JWT que es enviado como json
func postLoginHandler(c *fasthttp.RequestCtx) {
	// estructura para parsear la entrada, se espera un json válido
	login := models.Login{}
	if err := json.Unmarshal(c.PostBody(), &login); err != nil {
		common.BadRequest(c)
		return
	}
	if err := login.LoginDevice(); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	// enviamos el token de acceso al usuario
	common.SendJSON(c, &bson.M{"token": login.Token})
}
