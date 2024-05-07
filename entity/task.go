package entity

type Task struct {
	ID         int
	Title      string
	Dodate     string
	CategoryID int
	Isdone     bool
	UserID     int
}
