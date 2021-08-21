# RESTful API in GO

This is an intermediate-level go project that running with a project structure optimized RESTful API service in Go. API's of that project is designed based on solid and common principles and connected to the local MySQL database.

<b>Highlights of that project are listed at below </b>
  - The RESTful API presents standard CRUD operations of a database table
  - This project has clean architecture and it has been covered with tests.
  - Data validation
  - It uses JWT-based authentication and auth middleware.
  - Each token is expired in 15 minutes to prevent system bugs.
  - Error handling is done with clear responses.
  - It presents structured logging with the username, password for 15 minutes with help of a token that is created by JWT. After that, the user can create TODOs. 
  - <b> The project used the following packages during the development time </b>
      - Routing : <a href="github.com/gorilla/mux">Mux</a>
      - Database driver : <a href="github.com/go-sql-driver/mysql">Mysql</a>
      - JWT : <a href="github.com/dgrijalva/jwt-go">go-jwt</a>
      - Test : <a href="github.com/stretchr/testify">Testify</a>

## Getting Started

If you're have not encountered Go before, you should visit this website <a target="_blank" href="https://golang.org/doc/install">here</a>

After installing Go , you should run the following commands to experience this project

```cmd
# download the starter code
git clone https://github.com/Kivanc10/golang-rest-api-with-mysql.git

# open the code
cd golang-rest-api-with-mysql

# start the database server and run the code
go run ./operate/operate.go
```
After that, you have a RESTful API that is running at `http://127.0.0.1:8080`. It provides us following endpoints
  - `GET /users` : it provides us the list of all users logged-in
  - `POST /signUp` : it allows the user to sign up. It saves the user info into db and creates token with JWT.It accepts attached data like that:
    - `PersonID` can be anything because db arranged to auto increment
    - ```JSON
      {
          "PersonID":0,
          "UserName":"sample user name",
          "Password":"12312321"
       }
      ```
  - `POST /signIn` : authenticates and login. It creates token again with JWT.It accepts attached data to , something like up above.
  - `GET /users/me` : It allows the user to access his information.To do this , user must be authenticated,otherwise system wil not let that happen
  - `PUT /users/update/me` : It updates the current authenticated user with accepted data.To do this user must be authenticated.It accepts attached data like:
      - ```JSON
        {
              "PersonID":0,
              "UserName":"new user name",
              "Password":"new password"
        }
        ```
  - `DELETE /user/me` : It deletes the current authenticated user.To do this user must be authenticated.
  - `GET /users/logout/me` : It allows the user to logout from all tokens.To do this user must be authenticated.The user will be not deleted from db.
  - `POST /todo` : It allows the user to create todos. To do this user must be authenticated.It accepts attached data like that :
      - ```JSON
            {
                "Context" : "sample todo"
             }
          ```
  - `GET /todos` : It lists all todos that created by authenticated users.
  - `GET /todos/me` : It lists just the todos belong to the current user authenticated


### If you have API client tools like `POSTMAN`, you can handle complicated operations easily

```
# sign up the user via POST /signUp
curl -X POST -H "Content-Type: application/json" -d `{ "PersonID":0,"UserName":"sample user name","Password":"12312321"}` http://localhost:8080/signUp
# it should return response.header with jwt and token

# sign in with user via POST /signIn
curl -X POST -H "Content-Type: application/json" -d `{ "PersonID":0,"UserName":"sample user name","Password":"12312321"}` http://localhost:8080/signIn
# it should return response.Header with jwt and token
# save token during the loggedin and inherit auth from postman environment,it handles itself
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8080/users/me

# to create todos for the user authenticated
curl -X POST -H "Authorization: Bearer ...JWT token here..." -d `{"Context" : "sample todo"}` http://localhost:8080/todo
# it returns the saved todo belong to the user authenticated
```


