package models

// Route - структура, опичывающая сущность маршрута спидрана.
type Route struct {
	Id        int    `json:"id" db:"id"`                 // Id - id маршрута.
	Start     string `json:"start" db:"start"`           // Start - ссылка на стартовую статью маршрута.
	Finish    string `json:"finish" db:"finish"`         // Finish - ссылка на финишную статью маршрута.
	CreatorId int    `json:"creator_id" db:"creator_id"` // CreatorId - id пользователя, создавшего маршрут.
}
