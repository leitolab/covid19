package controllers

import (
	"ieliot/src/common"

	"github.com/valyala/fasthttp"
)

const ieliot = `  _ ______ _      _____ ____ _______ 
(_)  ____| |    |_   _/ __ \__   __|
 _| |__  | |      | || |  | | | |   
| |  __| | |      | || |  | | | |   
| | |____| |____ _| || |__| | | |   
|_|______|______|_____\____/  |_|   
									
									`

// Default todo lo que no coincida con una ruta llega aquí
func Default(c *fasthttp.RequestCtx) {
	common.SendTEXT(c, ieliot)
}
