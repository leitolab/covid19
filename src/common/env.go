package common

import (
	"crypto/rsa"
	"io/ioutil"
	"os"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// DATABASE Base de datos en mongo que usará el sistema
	DATABASE string
	// Client Puntero de conexión a mongo, mongo administra el pool, se debe implementar con mogos
	Client *mongo.Client
	// VERIFYKEY Clave publica de verificación del token
	VERIFYKEY *rsa.PublicKey
	// SIGNKEY Clave de firmado de los tokens
	SIGNKEY *rsa.PrivateKey
)

// Configure tareas iniciales del sistema
func Configure() {
	// se leen las variables de entorno del archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// configuracion del nivel de logs para el sistema
	if "DEBUG" == os.Getenv("ENVIRONMENT") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	// configuracion del nombre de la base de datos a la cual se conectará
	DATABASE = os.Getenv("DATABASE")
	// carga de las claves de firmas de JWT, la clave privada es supremanente sensible
	SIGNKEY = loadRSAPrivateKeyFromDisk("keys/rs512-4096-private.pem")
	VERIFYKEY = loadRSAPublicKeyFromDisk("keys/rs512-4096-public.pem")
}

// Función encargada de cargar del disco la clave privada y parsearla de PEM a RSA key
func loadRSAPrivateKeyFromDisk(location string) *rsa.PrivateKey {
	keyData, err := ioutil.ReadFile(location)
	if err != nil {
		panic(err.Error())
	}
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if err != nil {
		panic(err.Error())
	}
	return signKey
}

// Función encargada de cargar del disco la clave publica y parsearla de PEM a RSA key
func loadRSAPublicKeyFromDisk(location string) *rsa.PublicKey {
	keyData, err := ioutil.ReadFile(location)
	if err != nil {
		panic(err.Error())
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		panic(err.Error())
	}
	return verifyKey
}
