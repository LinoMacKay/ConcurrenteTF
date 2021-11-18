package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	//"strconv"

	"github.com/LinoMacKay/ConcurrenteTF/tree/master/BackEnd/RF"
	//"math"
)

func main() {

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
	targets := make([]string, 0)

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

		targets = append(targets, target)

	}

	train_inputs := make([][]interface{}, 0)

	train_targets := make([]string, 0)

	test_inputs := make([][]interface{}, 0)
	test_targets := make([]string, 0)

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
	// fmt.Printf("Cantidad de inputs: %d", len(inputs))
	// fmt.Printf("Cantidad de targets: %d", len(targe

	// fmt.Printf("Train Inputs Len: %d\n", len(train_inputs))
	// fmt.Printf("Train Targets Len: %d\n", len(train_targets))
	// fmt.Printf("Test Inputs Len: %d\n", len(test_inputs))
	// fmt.Printf("Test Targets Len: %d\n", len(test_targets))

	//// EJEMPLO DE PRUEBA
	ejemplo := [15]string{"0", "1", "0", "0", "", "0", "1", "0", "0", "0", "0", "0", "0", "1", "0"}

	apattern := ejemplo[:len(ejemplo)-1]
	atarget := ejemplo[len(ejemplo)-1]

	ej := make([]interface{}, 0)

	for _, value := range apattern {
		ej = append(ej, value)
	}

	fmt.Println(ej)

	ainputs := make([][]interface{}, 0)
	atargets := make([]string, 0)

	ainputs = append(ainputs, ej)
	atargets = append(atargets, atarget)

	///////////////////

	//Construcción del objeto RF
	forest := RF.BuildForest(inputs, targets, 10, 1500, len(train_inputs[0])) //100 trees

	test_inputs = train_inputs
	test_targets = train_targets
	err_count := 0.0

	//El ejemplo del github
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

	////// Su webada de ellos
	test_inputs = ainputs
	test_targets = atargets
	// fmt.Println("input ", ainputs)
	// fmt.Println("targets ", atargets)
	// true_positive := 0.0
	// false_positive := 0.0
	// true_negative := 0.0
	// false_negative := 0.0
	fmt.Println(test_inputs[0])

	for i := 0; i < len(test_inputs); i++ {
		output := forest.Predicate(test_inputs[i])
		fmt.Println("Se predijo de output: ", output)

	}

	fmt.Println(time.Since(start))
}
