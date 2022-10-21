package db

type Resource struct {
	Base
	Week    uint32
	Value   uint32
	UserID  string
	User    User // added this. Inverse pointer for the database
	GroupID string
	Group   Group // added this. Inverse pointer for the database
}
