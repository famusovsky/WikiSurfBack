package models

import "time"

// Sprint - структура, представляющая сущность спидран-спринта.
type Sprint struct {
	Id           int       `json:"id" db:"id"`                   // Id - id спринта.
	UserId       int       `json:"user_id" db:"user_id"`         // UserId - id пользователя, проведшего спринт.
	RouteId      int       `json:"route_id" db:"route_id"`       // RouteId - id маршрута, к которому относится спринт.
	TournamentId int       `json:"tour_id" db:"tour_id"`         // TournamentId - id соревнования, к которому относится спринт.
	Path         string    `json:"path" db:"path"`               // Path - пройденный в ходе спринта путь в формате JSON.
	Success      bool      `json:"success" db:"path"`            // Success - успешность спринта.
	LengthTime   int64     `json:"length_time" db:"length_time"` // LengthTime - длительность спринта в ms.
	StartTime    time.Time `json:"start_time" db:"start_time"`   // StartTime - время старта спринта.
}
