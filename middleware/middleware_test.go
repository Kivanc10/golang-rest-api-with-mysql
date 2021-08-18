package middleware

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	var username, id = "john doe", 15

	//var samplePerson = dbop.Person{Id: 15, UserName: "john doe", Password: "john_pswrd"}
	token, err := CreateToken(uint64(id), username)
	assert.NoError(t, err)

	if assert.NotEqual(t, token, "") {
		//assert.Equal(t, err, errors.New("an error occured during the create token"))
		t.Log("succeed to create token for the user")
	}
}

func TestMiddleWare(t *testing.T) {

	//log.Fatal(http.ListenAndServe(":8080", nil))
	// create a new sample http request like ping then check if it able logged in with token or not
	testMainFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyODc0NDksInVzZXJfaWQiOjUxLCJ1c2VyX25hbWUiOiJwZXJzb24xIn0.4B2IVsenu2swaCHZM4aKV9mGUl-P-hk-5E4EeReYBco")
		w.Write([]byte("welcome to main test func"))
	})
	http.Handle("/test/main", MiddleWare(testMainFunc))
	log.Fatal(http.ListenAndServe(":9090", nil))
}
