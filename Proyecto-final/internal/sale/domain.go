package sale

import "time"

// Sale represents a system sale with metadata for auditing and versioning.
type Sale struct {
	ID        string
	user_id   string
	estado    string
	amount    float32
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

// UpdateFields represents the optional fields for updating a Sale.
// A nil pointer means “no change” for that field.
type UpdateFields struct {
	estado string
}
