package parte2

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

type Controlable interface {
	Encender()
	Apagar()
	EstadoActual()
}
type Dispositivo struct {
	Nombre string
	Estado bool
}

func (d *Dispositivo) Encender() error {
	if d.Estado {
		return errors.New("El dispositivo ya está prendido")
	}
	d.Estado = true
	return nil
}
func (d *Dispositivo) Apagar() error {
	if !d.Estado {
		return errors.New("El dispositivo ya está apagado")
	}
	d.Estado = false
	return nil
}

func (d Dispositivo) EstadoActual() string {
	if d.Estado {
		color.Green("Encendido")
		return "Encendido"
	}
	color.Red("Apagado")
	return "Apagado"
}

func NuevoDispositivo(nombre string, estado bool) (Dispositivo, error) {
	if len(nombre) > 10 {
		return Dispositivo{}, fmt.Errorf("el nombre no puede tener más de 10 caracteres")
	}
	return Dispositivo{Nombre: nombre, Estado: estado}, nil
}
