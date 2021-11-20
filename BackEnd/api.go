package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sintoma struct {
	Sintoma    string `json:"sintoma"`
	IsSelected int    `json:"isSelected"`
}

type Info struct {
	Tipo     string
	AddrNodo string
	Valor    string
}

type Persona struct {
	Nombre   string    `json:"name"`
	Sintomas []Sintoma `json:"sintomas"`
}
type Pacients struct {
	CreatedAt  time.Time          `bson:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at"`
	ID         primitive.ObjectID `bson:"_id"`
	Persona    Persona            `bson:"person"`
	Prediction string             `bson:"prediction"`
}

type Result struct {
	Prediction string `json:"prediction"`
}

//ips preseteados
//var apiIp = "localhost:9000"
var ips = []string{"host.docker.internal:9002", "host.docker.internal:9004", "host.docker.internal:9006"}

//bitacoras
var wg sync.WaitGroup
var wg2 sync.WaitGroup

var result string
var totalBitacora = make(chan []string)
var totalConfig = make(chan string, 3)
var confings []string

//var bitacoraConfg []string

func main() {
	handleRequests()

}

func predict(resp http.ResponseWriter, req *http.Request) {
	//validacion
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	//llamada por post
	if req.Method == "POST" {
		if req.Header.Get("Content-type") == "application/json" {
			jsonBytes, err := ioutil.ReadAll(req.Body)

			if err != nil {
				http.Error(resp, "Error al recuperar el body", http.StatusInternalServerError)
			} else {
				//var oPersona Persona
				//json.Unmarshal(jsonBytes, &oPersona)
				wg2.Add(1)
				sendPatienteToNode(jsonBytes)
				//personas = append(personas, jsonBytes)
				wg2.Wait()
				resultJson := Result{result}
				fmt.Println(resultJson)
				resp.Header().Set("Content-Type", "application/json")
				resp.WriteHeader(http.StatusOK)
				json.NewEncoder(resp).Encode(resultJson)

			}
		} else {
			http.Error(resp, "Contenido inválido", http.StatusBadRequest)
		}

	} else {
		http.Error(resp, "Metodo invalido", http.StatusMethodNotAllowed)
	}
}
func reciveData() {
	ln, _ := net.Listen("tcp", "0.0.0.0:9009")
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejarRespuetas(con)
	}
}

func sendPatienteToNode(jsonBytes []byte) {
	if len(confings) == 0 {
		go reciveData()
	}
	ip := ips[rand.Intn(len(ips))]
	if len(confings) == 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			con, _ := net.Dial("tcp", ip)
			defer con.Close()

			toSend := &Info{"GETBITACORA", ip, ""}
			byteInfo, _ := json.Marshal(toSend)
			fmt.Fprintln(con, string(byteInfo))
		}()
		wg.Wait()
	}
	if len(totalBitacora) > 0 {
		ips = <-totalBitacora
	}
	if len(confings) == 0 {
		wg.Add(1)
		go dialForConfig(ips)
		wg.Wait()

		fmt.Println("IMPRIMIR VALORES")
		for i := 0; i < len(ips); i++ {
			confings = append(confings, <-totalConfig)
		}
	}
	fmt.Println("Configuraciones", confings)
	for i, v := range confings {
		if v == strconv.Itoa(1) {
			ipToSend := ips[i]
			go func() {
				con, _ := net.Dial("tcp", ipToSend)
				defer con.Close()
				myString := string(jsonBytes[:])
				toSend := &Info{"GETJSON", ipToSend, myString}
				byteInfo, _ := json.Marshal(toSend)
				fmt.Fprintln(con, string(byteInfo))
				fmt.Println("ENVIE LOS VALORES", toSend)
			}()
		}
	}
}

func dialForConfig(bitacoras []string) {
	defer wg.Done()

	for i := 0; i < len(bitacoras); i++ {
		func() {
			con, _ := net.Dial("tcp", bitacoras[i])
			defer con.Close()

			toSend := &Info{"GETNODECONFIGURATION", "localhost:9009", ""}
			byteInfo, _ := json.Marshal(toSend)
			fmt.Fprintln(con, string(byteInfo))
		}()
	}

}

func manejarRespuetas(con net.Conn) {
	defer con.Close()
	bufferIn := bufio.NewReader(con)
	bInfo, _ := bufferIn.ReadString('\n')
	var info Info
	json.Unmarshal([]byte(bInfo), &info)

	//fmt.Println(info)

	if info.Tipo == "SENDCONFIGURATION" {
		totalConfig <- info.Valor
	}
	if info.Tipo == "SENDBITACORA" {
		bitacora := strings.Split(info.Valor, ",")
		totalBitacora <- bitacora
	}
	if info.Tipo == "SENDRESULT" {

		//fmt.Println("RESULTADO", info.Valor)
		var pacient2 Pacients
		json.Unmarshal([]byte(info.Valor), &pacient2)
		result = pacient2.Prediction
		fmt.Println("Resultado" + result)
		wg2.Done()
	}
}

func mostrarInicio(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, "Inicio")
}

func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			// and call next handler!
			next.ServeHTTP(w, req)
		})
}

func handleRequests() {

	r := mux.NewRouter()
	enableCORS(r)
	r.HandleFunc("/", mostrarInicio)
	//http.HandleFunc("/", mostrarInicio)
	r.HandleFunc("/predict", predict)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":9000", nil))

}
