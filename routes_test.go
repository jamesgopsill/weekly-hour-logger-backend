package main

import (
	"bytes"
	"encoding/json"
	"jamesgopsill/resource-logger-backend/internal/controllers/user"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegister(t *testing.T) {
	mockRequest := `{
		"name": "test",
		"email" : "test@test.com",
		"confirmEmail" : "test@test.com",
		"password" : "test",
		"confirmPassword" : "test"
	}`
	req, err := http.NewRequest("POST", "/user/register", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterAccountExists(t *testing.T) {
	mockRequest := `{
		"name": "test",
		"email" : "test@test.com",
		"confirmEmail" : "test@test.com",
		"password" : "test",
		"confirmPassword" : "test"
	}`
	req, err := http.NewRequest("POST", "/user/register", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnprocessableEntity {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestLogin(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	mockRequest := `{}`
	req, err := http.NewRequest("POST", "/user/update", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+invalidSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

type loginResponse struct {
	Error string
	Data  string
}

func TestUpdateUser(t *testing.T) {
	var mockRequest string
	mockRequest = `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	els := strings.Split(response.Data, " ")

	token, err := jwt.ParseWithClaims(els[1], &user.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	assert.NoError(t, err)
	claims, ok := token.Claims.(*user.MyCustomClaims)
	assert.Equal(t, true, ok)

	mockRequest = `{
		"id": "` + claims.ID + `",
		"name": "updated name",
		"email": "test@test.com"
	}`
	req, err = http.NewRequest("POST", "/user/update", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRefreshToken(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	req, _ = http.NewRequest("POST", "/user/refresh-token", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListAllUsers(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	req, _ = http.NewRequest("GET", "/user/list-all-users", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateGroup(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"name": "test group",
		"emails": ["test@test.com"]
	}`
	req, _ = http.NewRequest("POST", "/group/create-group", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateGroupAlreadyExists(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"name": "test group",
		"emails": ["this_shouldnt_be_here@test.com"]
	}`
	req, _ = http.NewRequest("POST", "/group/create-group", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestAddGroup(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"name": "test group",
		"emails": ["dbtest2@test.com", "dbtest3@test.com"]
	}`
	req, _ = http.NewRequest("POST", "/group/add-users", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAddGroupUserAlreadyInGroup(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"name": "test group",
		"emails": ["dbtest2@test.com"]
	}`
	req, _ = http.NewRequest("POST", "/group/add-users", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// test remove user. Remove 'dbtest2@test.com' from 'test group'.
func TestRemoveGroup(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"name": "test group",
		"emails": ["dbtest2@test.com"]
	}`
	req, _ = http.NewRequest("POST", "/group/remove-users", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRemoveGroupUserNotInGroup(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"name": "test group",
		"emails": ["db@test.com"]
	}`
	req, _ = http.NewRequest("POST", "/group/remove-users", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListGroupUsers(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"name": "test group"
	}`
	req, _ = http.NewRequest("POST", "/group/list-users-in-group", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

// create resource entry
func TestAddResource(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"week": 7,
		"value": 500,
		"email": "test@test.com",
		"name": "test group"
	}`
	req, _ = http.NewRequest("POST", "/resource/add-resource", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

// resource - add resource entry to a group - user not in group
func TestAddResourceNotInGroup(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"week": 9,
		"value": 200,
		"email": "test@test.com",
		"name": "DB Group"
	}`
	req, _ = http.NewRequest("POST", "/resource/add-resource", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// resource - Fail: add resource entry to a group - entry already present
func TestAddResourceAlreadyExists(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"week": 7,
		"value": 200,
		"email": "test@test.com",
		"name": "test group"
	}`
	req, _ = http.NewRequest("POST", "/resource/add-resource", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// resource - update entry; new value, same week
func TestUpdateResourceSameWeek(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"oldWeek": 7,
		"newWeek": 7,
		"value": 400,
		"email": "test@test.com",
		"name": "test group"
	}`
	req, _ = http.NewRequest("POST", "/resource/update-resource", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

// resource - update, new week and new value
func TestUpdateResourceNewWeek(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"oldWeek": 7,
		"newWeek": 13,
		"value": 800,
		"email": "test@test.com",
		"name": "test group"
	}`
	req, _ = http.NewRequest("POST", "/resource/update-resource", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

// resource - update, old entry doesn't exist
func TestUpdateResourceOldWeekNotExisting(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"oldWeek": 9,
		"newWeek": 15,
		"value": 600,
		"email": "test@test.com",
		"name": "test group"
	}`
	req, _ = http.NewRequest("POST", "/resource/update-resource", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// update resource - fail: would create duplicate entries for a single week
func TestUpdateResourceCreatesDuplicate(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBufferString(mockRequest))
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response loginResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	mockRequest = `{
		"oldWeek": 9,
		"newWeek": 13,
		"value": 600,
		"email": "test@test.com",
		"name": "test group"
	}`
	req, _ = http.NewRequest("POST", "/resource/update-resource", bytes.NewBufferString(mockRequest))
	req.Header.Set("Authorization", response.Data)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
