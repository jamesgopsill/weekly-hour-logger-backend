package main

import (
	"jamesgopsill/resource-logger-backend/internal/db"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm/clause"
)

func TestDB(t *testing.T) {

	log.Info().Msg("hello")

	user := db.User{
		Name:         "DB Test",
		Email:        "db@test.com",
		PasswordHash: "db_hash",
		Scopes:       []string{"user", "admin"},
	}

	result := db.Connection.Create(&user)
	assert.NoError(t, result.Error)

	log.Info().Msg(user.ID)

	users := []db.User{user}

	group := &db.Group{
		Name:  "DB Group",
		Users: users,
	}

	result = db.Connection.Create(&group)
	assert.NoError(t, result.Error)

	user = db.User{
		Name:         "DB Test User 2",
		Email:        "dbtest2@test.com",
		PasswordHash: "db_hash",
		Scopes:       []string{"user", "admin"},
	}

	result = db.Connection.Create(&user)
	assert.NoError(t, result.Error)

	user = db.User{
		Name:         "DB Test User 3",
		Email:        "dbtest3@test.com",
		PasswordHash: "db_hash",
		Scopes:       []string{"user", "admin"},
	}

	result = db.Connection.Create(&user)
	assert.NoError(t, result.Error)

	user = db.User{
		Name:         "DB Test User 4 - no group",
		Email:        "dbtest4@test.com",
		PasswordHash: "db_hash",
		Scopes:       []string{"user", "admin"},
	}

	result = db.Connection.Create(&user)
	assert.NoError(t, result.Error)

	resource := &db.Resource{
		Week:    1,
		Value:   18,
		UserID:  user.ID,
		GroupID: group.ID,
	}

	result = db.Connection.Create(&resource)
	assert.NoError(t, result.Error)

	resource = &db.Resource{
		Week:    2,
		Value:   24,
		UserID:  user.ID,
		GroupID: group.ID,
	}

	result = db.Connection.Create(&resource)
	assert.NoError(t, result.Error)

	db.Connection.Preload(clause.Associations).Find(&user)

	for _, r := range user.Resources {
		log.Info().Msg(r.ID)
	}

}
