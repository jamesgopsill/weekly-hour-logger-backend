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
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/user/register", mockRequestBufferString)
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
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/user/register", mockRequestBufferString)
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
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/user/login", mockRequestBufferString)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response apiResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// This is being assigned to a global var for future use. Check init_test.go
	validUserSignedString = response.Data
	els := strings.Split(response.Data, " ")

	token, err := jwt.ParseWithClaims(els[1], &user.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	assert.NoError(t, err)

	// Got to initialise before because the fcn will try and create two new vars rather than use an existing global var.
	var ok bool
	// This is being assigned to a global var for future use. Check init_test.go
	validUserClaims, ok = token.Claims.(*user.MyCustomClaims)
	assert.True(t, ok)
	assert.NotNil(t, validUserClaims)
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	mockRequest := `{}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/user/update", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+invalidSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateScopesAddAdmin(t *testing.T) {
	mockRequest := `{
		"id": "` + validUserClaims.ID + `",
		"scopes": ["admin"]
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/user/update-scopes", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLoginAdmin(t *testing.T) {
	mockRequest := `{
		"password": "test",
		"email": "test@test.com"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/user/login", mockRequestBufferString)
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var response apiResponse
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// This is being assigned to a global var for future use. Check init_test.go
	validUserSignedString = response.Data
	els := strings.Split(response.Data, " ")

	token, err := jwt.ParseWithClaims(els[1], &user.MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	assert.NoError(t, err)

	// Got to initialise before because the fcn will try and create two new vars rather than use an existing global var.
	var ok bool
	// This is being assigned to a global var for future use. Check init_test.go
	validUserClaims, ok = token.Claims.(*user.MyCustomClaims)
	assert.True(t, ok)
	assert.NotNil(t, validUserClaims)
}

type loginResponse struct {
	Error string
	Data  string
}

func TestUpdateUser(t *testing.T) {
	log.Info().Msg(validUserSignedString)
	mockRequest := `{
		"id": "` + validUserClaims.ID + `",
		"name": "updated name",
		"email": "test@test.com"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/user/update", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRefreshToken(t *testing.T) {
	bufferMockRequest := bytes.NewBufferString("")
	req, err := http.NewRequest("POST", "/user/refresh-token", bufferMockRequest)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListAllUsers(t *testing.T) {
	bufferMockRequest := bytes.NewBufferString("")
	req, err := http.NewRequest("GET", "/user/list-all-users", bufferMockRequest)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateGroup(t *testing.T) {
	mockRequest := `{
		"name": "test group",
		"emails": ["test@test.com"]
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/group/create-group", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateGroupAlreadyExists(t *testing.T) {
	mockRequest := `{
		"name": "test group",
		"emails": ["this_shouldnt_be_here@test.com"]
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/group/create-group", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestAddGroup(t *testing.T) {
	mockRequest := `{
		"name": "test group",
		"emails": ["dbtest2@test.com", "dbtest3@test.com"]
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/group/add-users", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAddGroupUserAlreadyInGroup(t *testing.T) {
	mockRequest := `{
		"name": "test group",
		"emails": ["dbtest2@test.com"]
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/group/add-users", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code == http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// test remove user. Remove 'dbtest2@test.com' from 'test group'.
func TestRemoveGroup(t *testing.T) {
	mockRequest := `{
		"name": "test group",
		"emails": ["dbtest2@test.com"]
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/group/remove-users", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRemoveGroupUserNotInGroup(t *testing.T) {
	mockRequest := `{
		"name": "test group",
		"emails": ["db@test.com"]
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/group/remove-users", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListGroupUsers(t *testing.T) {
	mockRequest := `{
		"name": "test group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("GET", "/group/list-users-in-group", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

// create resource entry
func TestAddResource(t *testing.T) {
	mockRequest := `{
		"week": 7,
		"value": 500,
		"email": "test@test.com",
		"name": "test group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/resource/add-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

// resource - add resource entry to a group - user not in group
func TestAddResourceNotInGroup(t *testing.T) {
	mockRequest := `{
		"week": 9,
		"value": 200,
		"email": "test@test.com",
		"name": "DB Group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/resource/add-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

// resource - Fail: add resource entry to a group - entry already present
func TestAddResourceAlreadyExists(t *testing.T) {
	mockRequest := `{
		"week": 7,
		"value": 200,
		"email": "test@test.com",
		"name": "test group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/resource/add-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// resource - update entry; new value, same week
func TestUpdateResourceSameWeek(t *testing.T) {
	mockRequest := `{
		"oldWeek": 7,
		"newWeek": 7,
		"value": 400,
		"email": "test@test.com",
		"name": "test group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/resource/update-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

// resource - update, new week and new value
func TestUpdateResourceNewWeek(t *testing.T) {

	mockRequest := `{
		"oldWeek": 7,
		"newWeek": 13,
		"value": 800,
		"email": "test@test.com",
		"name": "test group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/resource/update-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

// resource - update, old entry doesn't exist
func TestUpdateResourceOldWeekNotExisting(t *testing.T) {

	mockRequest := `{
		"oldWeek": 9,
		"newWeek": 15,
		"value": 600,
		"email": "test@test.com",
		"name": "test group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/resource/update-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// update resource - fail: would create duplicate entries for a single week
func TestUpdateResourceCreatesDuplicate(t *testing.T) {

	mockRequest := `{
		"oldWeek": 9,
		"newWeek": 13,
		"value": 600,
		"email": "test@test.com",
		"name": "test group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/resource/update-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Delete resource: no entry found to delete
func TestDeleteResourceNoEntry(t *testing.T) {
	mockRequest := `{
		"week": 2,
		"value": 24,
		"email": "dbtest3@test.com",
		"name": "test group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/resource/delete-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteResource(t *testing.T) {
	mockRequest := `{
		"week": 8,
		"value": 500,
		"email": "test@test.com",
		"name": "test group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/resource/add-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var res loginResponse
	err = json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)

	mockRequest = `{
		"week": 8,
		"value": 500,
		"email": "test@test.com",
		"name": "test group"
	}`
	mockRequestBufferString = bytes.NewBufferString(mockRequest)
	req, err = http.NewRequest("POST", "/resource/delete-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

// Delete group: Deletes all resource entries from that group, and removes the link between group and user
func TestDeleteGroup(t *testing.T) {

	mockRequest := `{
		"name": "test group to delete",
		"emails": ["dbtest4@test.com"]
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("POST", "/group/create-group", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var res loginResponse
	err = json.NewDecoder(w.Body).Decode(&res)
	assert.NoError(t, err)

	mockRequest = `{
		"week": 7,
		"value": 300,
		"email": "dbtest4@test.com",
		"name": "test group to delete"
	}`
	mockRequestBufferString = bytes.NewBufferString(mockRequest)
	req, err = http.NewRequest("POST", "/resource/add-resource", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)

	var res2 loginResponse
	err = json.NewDecoder(w.Body).Decode(&res2)
	assert.NoError(t, err)

	mockRequest = `{
		"name": "test group to delete"
	}`
	mockRequestBufferString = bytes.NewBufferString(mockRequest)
	req, err = http.NewRequest("POST", "/group/delete-group", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListGroupResource(t *testing.T) {
	mockRequest := `{
		"name": "DB Group"
	}`
	mockRequestBufferString := bytes.NewBufferString(mockRequest)
	req, err := http.NewRequest("GET", "/group/list-resource-in-group", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListGroups(t *testing.T) {
	mockRequestBufferString := bytes.NewBufferString("")
	req, err := http.NewRequest("GET", "/group/list-groups", mockRequestBufferString)
	assert.NoError(t, err)
	req.Header.Set("Authorization", validUserSignedString)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		log.Info().Msg(w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}
