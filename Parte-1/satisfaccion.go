package main

import (
	"fmt"
)

func main() {
	var i int
	var valores []int
	general := 0
	var contadores [5]int
	for {
		fmt.Print("Ingresa tu puntuación: 1-5 o cualquier otra tecla para finalizar la carga.")
		_, err := fmt.Scan(&i)
		if err != nil {
			break
		}
		if i >= 1 && i <= 5 {
			if i > 3 {
				general++
			} else if i < 3 {
				general--
			}
			contadores[i-1]++
			valores = append(valores, i)
		} else {
			break
		}
	}
	fmt.Println("Primeros 10 puntajes: ")
	if cap(valores) > 10 {
		fmt.Println(valores[:10])
	} else {
		fmt.Println(valores)
	}
	fmt.Printf("# Votos con 5: %d \n", contadores[4])
	fmt.Printf("# Votos con 4: %d \n", contadores[3])
	fmt.Printf("# Votos con 3: %d \n", contadores[2])
	fmt.Printf("# Votos con 2: %d \n", contadores[1])
	fmt.Printf("# Votos con 1: %d \n", contadores[0])
	fmt.Println("Resultado general: ")
	if general > 0 {
		fmt.Println("¡Buen resultado!")
	} else {
		fmt.Println("Resultado mejorable")
	}
}
