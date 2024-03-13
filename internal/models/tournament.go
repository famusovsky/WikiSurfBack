package models

import "time"

// Tournament - структура, представляющая сущность соревнования.
type Tournament struct {
	Id        int       `json:"id" db:"id"`                 // Id - id соревнования.
	StartTime time.Time `json:"start_time" db:"start_time"` // StartTime - время начала соревнования.
	EndTime   time.Time `json:"end_time" db:"end_time"`     // EndTime - время конца соревнования.
	Pswd      string    `json:"pswd" db:"pswd"`             // Pswd - Код-пароль соревнования.
	Private   bool      `json:"private" db:"private"`       // Private - флаг, указывающий на закрытость соревнования.
	// Ratings string // XXX необходимы ли?
}

// TURelation - структура, представляющая отношение между соревнованием и пользователем.
type TURelation struct {
	TournamentId int `json:"tour_id" db:"tour_id"` // TournamentId - id соревнования, участвующего в отношении.
	UserId       int `json:"user_id" db:"user_id"` // UserId - id пользователя, участвующего в отношении.
}

// TURelation - структура, представляющая отношение между соревнованием и маршрутом.
type TRRelation struct {
	TournamentId int `json:"tour_id" db:"tour_id"`   // TournamentId - id соревнования, участвующего в отношении.
	RouteId      int `json:"route_id" db:"route_id"` // RouteId - id маршрута, участвующего в отношении.
}
