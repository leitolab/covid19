package controllers

import (
	"fmt"
	"ieliot/src/common"
	"ieliot/src/models"
	"time"

	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Route2Handler Servicio rest para la entidad People
func Route2Handler(c *fasthttp.RequestCtx) {
	if c.IsGet() {
		getRoute2Handler(c)
	} else {
		common.MethodNotAllowed(c)
	}
}

// Función de obtención de lugares a los que tuvo exposición
func getRoute2Handler(c *fasthttp.RequestCtx) {
	// validamos que el token del dispositivo sea válido y obtenemos la información contenida
	var err error
	origin := &models.Device{}
	if err := origin.ValidateToken(string(c.Request.Header.Peek("authorization")), 1); origin == nil || err != nil {
		common.Forbidden(c)
		return
	}

	//extraemos los argumentos de los params
	loc, _ := time.LoadLocation("America/Bogota")
	now := time.Now().Add(time.Duration(-300) * time.Minute) // hace 5 horas

	t0 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	t1 := t0.AddDate(0, 0, 1)

	// buscamos los lugares con los que tuvo contacto en esta ventana de tiempo
	var places []bson.M
	if places, err = origin.FindContactByDay(t0, t1); err != nil {
		common.SendJSON(c, &bson.M{"err": err.Error()})
		return
	}

	var place models.Place
	// iteramos para cada place y contamos la cantidad de usuarios infectados que han tenido contacto
	for _, placeBson := range places {
		// si hay contactos contamos cuantos estan con marcador positivo
		placeBson["cases"] = place.GetInfected(placeBson["_id"].(string))
	}

	var total0 time.Time
	var total1 time.Time
	lenPlaces := len(places)
	for i := 0; i < lenPlaces-1; i++ {
		timerange0 := places[i]["timeRange"].(bson.M)
		timerange1 := places[i+1]["timeRange"].(bson.M)

		timerange0["end"] = timerange1["start"]
		t0 := timerange0["start"].(primitive.DateTime).Time()
		t1 := timerange0["end"].(primitive.DateTime).Time()

		if i == 0 {
			total0 = t0
		}
		if i == lenPlaces-2 {
			total1 = t1
		}

		timerange0["start"] = fmt.Sprintf("%d:%d", t0.Hour(), t0.Minute())
		timerange0["end"] = fmt.Sprintf("%d:%d", t1.Hour(), t1.Minute())
		places[i]["timeRange"] = timerange0

		places[i]["duration"] = fmt.Sprintf("%v", t1.Sub(t0))
	}

	if lenPlaces > 1 {
		timerange0 := places[lenPlaces-1]["timeRange"].(bson.M)
		t0 = timerange0["start"].(primitive.DateTime).Time()
		t1 = timerange0["end"].(primitive.DateTime).Time()
		total1 = t1

		timerange0["start"] = fmt.Sprintf("%d:%d", t0.Hour(), t0.Minute())
		timerange0["end"] = fmt.Sprintf("%d:%d", t1.Hour(), t1.Minute())
		places[lenPlaces-1]["duration"] = "0h0m0s"
		places[lenPlaces-1]["timeRange"] = timerange0
	}

	for i := 0; i < lenPlaces; i++ {
		delete(places[i], "_id")
	}
	duration := fmt.Sprintf("%v", total1.Sub(total0))

	common.SendJSON(c, &bson.M{"places": places, "duration": duration})
}
