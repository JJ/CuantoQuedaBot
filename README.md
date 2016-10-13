# CuantoQuedaBot

Bot que te dice cuantos días quedan para entregar hitos y otras cosas
relacionadas con la asignatura
[Infraestructura Virtual](http://jj.github.io/IV)

Si tienes instalado Go, para que esto funcione dale a

	go get

(tendrás que definir antes GOPATH y GOBIN para que instale
correctamente la librería).


Consíguete despues
[tu propio API key para Telegram creando tu robot](http://bytelix.com/guias/crear-propio-bot-telegram/). Asígnaselo
a una variable de entorno con

	export BOT_TOKEN=my-bot-token-super-tocho


Y luego simplemente

	go run CuantoQuedaBot.go

Lo puedes ejecutar desde tu ordenador o cualquier otro que se quede
encendido todo el tiempo (o casi) como Nitrous o Cloud9.

Si quieres cambiar las contestaciones, no tienes más que editar el
fichero [`hitos.json`](hitos.json) y poner los títulos y URLs que
quieras. 


## Para probar como funciona 

Existe un bot que se ejecuta en un servidor de Cloud9 y utiliza el código de este repositorio, actualizándose automáticamente. 
Para ver qué es lo que hace basta con abrir una conversación de telegram con él a través de este enlace: 
[https://telegram.me/CuantoQuedaBot](https://telegram.me/CuantoQuedaBot)
