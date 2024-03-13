package models

// RouteRating - структура, представляющая блок рейтинга маршрута для пользователя.
type RouteRating struct {
	UserId            int    `json:"user_id" db:"user_id"`         // UserId - id пользователя, которого представляет блок.
	SprintId          int    `json:"sprint_id" db:"sprint_id"`     // SprintId - id лучшего спринта по маршруту для пользователя.
	SprintLengthTime  int64  `json:"length_time" db:"length_time"` // SprintLengthTime - длительность лучшего спринта по маршруту для пользователя в ms.
	SprintPath        string `json:"-" db:"path"`                  // SprintPath - путь лучшего спринта по маршруту для пользователя.
	SprintLengthSteps int    `json:"length_steps" db:"-"`          // SprintLengthSteps - количество шагов в лучшем спринте по маршруту для пользователя.
}

// TourRating - структура, представляющая блок рейтинга соревнования для пользователя.
type TourRating struct {
	UserName string `json:"user_name"` // UserName - имя пользователя, которого представляет блок.
	Points   int    `json:"-"`         // Points - количество очков пользователя.
}
