package smarthouse

type Controlable interface {
	Encender() error
	Apagar() error
	EstadoActual() string
}
