# BLion Transaction

## Descripción

BLion transaction es un servicio perteneciente a un conjunto de servicios core de la arquitectura de BLion que tiene
como fin controlar y gestionar las transacciones que se generan en la blockchain, además de todos los métodos que se
requieren para las transacciones.

## Alcance

El presente servicio está diseñado para ser utilizado por los servicios internos de BLion definidos en la arquitectura.
Además de que dichos servicios deben de estar correctamente autenticados para poder usar el servicio.

Este servicio no tiene alcance externo por sí mismo y no está diseñado para ser consumido de manera directa por
aplicaciones de terceros.

## Casos de uso

* Se requiere para la creación de un bloque inicial conocido como Genesis con la transacción inicial
* Se requiere la creación de bloques temporales
* Se requiere que dependiendo del minado los bloques temporales pasen a formar parte de la cadena de bloques
* Se requiere poder obtener un listado de los bloques de la cadena de bloques para poder verlos
* Se requiere poder obtener un bloque de la cadena de bloques por medio de su hash o identificador
* Se requiere poder minar los bloques temporales
* Se requiere que la información del bloque esté encriptada de extremo a extremo con algoritmo AES-256

## Casos de uso no soportados

* Se requiere poder obtener los datos de un bloque
* Se requiere poder obtener bloques específicos de la cadena de bloques
* Se requiere poder actualizar los datos de un bloque de la cadena de bloques

## Arquitectura

![Arquitectura BLion](https://lh3.googleusercontent.com/fife/AAWUweWqgKXdZdpzQmQ5aKjeDtD9I6gNNmItwG8SqAs_94SIIo40f7ndPup6YL92THAMLBgW4riBpxwGhAJdy6ajILyOqE-bPO-mtArb-TC_g0esL-0dk3GWB5uZCc3cJWh3G8RbnKo2uNxpVpVrxOGvRBCDjPU33rEAfJdEA2DFWTTpI_45-UtYsKaQf40HlaqWRtrFf_7147PVtJqZtveY1w3JRmtE8kAqX751RWtMfRVYRgeox9meIzr6PghGJ2N4vIW8fU-azuHVlpQqe_OM6rBWfCPlFFd_254AbPlLeNL7r-IW1wZPxjcAf6gT9Pe0jcK9Vu2-4p2PSsBTKBBdxlSLJBkQqiVRqzEy8zEilIcMpNrL42LspBUPH4CrZhKP11BuNHdDUFQRWMgRBuVPsOdWfVF--IV8S6cKnDt7DIo_PvI5S5XHwD9bbFzq8Nu4eavCPeQHrWqcKjqPVW9311r5Bju-5USBM2MgGM-ebooxz0S-Aysrnt1v7_WkAdu1y3jMkDtofLTxeEaXC-__DW1K2xYVGhscdEA-jXM8aEI-29sk8diJY5IntNUFw6aHHeOPUAm3zRclKx6FQjgpRPjUyNzgaeiIzwZu1vo1383CZkCI772kWECMUeAE1R430ZvsMBfgfHTtx5LaO5P2fqFImVjKi5DHSeh-rJCqL7EcIjChu3MwZ305EcSy7gZhIdSWcH1NtgeofeaMhxCa5cZ_PzdKGthoh-GpqeDeSlKKNzXHouZ2SBwCQHNpP-QlUSJKdtZkjjc6eDmUL7DInGfWRiiRyELYX_8B486KJJtPtVv-pIX_lu6-Fo0RAeIYqOcAQwTVz8XnmOQpJzBwQNOnjCNQRMsCvLLVa7B2BDqa8jAdS38fWpDlGwaHTyZuIl13dVMMbohaY6ivYEOHUO_yvIffh5NTD3BBG3lBGImjwywuGNCjS2W95GJJret4zCBjKCfp-hzEeVkJACNg7Kbsijdw8saoRH4uQ6GPifB9zjSfGuP8frV-SH1wmLosLB7wW5JtUM9L61scgYS2zlhUbKoxgac5ngpvqnfl3DstbevwMpJBBsQsVlgB3jqc_gngN-ub9COrJZDcPFG9zUXfYURk7N2dQB74i-m4BzeCF7PUTuh__3ViANK0HiL3cUZbHMpnLJfDLXH9ZbJiexCHPzyopW9aKdUl61n6wZwiQVuDf3S8vWlOlKi5_EsjO0PoVV9rITNcvXSch9ivOIYr9l5dwloucdOjLElWphWnu8DqBE7yLnU7nkaPC5bfvhwyjRXYMt75j42I37r6gFisUCrmtobbOodeeEikbEQr=w977-h937)

### Diagramas

### Modelos de datos

El presente diagrama muestra el modelo de datos que se utiliza para el servicio.

![Modelo Relacional](https://lh3.googleusercontent.com/fife/AAWUweUzIfkWY7dyA5jX-enrKmZmKuJewwNMqyq5vdfQsSqCpBcUUV5O4YJUJVKcED6VmGWoRuRFW5lDjKWIdYZ5MWOTZ4J7mKSNiMIsNAvJTCzzXgyXqZ6HssrO8XNYjFUP493sNEQUf9WaOtwXkLxpBYRtSoauGDZcoFQla6b_wCxMP33ndI_UNwQ89BrzHqcFeSfWalse3ZJheSW8-ubsj5QYoeaNJKrDyk6JuJWGUFS4ako6pi-TayEWu7qsXCnxJ_iA9N512iUUVAfrwKiOSn_yNLMAwnG2gIBmYQg-KSqrH0cVPtiQd0rtn5Z5rYWICeV-L9qjeMh0qhimfTqtbOTxEod6PvoGc5K5jBlEmQoYKE1EaNnfUzYN4DQMDWT9iH4k74etM6Fd18jUE4s86gUPOH2iwefozjJmTKjUQ8OoWBQreFW0kTHBoIXEcx7GzbFtDp5o0E2dHQYQ77r-WXYB2mvuBKlBWD5WZ8a6jkUGhHGx_3cWyEQWCEm5XH5zXGoaGdWeMIb6dwIbjGR4tN_niuEOy2pXSzXmWof0nDBiQMQuGvu5D3pIVg9uKwusfRP20fTulVxDs8GSs9bGrmVPFMyqrta5xrFFfbhCmOpkfrASBXq49TeqJ75UjfWt4kPCFPU9MG3xklstTnsI69eLwpfwj8I7La5bfNDMVCVGIMjU8v0lgqm7W29ZpB7iXlT5Suatp8w9RD2lA6XsGvlecrwEiEgSeMbf4p-M1STIEBz-ynojP5Sy-13bWRSoKvZBCA1UN5SSEuUp_4We-GopatXD69qjfmx3Ws9YBHjgF9bDGYaeGjfUE-59BKzxh_Tg0oPIds9gQ2St8jQxTg9ctBEaWlW0CdKafmn2veFM34voN3tJCKdmGx0JHbsgopRDD96CnpHl44rkgWIqUnIVG7HkW3OXzapJr5vkAK2yDP1w01JcAG5GJy8B5ay04OvJcspLlG_TexLMWa5qZjPXRNRzt93yXuE3RpRvubMFvZWamvIS2RJy55k3OszZsJPXnYzoilX5AkFSRzalDoC3sK6LDr0yecZOL3x_7aaObp3ozcEZ0mYtCwjnjluoc02sXICu3FlFIlDPlk9itzhop6gT9jwM2DmS0nHxiXRjR9ZZPDerRTGa_9fN3IqBMVict5NqmjUm-l7OSZ5Mo-JFgXk81_tiJAM3TMrmRJ_TTgVDTcmpmm1tqExXcsUf7SHTpcCscrr6dR1uo2wWbJ-Wq1560mXXpDpmz-F4ivz8vxaoP8b7N60cLQB_cvjhWNfO1622OLVgC3CI59OU9mURxhBCLdzGkouk-UmCJRol=w1920-h954)

## Limitaciones

Este servicio presenta las siguientes limitaciones con respecto a sus funciones y/o características:

* No se puede ejecutar en el sistema operativo MacOS

## Costos

La creación de las transacciones tiene un costo que depende de la información que se escribe en ella misma, tener en
cuenta que se cobrara un costo adicional al monto transferido
que son la comisión por escribir en la blockchain.
Dichos costos o aproximados se pueden obtener de la siguiente en el simulador de la página
de [BLion](https://www.bjungle.net).

## Instalación y Ejecución

Este servicio cuenta con un archivo ejecutable para los diferentes sistemas operativos,
que se encuentran junto a este archivo. Para generar dichos ejecutables se deben ejecutar los siguientes comandos:

#### Comandos de creación y compilación

##### Linux

````bash
    GOOS=linux  GOARCH=amd64 go build
````

##### Windows

````bash
    GOOS=windows  GOARCH=amd64 go build
````

#### Ejecutables

    - `blockchain-transaction.exe - windows`
    - `blockchain-transaction - linux`

##### Maneras y entornos para ejecutar el servicio:

Si está en desarrollo, se puede ejecutar el siguiente comando:

````bash
go run main.go
````

De estar en un entorno productivo debe ejecutar el archivo del servicio dependiendo del sistema operativo en donde
se encuentre instalado.

##### Windows

En el caso de windows dar clic derecho sobre el archivo ejecutable y seleccionar "Run as administrator"

##### Linux

En el caso de linux ingresar al directorio del ejecutable por medio de la terminal y ejecutar los siguientes comandos:

````bash
chmod 777 transaction
````

````bash
./blockchain-transaction &&
````

Tener en cuenta que para que el servicio se pueda ejecutar correctamente usted tiene que estar en modo superusuario (
ROOT).

#### Instalación y ejecución de GRPC

Para instalar ``GRPC`` y poder usarlo se tiene que primero
instalar [Protoc](https://github.com/protocolbuffers/protobuf/releases) para el sistema operativo adecuado además
se debe configurar la variable de entorno `protoc` y que apunta al directorio de instalación a la
carpeta `bin` (`c:\protoc\bin`), luego se necesita instalar los plugins para golang en las versiones más actuales:

````bash
 go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
````

````bash
 go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
````

## <span style="color:yellow">Atención</span>

Para versiones anterioeres a protoc v1.27.1 se ejecuta este comando

````bash
protoc -I api/grpc/proto --go_out=internal/grpc api/grpc/proto/*.proto
````

Para versiones superiores a protoc v1.27.1 se ejecuta este comando

````bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=false:. --go-grpc_opt=paths=source_relative api/grpc/proto/*.proto
````
