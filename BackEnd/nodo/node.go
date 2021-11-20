package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fxsjy/RF.go/RF/Regression"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var apiIP = "localhost:9009"

var collection *mongo.Collection
var ctx = context.TODO()

type Persona struct {
	Nombre   string    `json:"name"`
	Sintomas []Sintoma `json:"sintomas"`
}
type Sintoma struct {
	Sintoma    string `json:"sintoma"`
	IsSelected int    `json:"isSelected"`
}
type Pacients struct {
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	ID         primitive.ObjectID `bson:"_id"`
	Persona    Persona            `bson:"person"`
	Prediction string             `bson:"prediction"`
}
type Info struct {
	Tipo     string
	AddrNodo string
	Valor    string
}

var wg sync.WaitGroup

var confings []string
var localhostReg string //localhost:9001
var localhostNot string //localhost:9002
var actualConfiguration int
var remotehost string
var totalConfig = make(chan string, 3)

var bitacoraAddr []string //todos los localhost + puerot de notificaicon

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

	if info.Tipo == "GETBITACORA" && info.AddrNodo == localhostNot {
		fmt.Println("Me pidieron mi bitacora", info.AddrNodo)
		con, _ := net.Dial("tcp", apiIP)
		defer con.Close()

		bitacoatemp := bitacoraAddr
		bitacoatemp = append(bitacoatemp, localhostNot)
		justString := strings.Join(bitacoatemp, ",")

		toSend := &Info{"SENDBITACORA", localhostNot, justString}
		byteInfo, _ := json.Marshal(toSend)
		fmt.Fprintln(con, string(byteInfo))
	}

	if info.Tipo == "GETNODECONFIGURATION" {
		fmt.Println("Me pidieron mi configuración", info.AddrNodo)
		con, _ := net.Dial("tcp", info.AddrNodo)
		defer con.Close()

		fmt.Println("Envie mi configuracion", info.AddrNodo)

		toSend2 := &Info{"SENDCONFIGURATION", localhostNot, strconv.Itoa(actualConfiguration)}
		byteInfo2, _ := json.Marshal(toSend2)
		fmt.Fprintln(con, string(byteInfo2))
	}
	if info.Tipo == "GETJSON" && actualConfiguration == 1 {
		//fmt.Println("Me pasaron el json", info.Valor)

		var person Persona
		json.Unmarshal([]byte(info.Valor), &person)
		//fmt.Println(reflect.TypeOf(info.Valor))
		fmt.Println(person)
		searchDBNode(person)

	}
	if info.Tipo == "SENDCONFIGURATION" {
		totalConfig <- info.Valor
	}
	if info.Tipo == "GETPERSON" && actualConfiguration == 2 {
		//fmt.Println("Me pasaron el json", info.Valor)
		var person Persona
		json.Unmarshal([]byte(info.Valor), &person)
		addToDatabase(person)
	}
	if info.Tipo == "GETPACIENT" && actualConfiguration == 3 {
		var pacient Pacients
		json.Unmarshal([]byte(info.Valor), &pacient)
		fmt.Println(pacient)
		doMLProcess(pacient)
	}

}

func dialForConfig(bitacoras []string) {
	defer wg.Done()
	for i := 0; i < len(bitacoras); i++ {
		func() {
			con, _ := net.Dial("tcp", bitacoras[i])
			defer con.Close()

			toSend := &Info{"GETNODECONFIGURATION", localhostNot, ""}
			byteInfo, _ := json.Marshal(toSend)
			fmt.Fprintln(con, string(byteInfo))
		}()
	}

}

func searchDBNode(person Persona) {
	wg.Add(1)
	go dialForConfig(bitacoraAddr)
	wg.Wait()
	if len(confings) == 0 {
		fmt.Println("INICIAR GUARDADO DE CONFIGS")
		for i := 0; i < len(bitacoraAddr); i++ {
			confings = append(confings, <-totalConfig)
		}
	}
	fmt.Println("Configuraciones", confings)
	for i, v := range confings {
		if v == strconv.Itoa(2) {
			ipToSend := bitacoraAddr[i]
			go func() {
				con, _ := net.Dial("tcp", ipToSend)
				defer con.Close()
				personToSend, _ := json.Marshal(person)
				//Se envia mi persona
				toSend := &Info{"GETPERSON", "ip", string(personToSend)}
				byteInfo, _ := json.Marshal(toSend)
				fmt.Fprintln(con, string(byteInfo))
			}()
		}

	}
}

func addToDatabase(person Persona) {
	wg.Add(1)
	clientOptions := options.Client().ApplyURI("mongodb+srv://mongouser:raulino12@cluster0.qvc5e.mongodb.net/Cluster0?retryWrites=true&w=majority")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("Concurrente").Collection("Pacients")
	pacient := Pacients{time.Now(), time.Now(), primitive.NewObjectID(), person, ""}

	go func() {
		defer wg.Done()
		collection.InsertOne(ctx, pacient)
	}()

	fmt.Println("Se agregó correctamente", pacient)
	wg.Wait()
	searchMLNode(bitacoraAddr, pacient)
}

func searchMLNode(bitacoraAddr []string, pacient Pacients) {
	if len(confings) == 0 {
		wg.Add(1)
		go dialForConfig(bitacoraAddr)
		wg.Wait()
		fmt.Println("INICIAR GUARDADO DE CONFIGS")
		for i := 0; i < len(bitacoraAddr); i++ {
			confings = append(confings, <-totalConfig)
		}
	}
	fmt.Println("Configuraciones", confings)
	for i, v := range confings {
		if v == strconv.Itoa(3) {
			ipToSend := bitacoraAddr[i]
			go func() {
				con, _ := net.Dial("tcp", ipToSend)
				defer con.Close()
				pacientToSend, _ := json.Marshal(pacient)
				//Se envia mi persona
				toSend := &Info{"GETPACIENT", "ip", string(pacientToSend)}
				byteInfo, _ := json.Marshal(toSend)
				fmt.Fprintln(con, string(byteInfo))
			}()
		}
	}
}

func doMLProcess(pacient Pacients) {

	/*ML PROCESS*/

	var sintomas []float64

	for _, v := range pacient.Persona.Sintomas {
		var y float64 = float64(v.IsSelected)

		sintomas = append(sintomas, y)
	}
	sintomas = append(sintomas, 1)

	prediction := MLProcess(sintomas)

	clientOptions := options.Client().ApplyURI("mongodb+srv://mongouser:raulino12@cluster0.qvc5e.mongodb.net/Cluster0?retryWrites=true&w=majority")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("Concurrente").Collection("Pacients")
	update := pacient

	update.Prediction = fmt.Sprintf("%f%", prediction*100)
	update.UpdatedAt = time.Now()
	collection.FindOneAndReplace(ctx, bson.M{"_id": pacient.ID}, update)
	sendResult(update)
}

func MLProcess(sintomas []float64) float64 {
	start := time.Now()
	//Leer Dataset en el repositorio remoto
	resp, err := http.Get("https://raw.githubusercontent.com/LinoMacKay/ConcurrenteTF/master/BackEnd/dataset_covid.csv")
	if err != nil {
		print(err)
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll((resp.Body))
	s_content := string(content)
	lines := strings.Split(s_content, "\n")
	inputs := make([][]interface{}, 0)
	targets := make([]float64, 0)

	for _, line := range lines {
		line = strings.TrimRight(line, "\r\n")
		if len(line) == 0 {
			continue
		}
		tup := strings.Split(line, ",")
		pattern := tup[:len(tup)-1]
		target := tup[len(tup)-1]
		X := make([]interface{}, 0)
		for _, x := range pattern {
			X = append(X, x)
		}
		inputs = append(inputs, X)
		bitSize := 64
		floatNum, _ := strconv.ParseFloat(target, bitSize)
		targets = append(targets, floatNum)
	}

	train_inputs := make([][]interface{}, 0)
	train_targets := make([]float64, 0)
	test_inputs := make([][]interface{}, 0)
	test_targets := make([]float64, 0)
	for i, x := range inputs {
		if i%3 == 1 {
			test_inputs = append(test_inputs, x)
		} else {
			train_inputs = append(train_inputs, x)
		}
	}
	for i, y := range targets {
		if i%3 == 1 {
			test_targets = append(test_targets, y)
		} else {
			train_targets = append(train_targets, y)
		}
	}
	ejemplo := sintomas
	apattern := ejemplo[:len(ejemplo)-1]
	atarget := ejemplo[len(ejemplo)-1]
	ej := make([]interface{}, 0)
	for _, value := range apattern {
		ej = append(ej, value)
	}
	fmt.Println(ej)
	forest := Regression.BuildForest(train_inputs, train_targets, 50, len(inputs), 10)

	ainputs := make([][]interface{}, 0)
	atargets := make([]float64, 0)
	ainputs = append(ainputs, ej)
	atargets = append(atargets, atarget)
	test_inputs = train_inputs
	test_targets = train_targets
	err_count := 0.0
	fmt.Println(test_inputs[0])
	for i := 0; i < len(test_inputs); i++ {
		output := forest.Predicate(test_inputs[i])
		// fmt.Println(output)
		expect := test_targets[i]
		//fmt.Println(output,expect)
		if output != expect {
			err_count += 1
		}
		if i == 0 {
			fmt.Println("Se predijo de output:", output)
		}
	}
	fmt.Println("success rate:", 1.0-err_count/float64(len(test_inputs)))
	test_inputs = ainputs
	test_targets = atargets
	fmt.Println(test_inputs[0])
	var output float64
	for i := 0; i < len(test_inputs); i++ {
		output = forest.Predicate(test_inputs[i])
		fmt.Println("Se predijo de output: ", output)

	}
	fmt.Println(time.Since(start))
	return output
}

func sendResult(pacient Pacients) {
	func() {
		con, _ := net.Dial("tcp", apiIP)
		defer con.Close()
		pacientToSend, _ := json.Marshal(pacient)
		toSend := &Info{"SENDRESULT", "ip", string(pacientToSend)}
		byteInfo, _ := json.Marshal(toSend)
		fmt.Fprintln(con, string(byteInfo))
	}()
}
