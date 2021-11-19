# ConcurrenteTF

Trabajo final del curso de Concurrente
Para correr el trabajo final es necesario tener instalado NodeJs y Golang
En primer lugar, ingresar a la carpeta Frontend y utilizar el comando `npm init` para descargar los paquetes necesarios para correr el frontend
Luego de descargar todos los paquetes debe utilizar el comando `npm start` para iniciar el proyecto frontend

En el caso del backend, es necesario tener instalado Golang. En primer lugar, realizar un pull de la imagen dockerizada del api con el comando
`docker pull raulinoxx/apigolang` y luego es necesario correr la imagen en un contenedor utilizando el siguiente comando `docker run --rm -p 9000:9000 -p 9009:9009/tcp --network bridge   -i -d  --name apigolang raulinoxx/apigolang`

Finalmente, es requerido inicializar los 3 nodos a formar parte del cluster. 
