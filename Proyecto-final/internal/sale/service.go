package sale

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// Service provides high-level sale management operations on a LocalStorage backend.
type Service struct {
	// storage is the underlying persistence for Sale entities.
	storage *LocalStorage
}

// NewService creates a new Service.
func NewService(storage *LocalStorage) *Service {
	return &Service{
		storage: storage,
	}
}

// Create adds a brand-new sale to the system.
// It sets CreatedAt and UpdatedAt to the current time and initializes Version to 1.
// Returns ErrEmptyID if sale.ID is empty.
func (s *Service) Create(sale *Sale) error {
	sale.ID = uuid.NewString()
	// esto es asi porque sino los test de update fallan al crear aleatorio
	if sale.Estado == "" {
		estado := []string{"pending", "approved", "rejected"}
		sale.Estado = estado[rand.Intn(3)]
	}
	now := time.Now()
	sale.CreatedAt = now
	sale.UpdatedAt = now
	sale.Version = 1

	return s.storage.Set(sale)
}

// Get retrieves a sale by its ID.
// Returns ErrNotFound if no sale exists with the given ID.
func (s *Service) Get(id string) (*Sale, error) {
	return s.storage.Read(id)
}

// Update modifies an existing sale's data.
// It updates Name, Address, NickName, sets UpdatedAt to now and increments Version.
// Returns ErrNotFound if the sale does not exist, or ErrEmptyID if sale.ID is empty.
func (s *Service) Update(id string, sale *UpdateFields) (*Sale, error) {
	existing, err := s.storage.Read(id)
	// controlo existencia
	if err != nil {
		return nil, ErrSaleNotFound
	}
	// reviso que el estado anterior es valido
	if existing.Estado != "pending" {
		return nil, ErrInvalidStateChange //Solo permite cambiar si el estado anterior es == pending
	}
	// me fijo que el estado nuevo es de los dos validos
	if !(sale.Estado == "approved" || sale.Estado == "rejected") {
		return nil, ErrInvalidNewState
	}

	existing.Estado = sale.Estado
	existing.UpdatedAt = time.Now()
	existing.Version++

	if err := s.storage.Set(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes a sale from the system by its ID.
// Returns ErrNotFound if the sale does not exist.
func (s *Service) Delete(id string) error {
	return s.storage.Delete(id)
}

// Crear endpoint GET /sales con filtros por user_id y status.
func (s *Service) ListByUserAndStatus(userID, status string) ([]Sale, error) {
	var filtered []Sale
	for _, sale := range s.storage.GetAll() {
		if sale.UserID == userID && sale.Estado == status {
			filtered = append(filtered, sale)
		}
	}
	return filtered, nil
}
