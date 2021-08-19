package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
	/* you should run the program first , go run ./operate/operate.go */
	//log.Fatal(http.ListenAndServe(":8080", nil))
	// create a new sample http request like ping then check if it able logged in with token or not
	testFunc := func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyODc0NDksInVzZXJfaWQiOjUxLCJ1c2VyX25hbWUiOiJwZXJzb24xIn0.4B2IVsenu2swaCHZM4aKV9mGUl-P-hk-5E4EeReYBco")
		w.Write([]byte("welcome to main test func"))
	}
	// get me func --> you will not access to yourself becuase your token is expired
	req, err := http.NewRequest("GET", "http://localhost:8080/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testFunc)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code : got %v want %v\n", status, http.StatusOK)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		t.Errorf("the http test request failed with err %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		expected := []byte("you are Unauthorized or your token is expired")
		if !bytes.Equal(data, expected) {
			t.Errorf("handler returned unexpected body : got %v want %v", data, expected)
		}
	}
}
