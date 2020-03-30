package main

import (
	"context"
	"flag"
	"net"
	"os"
	"time"
	"unsafe"

	"ieliot/src/common"
	"ieliot/src/controllers"

	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Funcion principal del servidor por aquí se inicia todo el proceso
func main() {
	// se generan las variables globales a partir de variables del entorno
	common.Configure()

	// puerto estandar para desplegar el servicio
	bindHost := flag.String("bind", "0.0.0.0:8080", "set bind host")
	flag.Parse()

	var err error
	s := &fasthttp.Server{
		Handler: mainHandler,
		Name:    "ieliot",
	}

	// Contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// conexion al servidor de mongo
	common.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		log.Fatal(err)
	}

	// creacion del listener a las variables bind host
	var ln net.Listener
	ln = common.GetListener(*bindHost)

	// levantamiento del servidor
	if err = s.Serve(ln); err != nil {
		log.Fatalf("Error when serving incoming connections: %s", err)
	}
}

// Router del sistema dado que no se manejara un gran grupo de rutas esta solución es suficiente y eficiente
func mainHandler(c *fasthttp.RequestCtx) {
	path := c.Path()
	switch *(*string)(unsafe.Pointer(&path)) {

	case "/rest/v1/emit/":
		controllers.IotHandler(c)

	case "/rest/v1/contact/":
		controllers.ContactHandler(c)

	case "/rest/v1/people/":
		controllers.PeopleHandler(c)

	case "/rest/v1/route/":
		controllers.RouteHandler(c)

	case "/rest/v1/place/":
		controllers.PlaceHandler(c)

	case "/rest/v1/device/":
		controllers.DeviceHandler(c)

	case "/rest/v1/login/":
		controllers.LoginHandler(c)

	case "/rest/v1/client/":
		controllers.ClientHandler(c)

	default:
		controllers.Default(c)
	}
}
