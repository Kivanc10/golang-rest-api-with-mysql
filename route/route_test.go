package route

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	dbop "kivancaydogmus.com/apps/userApp/dbOp"
)

func Test_getUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getUsers)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code : got %v want %v\n", status, http.StatusOK)
	}
	// expectedFirst := []dbop.Person{
	// 	{Id: 51, UserName: "person1", Password: "I'll do it myself", Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyNDM0ODMsInVzZXJfaWQiOjAsInVzZXJfbmFtZSI6InBlcnNvbjEifQ.E-3Mkm7WulyKdJ40OiwFsFttabnTTRdobGVygaj0rm0"},
	// 	{Id: 52, UserName: "sample user name", Password: "sample password", Token: ""},
	// }
	//expectedLat, _ := json.Marshal(expectedFirst)

	expectedLat := `[{"PersonID":51,"UserName":"person1","Password":"I'll do it myself","Token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE2MjkyNDM0ODMsInVzZXJfaWQiOjAsInVzZXJfbmFtZSI6InBlcnNvbjEifQ.E-3Mkm7WulyKdJ40OiwFsFttabnTTRdobGVygaj0rm0"},{"PersonID":52,"UserName":"sample user name","Password":"sample password","Token":""}]`

	if rr.Body.String() != expectedLat {
		t.Errorf("handler returned unexpected body : got %v want %v", rr.Body.String(), expectedLat)
	}
}

func Test_addUser(t *testing.T) {
	samplePerson := dbop.Person{Id: 1, UserName: "username", Password: "password"}
	bytePerson, _ := json.Marshal(samplePerson)

	req, err := http.NewRequest("POST", "/signUp", bytes.NewReader(bytePerson))

	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	//rr.Body = bytes.NewBuffer(bytePerson)
	handler := http.HandlerFunc(addUser)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code : got %v want %v\n", status, http.StatusOK)
	}

	req.Header.Set("Token", samplePerson.Token)

}

func Test_login(t *testing.T) {
	samplePerson := dbop.Person{Id: 1, UserName: "username", Password: "password"}
	bytePerson, _ := json.Marshal(samplePerson)
	req, err := http.NewRequest("POST", "/signIn", bytes.NewReader(bytePerson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(login)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code : got %v want %v\n", status, http.StatusOK)
	}
	req.Header.Set("Token", samplePerson.Token)
}
