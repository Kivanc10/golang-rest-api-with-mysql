package dbop

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	var persons = GetAllUsers()
	expected := []Person{{Id: 51, UserName: "person1", Password: "I'll do it myself", Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyNDM0ODMsInVzZXJfaWQiOjAsInVzZXJfbmFtZSI6InBlcnNvbjEifQ.E-3Mkm7WulyKdJ40OiwFsFttabnTTRdobGVygaj0rm0"}, {Id: 52, UserName: "sample user name", Password: "sample password", Token: ""}}
	actual := persons
	assert.Equal(t, expected, actual)
}

func TestGetTodo(t *testing.T) {
	expected := []string{"sample todo"}
	actual := GetTodo("person1")
	assert.Equal(t, expected, actual)
}

func Test_insert(t *testing.T) {
	db := prepareDb(dbname)
	defer db.Close()
	var samplePerson = Person{Id: 0, UserName: "sample user name", Password: "sample password", Token: ""}
	_, err := insert(db, samplePerson)
	assert.NoError(t, err)
}

func Test_getPersonToLogin(t *testing.T) {
	expectedPerson := Person{Id: 50, UserName: "sample user name", Password: "sample password", Token: ""}
	actualPerson := getPersonToLogin(expectedPerson.UserName, expectedPerson.Password)
	assert.Equal(t, expectedPerson, actualPerson)
}

func Test_deleteUserById(t *testing.T) {
	db := prepareDb(dbname)
	defer db.Close()
	_, err := deleteUserById(db, 52)
	assert.NoError(t, err)
}

func Test_deleteAllTokens(t *testing.T) {
	db := prepareDb(dbname)
	defer db.Close()
	p := Person{Id: 50, UserName: "sample user name", Password: "sample password", Token: ""}
	_, err := deleteAllTokens(&p)
	assert.NoError(t, err)
}

func Test_deleteAllTodos(t *testing.T) {
	db := prepareDb(dbname)
	defer db.Close()
	p := Person{Id: 50, UserName: "sample user name", Password: "sample password", Token: ""}
	_, err := deleteAllTodos(&p)
	assert.NoError(t, err)
}

func TestLogOutFromAllSession(t *testing.T) {
	db := prepareDb(dbname)
	defer db.Close()
	_, err := LogOutFromAllSession("sample user name")
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	p := Person{Id: 52, UserName: "updated name", Password: "updated pswrd", Token: ""}
	reqBody, err := json.Marshal(p)
	assert.NoError(t, err)

	_, err = UpdateUser(reqBody, "sample user name")
	assert.NoError(t, err)
}

func TestGetMe(t *testing.T) {
	expected := Person{51, "person1", "I'll do it myself", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyNDM0ODQsInVzZXJfaWQiOjUxLCJ1c2VyX25hbWUiOiJwZXJzb24xIn0.YyKwqSRLtAXcDWEsxA0eWN-y5k3MiDfjiKkbrWWA1D4"}
	actual := GetMe(expected.UserName)
	assert.Equal(t, expected, actual)
}

func TestGetLastLoginToken(t *testing.T) {
	type testcase struct {
		username      string
		expectedToken string
		actualToken   string
	}
	users := []testcase{
		{"sample user name", "", ""},
		{"person1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyNDM0ODQsInVzZXJfaWQiOjUxLCJ1c2VyX25hbWUiOiJwZXJzb24xIn0.YyKwqSRLtAXcDWEsxA0eWN-y5k3MiDfjiKkbrWWA1D4", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyNDM0ODQsInVzZXJfaWQiOjUxLCJ1c2VyX25hbWUiOiJwZXJzb24xIn0.YyKwqSRLtAXcDWEsxA0eWN-y5k3MiDfjiKkbrWWA1D4"},
	}
	for _, user := range users {
		if user.actualToken == GetLastLoginToken(user.username) && user.expectedToken == user.actualToken {

		} else {
			t.Error("unappropriate token for user")
		}
	}
	/*
		token := GetLastLoginToken("person1")
		expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyNDM0ODQsInVzZXJfaWQiOjUxLCJ1c2VyX25hbWUiOiJwZXJzb24xIn0.YyKwqSRLtAXcDWEsxA0eWN-y5k3MiDfjiKkbrWWA1D4"
		actual := token
		assert.Equal(t, expected, actual)*/
}
