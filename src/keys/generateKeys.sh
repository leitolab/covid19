#! bin/bash
# script de apoyo para la generacion de claves publica y privada para el uso de JWT

# se genera una clave privada de 4096 bytes de longitud
openssl genrsa -out rs512-4096-private.pem 4096
# se genera una clave publica a partir de la clave privada
openssl rsa -in rs512-4096-private.pem -pubout > rs512-4096-public.pem