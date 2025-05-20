package sale

import "time"

type Sale struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"` // antes: user_id
	Estado    string    `json:"estado"`  // antes: estado
	Amount    float32   `json:"amount"`  // antes: amount
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
}

type UpdateFields struct {
	Estado string `json:"estado"` // antes: estado
}
