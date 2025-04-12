package connection

// User represents the authenticated user in a WebSocket connection
type User struct {
	ID       int64
	Username string
	Email    string
	// Add other user fields as needed
}

// NewUser creates a new User instance
func NewUser(id int64, username, email string) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
	}
}
