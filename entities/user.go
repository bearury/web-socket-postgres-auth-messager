package entities

type User struct {
	Id       string `json:"-" db:"id"`               // json:"-" - при сериализации в JSON поле игнорируется
	Name     string `json:"name" binding:"required"` // binding:"required" Это тег фреймворка Gin
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
