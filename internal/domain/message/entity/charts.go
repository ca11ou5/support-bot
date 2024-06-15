package entity

type ChartData struct {
	ID     string // Уникальный идентификатор для div элемента
	Option interface{}
}

type PageData struct {
	Charts []ChartData
}
