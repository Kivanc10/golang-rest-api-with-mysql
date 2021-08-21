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

