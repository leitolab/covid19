package models

import (
	"context"
	"errors"
	"ieliot/src/common"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// Login estructura del para realizar login y generacion de tokens
type Login struct {
	Username string  `json:"username,omitempty" bson:"_id,omitempty"` // _id de mongo para login
	Password string  `json:"password,omitempty"`                      // Se debe omitir siempre en los json
	Token    *string `json:"token,omitempty"`                         // JWT generado
}

// LoginDevice ...
func (login *Login) LoginDevice() error {
	if login.Password == "" {
		return errors.New("params are required: password")
	}

	// contexto timeout para la solicitud a mongo
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Lo buscamos dentro de mongo, en caso de no existir se retorna un error estándar por seguridad
	var device Device
	collection := common.Client.Database(common.DATABASE).Collection("devices")
	if err := collection.FindOne(ctx, bson.M{"data.phone": login.Username}).Decode(&device); err != nil {
		return err
	}

	if device.Status != 1 {
		return errors.New("device isn't ready")
	}
	// generamos el token de acceso
	login.Token = device.MakeToken()

	// validamos que la contraseña concuerde, en caso de que no retornamos un error estándar por seguridad
	if ok := checkPasswordHash(login.Password, device.Password); !ok {
		return errors.New("params are required: password")
	}

	return nil
}

// funcion para generar un bcrypt hash a partir de un password plano
func hashPassword(password string) (string, error) {
	// se emplea u factor de trabajo de 10 para balance seguridad rendimiento
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

// funcion que compara un password plano con su hash para verificar su coincidencia
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
