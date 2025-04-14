package parte2

import (
	"fmt"
)
func main() {
	var i int
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
}
