package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"chubb.com/m/dto"
	"github.com/gorilla/mux"
)

var certs []dto.Cert

// Get all certs
func GetCerts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(certs)
}

// Get cert by ID
func GetCert(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range certs {
		if item.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Certificate not found", http.StatusNotFound)
}

// Create new cert
func CreateCert(w http.ResponseWriter, r *http.Request) {
	var cert dto.Cert
	if err := json.NewDecoder(r.Body).Decode(&cert); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Prepare request to currencyconverter
	var currencyReq dto.CurrencyRequest
	currencyReq.Amount = cert.PremiumAmount
	currencyReq.CountryCode = cert.Contry

	jsonReq, err := json.Marshal(currencyReq)
	if err != nil {
		http.Error(w, "Failed to encode currency request", http.StatusInternalServerError)
		return
	}

	// Make the POST call to currencyconverter
	resp, err := http.Post("http://localhost:8081/currencyconverter", "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		http.Error(w, "Currency conversion service failed", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse the currencyconverter response
	var currencyResponse dto.CurrencyResponse
	if err := json.NewDecoder(resp.Body).Decode(&currencyResponse); err != nil {
		http.Error(w, "Failed to parse currency response", http.StatusInternalServerError)
		return
	}

	// Set the converted values in cert
	cert.PremiumAmount = currencyResponse.ConvertedAmount
	cert.CurrencyType = currencyResponse.CurrencyCode // or currencyResponse.CurrencyName if preferred
	cert.ID = strconv.Itoa(rand.Intn(1000000))        // Random ID

	certs = append(certs, cert)

	http.ServeFile(w, r, "static/success.html")
}

// Update cert by ID
func UpdateCert(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, item := range certs {
		if item.ID == params["id"] {
			certs = append(certs[:i], certs[i+1:]...) // delete old
			var updatedCert dto.Cert
			_ = json.NewDecoder(r.Body).Decode(&updatedCert)
			updatedCert.ID = params["id"]
			certs = append(certs, updatedCert)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCert)
			return
		}
	}
	http.Error(w, "Certificate not found", http.StatusNotFound)
}

// Delete cert by ID
func DeleteCert(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, item := range certs {
		if item.ID == params["id"] {
			certs = append(certs[:i], certs[i+1:]...)
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	http.Error(w, "Certificate not found", http.StatusNotFound)
}

func main() {
	r := mux.NewRouter()
	// context path

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/home.html")
	})
	// Serve static files (like CSS/JS/images)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Registered routes:")
	http.DefaultServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Requested:", r.Method, r.URL.Path)
	})

	api := r.PathPrefix("/worldcert").Subrouter()
	api.HandleFunc("/getallcerts", GetCerts).Methods("GET")
	api.HandleFunc("/getcertbyid/{id}", GetCert).Methods("GET")
	api.HandleFunc("/createcert", CreateCert).Methods("POST")
	api.HandleFunc("/updatecerts/{id}", UpdateCert).Methods("PUT")
	api.HandleFunc("/deletecertbyid/{id}", DeleteCert).Methods("DELETE")

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
