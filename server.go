package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Supply struct {
	ID         int     `json:"id"`
	Barcode    string  `json:"barcode"`
	Quantity   int     `json:"quantity"`
	SupplyTime string  `json:"supply_time"`
	Price      float64 `json:"price"`
	SoldAmount float64 `json:"sold_amount"`
}

type Sale struct {
	ID       int     `json:"id"`
	Barcode  string  `json:"barcode"`
	Quantity int     `json:"quantity"`
	SaleTime string  `json:"sale_time"`
	Price    float64 `json:"price"`
	Margin   float64 `json:"margin"`
}

var db *sql.DB

func main() {
	var err error
	var dsn string
	dsn = "urbvbmxcu9morjp7:E8CMsE7fDdxDHPCcPOh9@tcp(bqm4pgy45ritx2m2zkgn-mysql.services.clever-cloud.com:3306)/bqm4pgy45ritx2m2zkgn?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/sales", getSales).Methods("GET")
	router.HandleFunc("/sales/{id}", getSale).Methods("GET")
	router.HandleFunc("/sales", postSale).Methods("POST")
	router.HandleFunc("/sales/{id}", patchSale).Methods("PATCH")
	router.HandleFunc("/sales/{id}", deleteSale).Methods("DELETE")

	router.HandleFunc("/supplies", getSupplies).Methods("GET")
	router.HandleFunc("/supplies/{id}", getSupply).Methods("GET")
	router.HandleFunc("/supplies", postSupply).Methods("POST")
	router.HandleFunc("/supplies/{id}", patchSupply).Methods("PATCH")
	router.HandleFunc("/supplies/{id}", deleteSupply).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func getSales(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, barcode, quantity, sale_time, price, margin FROM sale")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var sales []Sale
	for rows.Next() {
		var s Sale
		if err := rows.Scan(&s.ID, &s.Barcode, &s.Quantity, &s.SaleTime, &s.Price, &s.Margin); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sales = append(sales, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sales)
}

func getSale(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var s Sale
	err := db.QueryRow("SELECT id, barcode, quantity, sale_time, price, margin FROM sale WHERE id=$1", id).Scan(&s.ID, &s.Barcode, &s.Quantity, &s.SaleTime, &s.Price, &s.Margin)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func postSale(w http.ResponseWriter, r *http.Request) {
	var s Sale
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &s)

	sqlStatement := "INSERT INTO sale (barcode, quantity, sale_time, price, margin) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := db.QueryRow(sqlStatement, s.Barcode, s.Quantity, s.SaleTime, s.Price, s.Margin).Scan(&s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func patchSale(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var s Sale
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &s)

	sqlStatement := "UPDATE sale SET barcode=$1, quantity=$2, sale_time=$3, price=$4, margin=$5 WHERE id=$6"
	res, err := db.Exec(sqlStatement, s.Barcode, s.Quantity, s.SaleTime, s.Price, s.Margin, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteSale(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	sqlStatement := "DELETE FROM sales WHERE id=$1"
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getSupplies(w http.ResponseWriter, r *http.Request) {
	var supplies []Supply
	rows, err := db.Query("SELECT * FROM supply")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var s Supply
		err = rows.Scan(&s.ID, &s.Barcode, &s.Quantity, &s.SupplyTime, &s.Price, &s.SoldAmount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		supplies = append(supplies, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supplies)
}

func getSupply(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var s Supply
	err := db.QueryRow("SELECT * FROM supply WHERE id = $1", id).Scan(&s.ID, &s.Barcode, &s.Quantity, &s.SupplyTime, &s.Price, &s.SoldAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func postSupply(w http.ResponseWriter, r *http.Request) {
	var s Supply
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &s)

	sqlStatement := "INSERT INTO supply (barcode, quantity, supply_time, price, sold_amount) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err := db.QueryRow(sqlStatement, s.Barcode, s.Quantity, s.SupplyTime, s.Price, s.SoldAmount).Scan(&s.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func patchSupply(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var s Supply
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &s)

	sqlStatement := "UPDATE supply SET barcode = $1, quantity = $2, supply_time = $3, price = $4, sold_amount = $5 WHERE id = $6"
	result, err := db.Exec(sqlStatement, s.Barcode, s.Quantity, s.SupplyTime, s.Price, s.SoldAmount, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteSupply(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	sqlStatement := "DELETE FROM supply WHERE id = $1"
	result, err := db.Exec(sqlStatement, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
