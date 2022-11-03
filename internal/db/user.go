package db

type User struct {
	Base
	Name         string
	Email        string `gorm:"uniqueIndex"`
	Scopes       SerialisableStringArray
	PasswordHash string `json:"-"`
	GroupID      string
	Group        Group      // added this. Group ID above is the key, this is the inverse signpost?
	Resources    []Resource `gorm:"foreignKey:UserID"`
}
