package domain

// User es una estructura que representa a un usuario
type User struct {
	ID        string
	Following []string
	Tweets    []Tweet
}
