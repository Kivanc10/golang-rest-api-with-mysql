package dbop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
	var persons = GetAllUsers()
	expected := []Person{{Id: 48, UserName: "person1", Password: "I'll do it myself", Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyMzk4MTgsInVzZXJfaWQiOjAsInVzZXJfbmFtZSI6InBlcnNvbjEifQ.Ti_ZzS1c9Sp45vmXAMpzIzX09azqEI77KxGDQetNMr8"}, {Id: 49, UserName: "sample user name", Password: "sample password", Token: ""}}
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
	expectedPerson := Person{Id: 49, UserName: "sample user name", Password: "sample password", Token: ""}
	actualPerson := getPersonToLogin(expectedPerson.UserName, expectedPerson.Password)
	assert.Equal(t, expectedPerson, actualPerson)
}
