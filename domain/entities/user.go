package entities

import (
	"time"
)

// ✅ JSONTime - Custom time formatting for JSON responses
type JSONTime struct {
	time.Time
}

// ✅ Custom JSON formatting: "YYYY-MM-DD HH:MM:SS"
func (t JSONTime) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(`"` + t.Format("2006-01-02 15:04:05") + `"`), nil
}

// ✅ User Entity (Business Model)
type User struct {
	ID        int      `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	CreatedAt JSONTime `json:"created_at"`
}

// ✅ Business Logic Method (Example)
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
