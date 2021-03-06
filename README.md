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

## Desplegar el bot

Puedes ejecutarlo localmente o usar ansible + fabric, que tendrás que
instalar.

Empieza por usar Ansible para preparar el entorno. Define las
variables de entorno que uses para el bot, los logs y el resto de los
APIS, y ejecuta:

	ansible-playbook -i hosts go.yml

Esto instala los fuentes y algunas utilidades necesarias, aunque
supone que Go y Python están instalados.

A continuación tendrás que usar fabric para desplegar la última
versión y ejecutarlo remotamente. Sustituye

	env.hosts = [ '159.100.248.62' ]
	env.user = "root"
	env.release_path= "CuantoQuedaBot"

(que corresponde a un servidor cedido gentilmente por
[Exoscale](http://exoscale.ch) en `fabfile.py` por los valores para tu
programa en particular. Establece los valores de las variables de
entorno como el de arriba y otros para que vayan los logs y haz

	fab build
	fab start

Para echarlo a andar. Cuando te canses,

	fab stop

detendrá el proceso remoto

## Utilizar el bot

Existe un bot que se ejecuta en un servidor de Cloud9 (o de donde
pille) y utiliza el código de este repositorio, actualizándose automáticamente. 
Para utilizar el bot tal y como está programado en este repositorio basta con iniciar una conversación en telegram con él a través de este enlace:

[https://telegram.me/CuantoQuedaBot](https://telegram.me/CuantoQuedaBot)
