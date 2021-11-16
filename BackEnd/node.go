package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Info struct {
	Tipo     string
	AddrNodo string
	Valor    string
}

var localhostReg string //localhost:9001
var localhostNot string //localhost:9002
var actualConfiguration int
var remotehost string

var bitacoraAddr []string //todos los localhost + puerot de notificaicon
/*
func main() {
	bufferIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el puerto de registro: ")
	puerto, _ := bufferIn.ReadString('\n')
	puerto = strings.TrimSpace(puerto)
	localhostReg = fmt.Sprintf("localhost:%s", puerto)

	fmt.Print("Ingrese el puerto de notificacion: ")
	puerto, _ = bufferIn.ReadString('\n')
	puerto = strings.TrimSpace(puerto)
	localhostNot = fmt.Sprintf("localhost:%s", puerto)

	go activarServicioRegistro() //rol de servidor

	//rol de cliente
	fmt.Print("Ingrese del puerto del nodo a solicitar registro: ")
	puerto, _ = bufferIn.ReadString('\n')
	puerto = strings.TrimSpace(puerto)
	remotehost = fmt.Sprintf("localhost:%s", puerto) //solicito el punto de conexión para la red

	procesarConfiguracionActual()

	//si no es le primer nodo que crea la red
	if puerto != "" {
		registrarSolicitud(remotehost) //envio solicitud de registro a la red
		validarConfiguration()
	}
	//rol de servidor
	procesarNotificaciones()
}
*/
func procesarConfiguracionActual() {
	bufferIn := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese configuracion del nodo: \n")
	fmt.Print("1 Create Person: \n")
	fmt.Print("2 Save Person: \n")
	fmt.Print("3 Create Prediction: \n")
	strConf, _ := bufferIn.ReadString('\n')
	strConf = strings.TrimSpace(strConf)
	numConf, _ := strconv.Atoi(strConf)
	actualConfiguration = numConf
	fmt.Printf("Actual Configuration %s \n", strConf)
}

func validarConfiguration() {
	info := Info{"NODEFUNCTION", localhostNot, strconv.Itoa(actualConfiguration)}
	comunicarTodos(info)
}

func activarServicioRegistro() {
	//colocar en modo escucha por el puerto de registro
	ln, _ := net.Listen("tcp", localhostReg)
	defer ln.Close()
	//atencion a conexiones
	for {
		con, _ := ln.Accept()
		go manejadorRegistro(con) //atender varias conexiones de forma concurrente
	}
}

func manejadorRegistro(con net.Conn) {
	//1 recibir identificacion del nuevo nodo
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	ident, _ := bufferIn.ReadString('\n')
	ident = strings.TrimSpace(ident)
	//2 enviar la bitacora incluido su identificaion al nuevo nodo
	//codigicar la bitacora
	bitacoraBytes, _ := json.Marshal(append(bitacoraAddr, localhostNot))
	fmt.Fprintln(con, string(bitacoraBytes))

	//3 comunicar al resto de nodos que llegó uno nuevo
	info := Info{"REGISTRATION", localhostNot, ident}
	comunicarTodos(info)

	//4 actualizar la bitacora del nodo actual
	bitacoraAddr = append(bitacoraAddr, ident)

	//mostrar la bitacora
	fmt.Println(bitacoraAddr)
}

func comunicarTodos(info Info) {
	//recuperando la bitacora del nodo actual y recorrerla para enviar la notificación a cada uno
	for _, addr := range bitacoraAddr {
		notificar(addr, info)
	}
}

func notificar(addr string, info Info) {
	con, _ := net.Dial("tcp", addr)
	defer con.Close()

	byteInfo, _ := json.Marshal(info)
	fmt.Fprintln(con, string(byteInfo))

}

func registrarSolicitud(remotehost string) {
	//llamar al nodo remoto
	con, _ := net.Dial("tcp", remotehost)
	defer con.Close()

	fmt.Fprintln(con, localhostNot)

	//procesar lo que el nodo remoto responde
	bufferIn := bufio.NewReader(con)
	bitacoraNodo, _ := bufferIn.ReadString('\n')

	//guardar localmente
	var bitacoraTemp []string
	json.Unmarshal([]byte(bitacoraNodo), &bitacoraTemp)

	//asigna a la bitacora local
	bitacoraAddr = bitacoraTemp

	fmt.Println(bitacoraAddr)

}

func procesarNotificaciones() {
	ln, _ := net.Listen("tcp", localhostNot)
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejadorNotificacionesEnviadas(con)
	}
}

func manejadorNotificacionesEnviadas(con net.Conn) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	bInfo, _ := bufferIn.ReadString('\n')

	var info Info
	json.Unmarshal([]byte(bInfo), &info)

	//ident = strings.TrimSpace(ident)
	//actualizar la bitacora
	if info.Tipo == "REGISTRATION" {
		bitacoraAddr = append(bitacoraAddr, info.Valor)
		//imprimir la bitacora
		fmt.Println(bitacoraAddr)
	}
	if info.Tipo == "NODEFUNCTION" {
		numConf, _ := strconv.Atoi(info.Valor)

		if actualConfiguration == numConf {
			fmt.Println("Configuraciones Iguales")
			procesarConfiguracionActual()
			validarConfiguration()
		} else {
			fmt.Println("Configuraciones Diferentes")
		}
	}

	if info.Tipo == "GETBITACORA" {
		con, _ := net.Dial("tcp", info.AddrNodo)
		defer con.Close()
		bitacoraTemp := bitacoraAddr
		bitacoraTemp = append(bitacoraTemp, localhostNot)
		fmt.Fprintln(con, bitacoraTemp)
	}

	if info.Tipo == "GETNODECONFIGURATION" {
		con, _ := net.Dial("tcp", info.AddrNodo)
		defer con.Close()
		fmt.Fprintln(con, actualConfiguration)
	}

}
