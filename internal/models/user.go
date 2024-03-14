package models

// User - структура, описывающая сущность пользователя.
type User struct {
	Id       int    `json:"-" db:"id"`              // Id - id пользователя.
	Name     string `json:"name" db:"name"`         // Name - никнейм пользователя.
	Email    string `json:"email" db:"email"`       // Email - адрес электронной почты пользователя.
	Password string `json:"password" db:"password"` // Password - зашифрованный пароль пользователя.
}
