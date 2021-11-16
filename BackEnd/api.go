package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
)

type Sintoma struct {
	Sintoma    string `json:"sintoma"`
	isSelected int    `json:"isSelected"`
}

type Persona struct {
	Nombre   string    `json:"name"`
	Sintomas []Sintoma `json:"sintomas"`
}

var ips = []string{"localhost:9002", "localhost:9004", "localhost:9006"}

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

func agregarAlumno(resp http.ResponseWriter, req *http.Request) {
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
func mostrarInicio(resp http.ResponseWriter, req *http.Request) {
	io.WriteString(resp, "Inicio")
}

func handleRequests() {    
	
	r := mux.NewRouter()
	r.HandleFunc("/", mostrarInicio)
	http.HandleFunc("/listarPersonas", listarAlumnos)
	//http.HandleFunc("/", mostrarInicio)
	http.HandleFunc("/predict", agregarAlumno)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":9000", nil))


}
