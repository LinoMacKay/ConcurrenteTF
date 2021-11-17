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
	"strings"
	"sync"

	"github.com/gorilla/mux"
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

//ips preseteados
var apiIp = "localhost:9000"
var ips = []string{"localhost:9002"}

//bitacoras
var wg sync.WaitGroup
var totalBitacora = make(chan []string)
var bitacoraConfg []string

var personas []Persona

func main() {
	cargarAlumnos()
	handleRequests()

}

func cargarAlumnos() {
	personas = []Persona{
		{"Pedro", []Sintoma{
			{"Gripe", 1}}},
	}
}

func listarAlumnos(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("Content-Type", "application/json")
	jsonbytes, _ := json.Marshal(personas)
	io.WriteString(resp, string(jsonbytes))

}

func buscarPersonas(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	//Estas maneras son para obtener un request param
	//code := req.FormValue("codigo")
	//code := req.URL.Query()["codigo"][0]

	//Estas maneras son para obtener un path variable
	id := strings.TrimPrefix(req.URL.Path, "/buscar_alumnos/")

	for _, alumno := range personas {
		if alumno.Nombre == id {
			jsonbytes, _ := json.Marshal(alumno)
			io.WriteString(resp, string(jsonbytes))
		}
	}

	log.Println(id)
}

func predict(resp http.ResponseWriter, req *http.Request) {
	//validacion
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.Header().Set("Access-Control-Allow-Origin", "*")
	//llamada por post
	if req.Method == "POST" {
		if req.Header.Get("Content-type") == "application/json" {
			log.Println("Accede a agregar alumno")
			jsonBytes, err := ioutil.ReadAll(req.Body)

			if err != nil {
				http.Error(resp, "Error al recuperar el body", http.StatusInternalServerError)
			} else {
				var oPersona Persona
				json.Unmarshal(jsonBytes, &oPersona)

				fmt.Println(oPersona)
				sendPatienteToNode(oPersona)

				personas = append(personas, oPersona)
				resp.Header().Set("Content-Type", "application/json")
				io.WriteString(resp, `
					{
						"respuesta":"Registro satisfactorio"
					}
				`)
			}
		} else {
			http.Error(resp, "Contenido inválido", http.StatusBadRequest)
		}

	} else {
		http.Error(resp, "Metodo invalido", http.StatusMethodNotAllowed)
	}
}
func reciveData() {
	ln, _ := net.Listen("tcp", "localhost:9009")
	defer ln.Close()
	for {
		con, _ := ln.Accept()
		go manejarRespuetas(con)
	}
}

func sendPatienteToNode(oPersona Persona) {

	go reciveData()
	ip := ips[rand.Intn(len(ips))]
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
	fmt.Println("PEDIR CONFIGURACIONES")

	bitacoras := <-totalBitacora
	dialForConfig(bitacoras)
}

func dialForConfig(bitacoras []string) {

	for i := 0; i < 2; i++ {
		func() {
			con, _ := net.Dial("tcp", bitacoras[i])
			defer con.Close()

			toSend := &Info{"GETNODECONFIGURATION", bitacoras[i], ""}
			byteInfo, _ := json.Marshal(toSend)
			fmt.Println(toSend)
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

	fmt.Println(info)

	if info.Tipo == "SENDAEA" {
		fmt.Println("Llegó hasta pedir configuracion", info.Valor)
	}
	if info.Tipo == "SENDBITACORA" {
		bitacora := strings.Split(info.Valor, ",")
		totalBitacora <- bitacora
		fmt.Println("Llegó hasta pedir bitacora", totalBitacora)
	}
}

func mostrarInicio(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, "Inicio")
}

func enableCORS(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	}).Methods(http.MethodOptions)
	router.Use(middlewareCors)
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
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
	http.HandleFunc("/listarPersonas", listarAlumnos)
	//http.HandleFunc("/", mostrarInicio)
	r.HandleFunc("/predict", predict)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":9000", nil))

}
