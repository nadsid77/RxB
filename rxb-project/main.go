
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5555
	user     = "postgres"
	password = "postgres"
	dbname   = "dvdrental"
)


type Film struct {
	FilmID           int
	Title            string  `json:"Title,omitempty"`
	Category         string  `json:"Category,omitempty"`
	Rating           string  `json:"Rating,omitempty"`
	Language         string  `json:"Language,omitempty"`
	Description      string  `json:"Description,omitempty"`
	ReleaseYear      int     `json:"ReleaseYear,omitempty"`
	ActorFirstName   string  `json:"ActorFirstName,omitempty"`
	ActorLastName    string  `json:"ActorLastName,omitempty"`
	Length           int     `json:"Length,omitempty"`
	RentalDuration   int     `json:"RentalDuration,omitempty"`
	RentalRate       float32 `json:"RentalRate,omitempty"`
	ReplacementCost  float32 `json:"ReplacementCost,omitempty"`
	SpecialFeatures string  `json:"SpecialFeatures,omitempty"`
	CustomerId       int     `json:"CustomerId,omitempty"`
	CommentId        int     `json:"CommentId,omitempty"`
	Comment          string  `json:"comment,omitempty"`
}

type Customer struct {
	CustomerID  int
	FirstName string  `json:"FirstName,omitempty"`
	LastName string `json:"last_name,omitempty"`
}

type Comment struct {
	CommentId int
	CustomerID int
	FilmID     int
	Comment    string `json:"comment,omitempty"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/Ping", PingHandler).Methods("GET")
	r.HandleFunc("/", WelcomeHandler).Methods("GET")
	r.HandleFunc("/customer", GetCustomersHandler).Methods("GET")
	r.HandleFunc("/films", GetFilmsHandler).Methods("GET")
	r.HandleFunc("/films/ratings/{rating}", GetRatingHandler).Methods("GET")
	r.HandleFunc("/films/categories/{category}", GetCategoryHandler).Methods("GET")
	r.HandleFunc("/films/titles/{title}", GetFilmDetailHandler).Methods("GET")
	r.HandleFunc("/films/comment", PostCommentHandler).Methods("POST")
	r.HandleFunc("/films/{film_id}/comment/{customer_id}", GetCommentHandler).Methods("GET")

	http.ListenAndServe(":8080", r)
}

func connect() *sql.DB {

	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", conn)

	if err != nil {
		panic(err)
	}
	err = db.Ping()

	return db
}

func GetFilmsHandler(w http.ResponseWriter, r *http.Request) {

	db := connect()

	rows, err := db.Query(
		`SELECT DISTINCT ON (f.film_id)
			f.film_id,
			f.title,
			f.description,
			f.release_year,
       		l.name,
			f.rental_duration,
			f.rental_rate,
     		f.Length,
			f.replacement_cost,
       		f.rating,
			f.special_features
		FROM film f, language l
		WHERE 
			l.language_id = f.language_id `)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var films []*Film

	for rows.Next() {

		flm := new(Film)
		//go type
		rows.Scan(&flm.FilmID, &flm.Title, &flm.Description, &flm.ReleaseYear, &flm.Language,  &flm.RentalDuration, &flm.RentalRate, &flm.Length, &flm.ReplacementCost, &flm.Rating, &flm.SpecialFeatures)

		films = append(films, flm)
	}
	//marshalling
	res, _ := json.MarshalIndent(films, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func GetCustomersHandler(w http.ResponseWriter, r *http.Request) {

	db := connect()

	rows, err := db.Query(
		`SELECT customer_id, first_name, last_name FROM customer c`)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	var customers []*Customer

	for rows.Next() {

	customer := new(Customer)

	//go type
	rows.Scan(&customer.CustomerID, &customer.FirstName, &customer.LastName)

	customers = append(customers, customer)
	}

	//marshalling
	res, _ := json.MarshalIndent(customers, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func GetRatingHandler(w http.ResponseWriter, r *http.Request) {

	db := connect()

	vars := mux.Vars(r)

	ratings := vars["rating"]

	rows, err := db.Query(
		`SELECT DISTINCT ON (f.film_id)
			f.film_id,
			f.rating,
    		f.title,
			f.description,
			f.release_year,
			f.Length,
			l.name,
			f.rental_duration,
			f.rental_rate,
			f.replacement_cost,
			f.special_features
		FROM film f, language l
		WHERE 
			l.language_id = f.language_id AND
		rating=$1`, ratings)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var films []*Film

	for rows.Next() {

		flm := new(Film)

	//go type
		rows.Scan(&flm.FilmID,  &flm.Rating, &flm.Title, &flm.Description, &flm.ReleaseYear, &flm.Length, &flm.Language, &flm.RentalDuration, &flm.RentalRate, &flm.ReplacementCost, &flm.SpecialFeatures)

		films = append(films, flm)
	}

	//marshalling
	res, _ := json.MarshalIndent(films, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func GetCategoryHandler(w http.ResponseWriter, r *http.Request) {

	db := connect()

	vars := mux.Vars(r)

	categories := vars["category"]

	rows, err := db.Query(
		`SELECT DISTINCT ON (f.film_id)
			f.film_id,
    		c.name,
			f.title,
			f.rating,
			f.description,
			f.release_year,
			f.Length,
			l.name,
			f.rental_duration,
			f.rental_rate,
			f.replacement_cost,
			f.special_features
		FROM film f, category c, film_category fc, language l
		WHERE 
			fc.film_id = f.film_id AND 
			c.category_id = fc.category_id AND
			l.language_id = f.language_id AND
			c.name=$1`, categories)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var films []*Film

	for rows.Next() {

		flm := new(Film)

		rows.Scan(&flm.FilmID, &flm.Category, &flm.Title, &flm.Rating, &flm.Description, &flm.ReleaseYear, &flm.Length, &flm.Language, &flm.RentalDuration, &flm.RentalRate, &flm.ReplacementCost, &flm.SpecialFeatures)

		films = append(films, flm)
	}

	//marshalling
	res, _ := json.MarshalIndent(films, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func GetFilmDetailHandler(w http.ResponseWriter, r *http.Request) {

	db := connect()
	vars := mux.Vars(r)

	title := vars["title"]


	rows, err := db.Query(
		`SELECT DISTINCT ON (f.film_id)
			f.film_id,
			f.title,
			f.rating,
			f.description,
			f.release_year,
			f.Length,
			a.first_name,
			a.last_name,
			l.name,
			c.name,
			f.rental_duration,
			f.rental_rate,
			f.replacement_cost,
			f.special_features
		FROM film f, category c, film_category fc, language l, actor a, film_actor fa
		WHERE 
			fc.film_id = f.film_id AND 
			c.category_id = fc.category_id AND
			l.language_id = f.language_id AND
			fa.film_id = f.film_id AND
			a.actor_id = fa.actor_id AND f.title=$1`, title)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var films []*Film

	for rows.Next() {

		flm := new(Film)

		rows.Scan(&flm.FilmID, &flm.Title, &flm.Rating, &flm.Description, &flm.ReleaseYear, &flm.Length, &flm.ActorFirstName, &flm.ActorLastName, &flm.Language, &flm.Category, &flm.RentalDuration, &flm.RentalRate, &flm.ReplacementCost, &flm.SpecialFeatures)

		films = append(films, flm)
	}

	//marshall
	res, _ := json.MarshalIndent(films, "", "	")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {

	db := connect()

	var cmt Film

	//decode post args point it to cmt
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cmt)
	if err != nil {
		log.Fatal(err)
	}


	comment_id := 0
	comment := cmt.Comment
	customer_id := cmt.CustomerId
	film_id := cmt.FilmID


	err = db.QueryRow(
		`INSERT INTO comment(
			comment,
			customer_id,
			film_id)
		VALUES($1,$2,$3) RETURNING comment_id`, comment, customer_id, film_id).Scan(&comment_id)

	if err != nil {
		log.Fatal(err)
	}

	//marshall
	var message = strconv.Itoa(customer_id) + "Your Comment has been posted!" + comment+ "commentId" + strconv.Itoa(comment_id)
	res, _ := json.Marshal(WelcomeResponse{Message: message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func GetCommentHandler(w http.ResponseWriter, r *http.Request) {

	db := connect()
	vars := mux.Vars(r)

	film_id := vars["film_id"]
	customer_id := vars["customer_id"]

	rows, err := db.Query(
		`SELECT c.comment_id, c.comment, c.film_id, c.customer_id 
		FROM comment c WHERE
		c.film_id=$1 AND c.customer_id=$2`,  film_id, customer_id)

	if err != nil {
		log.Fatal(err)
	}

	var comment []*Comment
	for rows.Next() {

		cmt := new(Comment)

		rows.Scan(&cmt.CommentId, &cmt.Comment, &cmt.FilmID, &cmt.CustomerID)

		comment = append(comment, cmt)
	}

	//marshall
	res, _ := json.MarshalIndent(comment, "", "		")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}
func PingHandler(w http.ResponseWriter, r *http.Request){
	res, _ := json.MarshalIndent(PingResponse{Message: "Pong"}, "\n", " " )

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := json.MarshalIndent(WelcomeResponse{Message: "Mockbuster Portal "}, "\n", "	")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

type PingResponse struct{
	Message string `json:"message"`
}

type WelcomeResponse struct {
	Message string `json:"message"`
}
