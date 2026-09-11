package entities

type ChatType string

const (
	ChatTypeChat  ChatType = "chat"
	ChatTypeGroup ChatType = "group"
	ChatTypeChain ChatType = "chain"
)

type Chat struct {
	Id   string   `json:"-" db:"id"` // json:"-" - при сериализации в JSON поле игнорируется
	Name string   `json:"name" db:"name" binding:"required"`
	Type ChatType `json:"type" db:"type" binding:"required"`
}
