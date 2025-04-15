package smarthouse

import (
	"errors"
	"fmt"

	"github.com/fatih/color"
)

type Dispositivo struct {
	nombre string
	estado bool
}

func (d *Dispositivo) Encender() error {
	if d.estado {
		return errors.New("El dispositivo ya está prendido")
	}
	d.estado = true
	return nil
}

func (d *Dispositivo) Apagar() error {
	if !d.estado {
		return errors.New("El dispositivo ya está apagado")
	}
	d.estado = false
	return nil
}

func (d *Dispositivo) EstadoActual() string {
	if d.estado {
		color.Green("Encendido")
		return "Encendido"
	}
	color.Red("Apagado")
	return "Apagado"
}

func (d *Dispositivo) PrintObject() string {
	return d.nombre
}
func NuevoDispositivo(nombre string, estado bool) (Dispositivo, error) {
	if len(nombre) > 10 {
		return Dispositivo{}, fmt.Errorf("el nombre no puede tener más de 10 caracteres")
	}
	return Dispositivo{nombre: nombre, estado: estado}, nil
}
