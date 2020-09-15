package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-postgres-crud/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "root"
	dbname   = "users"
)

func createConnection() *sql.DB {
	// err := godotenv.Load(".env")
	psqlQuery := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlQuery)

	// if err != nil {
	// 	log.Fatalf("Error loading .env file", err)
	// 	panic(err)
	// }

	// db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err.Error)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully Connected")
	return db

}

//create user

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// creating empty user of type models from models.User
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("unable to decode the request body", err)
		panic(err)
	}

	insertId := insertUser(user)

	res := response{
		ID:      insertId,
		Message: "user created Successfully",
	}

	// send the response
	json.NewEncoder(w).Encode(res)

}

// get user - single user data
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to get user. %v", err)
	}

	json.NewEncoder(w).Encode(id)
}

// get All users - data

func GetAllUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// get all users from DB
	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to get all user. %v", err)
	}

	json.NewEncoder(w).Encode(users)

}

func getAllUsers() ([]models.User, error) {
	db := createConnection()

	defer db.Close()

	var users []models.User

	// select query
	sqlStatement := `SELECT * FROM users`

	// execute the query
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the statement", err)
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		users = append(users, user)
	}

	// return empty user on error
	return users, err

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Contenet-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	var user models.User

	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	updatedRows := updatingUser(int64(id), user)

	msg := fmt.Sprintf("User updated successfully. Total record affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

// update user in the DB
func updatingUser(id int64, user models.User) int64 {
	db := createConnection()
	defer db.Close()

	sqlQuery := `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlQuery, id, user.Name, user.Location, user.Age)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("err while checking the affected rows. %v", err)
	}

	return rowsAffected
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	deleteedRows := deletingUser(int64(id))

	msg := fmt.Sprintf("User Updated Successfully. Total row(s) Affected %v", deleteedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func deletingUser(id int64) int64 {
	db := createConnection()

	defer db.Close()

	sqlQuery := `DELETE FROM users WHERE userid=$1`

	res, err := db.Exec(sqlQuery, id)
	resPrint := &res
	fmt.Println("res=============", *resPrint)
	if err != nil {
		log.Fatalf("Unable to execute the query .%v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Println("Total row affected %v", rowsAffected)
	return rowsAffected

}

func insertUser(user models.User) int64 {
	//create the postgres db connection
	db := createConnection()

	// close db
	defer db.Close()

	sqlStatement := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`
	var id int64

	err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query .%v", err)
	}

	return id
}
