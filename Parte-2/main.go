package main

import (
	"fmt"

	"taller_go/parte_dos/smarthouse"

	"github.com/fatih/color"
)

func main() {
	var i int
	var list []smarthouse.Controlable
	for {
		fmt.Print("Menu.\n1.Crear un nuevo dispositivo indicando su nombre.\n2. Elegir uno para prenderlo o apagarlo. \n3. Visualizar todos los dispositivos indicando si están prendidos o apagados.\n\n4.Salir.\n ")
		fmt.Printf("Opción seleccionada: ")
		_, err := fmt.Scan(&i)
		if err != nil {
			color.Red("Error al leer la opcion")
		}
		switch i {
		case 1:
			list = agregarDispositivo(list)
		case 2:
			if len(list) > 0 {
				list = chooseDisp(list)
			}
		case 3:
			showStates(list)
		case 4:
			return
		}
	}
}

func agregarDispositivo(list []smarthouse.Controlable) []smarthouse.Controlable {
	var input string
	for {
		fmt.Println("Ingresa el nombre del dispositivo (solo 10 caracteres):")

		_, err := fmt.Scan(&input)
		if err != nil {
			color.Red("Error al leer el  nombre")
			continue
		}

		dispositivo, err := smarthouse.NuevoDispositivo(input, false)
		var aux smarthouse.Controlable
		aux = &dispositivo
		if err != nil {
			color.Red(err.Error())
			continue
		}
		list = append(list, aux)
		break
	}
	return list
}

func chooseDisp(list []smarthouse.Controlable) []smarthouse.Controlable {
	var i int
	var input int

	for {
		fmt.Println("Seleccioná tu dispositivo")
		for i = range list {
			p, ok := list[i].(smarthouse.Printeable)
			if ok {
				fmt.Printf("%d. %s. \n", i+1, p.PrintObject())
			}

		}
		i = i + 2
		fmt.Printf("%d. Volver. \n", i)
		fmt.Printf("Opción seleccionada: ")

		_, err := fmt.Scan(&input)
		if err != nil {
			continue
		}
		if input == i {
			break
		}
		if input >= 1 && input < i {
			var aux smarthouse.Controlable
			aux = list[input-1]
			for {
				fmt.Println("Que deseas hacer: ")
				fmt.Println("1. Encender. ")
				fmt.Println("2. Apagar. ")
				fmt.Println("3. Ver estado. ")
				fmt.Println("4. Volver. ")

				fmt.Printf("Opción seleccionada: ")

				_, err := fmt.Scan(&input)
				if err != nil {
					continue
				}
				if input == 4 {
					break
				}
				switch input {
				case 1:
					err = aux.Encender()
					if err != nil {
						color.Red(err.Error())
					}
				case 2:
					err = aux.Apagar()
					if err != nil {
						color.Red(err.Error())
					}
				case 3:
					aux.EstadoActual()
				}

			}
		}

	}
	return list
}
func showStates(list []smarthouse.Controlable) {
	for i := range list {

		p, ok := list[i].(smarthouse.Printeable)
		if ok {
			fmt.Println(p.PrintObject())
		}
		list[i].EstadoActual()
	}
}
