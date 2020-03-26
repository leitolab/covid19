package common

import (
	"encoding/json"
	"net"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"gopkg.in/mgo.v2/bson"
)

// GetListener se obtiene un listener tcp4 del sistema con el puerto solicitado
func GetListener(listenAddr string) net.Listener {
	ln, err := net.Listen("tcp4", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	return ln
}

// MethodNotAllowed respuesta JSON en caso de no tener implementado el metodo solicitado
func MethodNotAllowed(c *fasthttp.RequestCtx) {
	c.SetStatusCode(405)
	c.SetContentType("application/json")
	c.WriteString(`{"err":"method not allowed"}`)
}

// BadRequest respuesta JSON en caso de una peticion malformada
func BadRequest(c *fasthttp.RequestCtx) {
	c.SetStatusCode(400)
	c.SetContentType("application/json")
	c.WriteString(`{"err":"bad request"}`)
}

// Forbidden respuesta JSON en caso de no contar con las credenciales adecuadas
func Forbidden(c *fasthttp.RequestCtx) {
	c.SetStatusCode(403)
	c.SetContentType("application/json")
	c.WriteString(`{"err":"forbidden"}`)
}

// SendJSON respuesta JSON a partir de una estructura flexible
func SendJSON(c *fasthttp.RequestCtx, j *bson.M) {
	jb, err := json.Marshal(j)
	if err != nil {
		c.SetStatusCode(500)
		c.SetContentType("application/json")
		c.WriteString(`{"err":"response can't be parsed"}`)
		return
	}
	c.SetContentType("application/json")
	c.Write(jb)
}

// SendTEXT respuesta tipo texto plano con alguna cadena en particular
func SendTEXT(c *fasthttp.RequestCtx, s string) {
	c.SetContentType("text/plain")
	c.WriteString(s)
}
