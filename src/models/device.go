package models

import (
	"context"
	"errors"
	"ieliot/src/common"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// Device estructura del dispositivo a nivel de base de datos pero tambien a nivel de JWT
// Estado de infección en data.infected 0: no infectado, 1: infectado, 2 recuperado, 3: muerto
type Device struct {
	ID          string     `json:"_id,omitempty" bson:"_id,omitempty"`  // _id de mongo
	Data        bson.M     `json:"data,required"`                       // data del device
	Phone       string     `json:"phone,bson:"phone,"`                  // telefono del registro
	Password    string     `json:"password,omitempty"`                  // Se debe omitir siempre en los json
	Client      string     `json:"client,0mitempty"`                    // propietario del dispositvo
	Product     string     `json:"product,omitempty"`                   // producto al cual pertenece el dispositivo, ej: gs
	Expires     int64      `json:"expires,omitempty" bson:",omitempty"` // Solo se una a nivel de JWT
	Scope       int8       `json:"scope,omitempty" `                    // alcance de privilegios del dispositivo
	Status      int8       `json:"status,omitempty"`                    // bandera para saber si el dispositivo esta activo en la plataforma
	SymptomLast []int64    `json:"symptom" bson:"symptom_last"`         // ultimos sintomas reportados por el usuario
	Created     *time.Time `json:"created,omitempty"`                   // fecha de creacion del usuario
	Updated     *time.Time `json:"updated,omitempty"`                   // fecha de actualizacion del usuario
}

// FindOne Obtenemos de la base de datos un device dado su id
func (device *Device) FindOne() error {
	// Parseamos el id a una estructura de id de mongo
	deviceID, err := primitive.ObjectIDFromHex(device.ID)
	if err != nil {
		return err
	}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("devices")
	if err := collection.FindOne(ctx, bson.M{"_id": deviceID, "product": device.Product}).Decode(device); err != nil {
		return err
	}

	device.Password = ""
	return nil
}

// FindByClient Obtenemos de la base de datos todos los devices de un cliente
func (device *Device) FindByClient() (*[]Device, error) {

	var devices []Device

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("devices")
	cur, err := collection.Find(ctx, bson.M{"product": device.Product, "client": device.Client})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		// una variable para cada device encontrado
		var mDevice Device
		err := cur.Decode(&mDevice)
		if err != nil {
			return nil, err
		}
		mDevice.Password = ""
		// agregamos el device al array
		devices = append(devices, mDevice)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return &devices, nil
}

// cambia el password en el objeto device
func (device *Device) changePassword() error {
	if len(device.Password) < 8 {
		return errors.New("The minimum password length is 8 characters")
	}
	var err error
	if device.Password, err = hashPassword(device.Password); err != nil {
		return err
	}
	return nil
}

// InsertOne Insertamos en la base de datos un device
func (device *Device) InsertOne() error {

	if device.Password == "" {
		return errors.New("params are required: password")
	}
	if device.Client == "" {
		return errors.New("params are required: client")
	}
	if device.Scope == 0 {
		return errors.New("params are required: scope")
	}
	if device.Status == 0 {
		return errors.New("params are required: status")
	}
	if device.Data == nil {
		return errors.New("params are required: data")
	}
	device.Phone = device.Data["phone"].(string)
	if device.Phone == "" {
		return errors.New("params are required: phone")
	}

	if err := device.changePassword(); err != nil {
		return err
	}

	device.ID = ""
	device.Status = 1
	now := time.Now().UTC()
	device.Created = &now
	device.Updated = &now

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("devices")
	res, err := collection.InsertOne(ctx, device)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		device.ID = oid.Hex()
	}
	device.Password = ""
	return nil
}

// UpdateOne Actualizamos en la base de datos un device
func (device *Device) UpdateOne() error {
	if device.ID == "" {
		return errors.New("params are required: _id")
	}
	if device.Data == nil {
		return errors.New("params are required: data")
	}

	// Parseamos el id a una estructura de id de mongo
	deviceID, err := primitive.ObjectIDFromHex(device.ID)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	device.Updated = &now

	filter := bson.M{"_id": deviceID, "product": device.Product}
	update := bson.M{"$set": bson.M{
		"data.infected":           device.Data["infected"],
		"data.infected_timestamp": time.Now().UTC(),
	}}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("devices")
	res, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("The document to update was not found")
	}

	device.Password = ""
	return nil
}

// DeleteOne Eliminamos de la base de datos un elemento
func (device *Device) DeleteOne() error {
	// Parseamos el id a una estructura de id de mongo
	deviceID, err := primitive.ObjectIDFromHex(device.ID)
	if err != nil {
		return err
	}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("devices")
	res, err := collection.DeleteOne(ctx, bson.M{"_id": deviceID, "product": device.Product})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("The document to delete was not found")
	}

	return nil
}

// FindContactByTime encontramos los lugares con los cuales tuvo contacto en una ventana de tiempo
func (device *Device) FindContactByTime(t0, t1 time.Time) ([]interface{}, error) {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := []bson.M{
		bson.M{"$match": bson.M{"t": bson.M{"$gte": t0, "$lte": t1}, "a": device.ID}},
		bson.M{"$project": bson.M{"_id": 0, "place": 1}},
		bson.M{"$sort": bson.M{"t": 1}},
		bson.M{"$limit": 100},
	}

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("contacts_places")
	cur, err := collection.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// obtenemos los ids de contacto
	var places []interface{}
	for cur.Next(ctx) {
		var mPlaceInt bson.M
		err := cur.Decode(&mPlaceInt)
		if err != nil {
			return nil, err
		}

		mPlace := mPlaceInt["place"]
		places = append(places, mPlace)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return places, nil
}

func mRiskLevel(risk float64) string {
	if risk < 33 {
		return "low"
	} else if risk < 66 {
		return "meddium"
	}
	return "high"
}

// FindContactByDay encontramos los lugares con los cuales tuvo contacto en un dia
func (device *Device) FindContactByDay(t0, t1 time.Time) ([]bson.M, error) {
	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := []bson.M{
		bson.M{"$match": bson.M{"t": bson.M{"$gte": t0, "$lte": t1}, "a": device.ID}},
		bson.M{"$project": bson.M{"_id": 0, "place": 1, "t": 1}},
		bson.M{"$sort": bson.M{"t": 1}},
		bson.M{"$limit": 300},
	}

	// se ejecuta la insercion a la base de datos
	collection := common.Client.Database(common.DATABASE).Collection("contacts_places")
	cur, err := collection.Aggregate(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	// obtenemos los ids de contacto
	var places []bson.M
	for cur.Next(ctx) {
		var mPlaceInt bson.M
		err := cur.Decode(&mPlaceInt)
		if err != nil {
			return nil, err
		}
		var mPlace Place
		bsonBytes, _ := bson.Marshal(mPlaceInt["place"])
		bson.Unmarshal(bsonBytes, &mPlace)
		front := bson.M{
			"_id":         mPlace.ID,
			"place":       mPlace.Data["title"],
			"description": mPlace.Data["description"],
			"riskLevel":   mRiskLevel(mPlace.Data["risk"].(float64)),
			"cases":       0,
			"timeRange": bson.M{
				"start": mPlaceInt["t"],
				"end":   mPlaceInt["t"],
			},
		}
		places = append(places, front)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return places, nil
}

// ====================================================================================================
//
// 	Funciones con fines utilitarios para la estructura device
//
// ====================================================================================================

// MakeToken funcion generadora de JWT para un dispositivo especifico
func (device *Device) MakeToken() *string {
	// hora actual para el conteo del token en UTC
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"dev": device.ID,                                // id del dipositivo
		"cli": device.Client,                            // id del cliente del cual pertenece el dispositivo
		"pro": device.Product,                           // producto que usa el core, ej: gs
		"scp": device.Scope,                             // alcance de privilegios del dispositivo
		"nbf": now.Unix(),                               // no valido antes de
		"exp": now.Add(365 * 2 * 24 * time.Hour).Unix(), // fecha de expiracion en 2 años
	}

	// generacion del token
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, &claims)
	// firma del token
	tokenSigned, err := token.SignedString(common.SIGNKEY)
	if err != nil {
		panic(err.Error())
	}

	// retorno del token
	return &tokenSigned
}

// ValidateToken validams el token de autenticacion
func (device *Device) ValidateToken(tokenString string, scope int8) error {
	// parseamos el token y validamos con la clave publica
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return common.VERIFYKEY, nil
	})
	if err != nil || !token.Valid {
		return err
	}
	//e xtraemos los datos del token
	claims := token.Claims.(jwt.MapClaims)

	// parseamos las variables de los claims a la estructura Device, si hay algun problema el dispositivo es nil
	var ok bool
	var tmp float64
	if tmp, ok = claims["scp"].(float64); !ok {
		return err
	}
	device.Scope = int8(tmp)
	if device.Scope < scope {
		return errors.New("invalid scope")
	}

	if device.ID, ok = claims["dev"].(string); !ok {
		return err
	}
	if device.Client, ok = claims["cli"].(string); !ok {
		return err
	}
	if device.Product, ok = claims["pro"].(string); !ok {
		return err
	}
	if tmp, ok = claims["exp"].(float64); !ok {
		return err
	}
	device.Expires = int64(tmp)

	// si todo esta ok se regresa el dispositivo
	return nil
}
