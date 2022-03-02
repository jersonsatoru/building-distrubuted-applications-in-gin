package models

// swagger:model user
type User struct {
	// User's password
	//
	// required: true
	Username string `json:"username"`

	// User's login
	//
	// required: true
	Password string `json:"password"`
}
