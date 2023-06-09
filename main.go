package main

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"sync"
	"strconv"
	"log"
	"time"
	"fmt"
	"io"
	
	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
)

type Promotion struct {
	ID             string   
	Price          float64  
	ExpirationDate time.Time
}

type PromotionMutex struct {
    sync.RWMutex
    promotionsMap map[string]Promotion
}
var promotions = PromotionMutex {
	promotionsMap: make(map[string]Promotion),
}

type PromotionResponse struct {
    ID             string  `json:"id"`
    Price          string `json:"price"`
    ExpirationDate string  `json:"expiration_date"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	csvFilePath := os.Getenv("CSV_FILE_PATH")
	sleepTime := os.Getenv("SLEEP_TIME")
	serverPort := os.Getenv("SERVER_PORT")
	sleepDuration, err := time.ParseDuration(sleepTime)

	go func() {
		for {
			tempPromotionsMap, err := updateCsvData(csvFilePath)
			if err != nil {
				log.Fatal(err)
			}

			promotions.Lock()
			promotions.promotionsMap = tempPromotionsMap
			promotions.Unlock()

			time.Sleep(sleepDuration)
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/promotions/{id}", getPromotion).Methods("GET")
	http.ListenAndServe(":"+serverPort, r)
}

func getPromotion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	promotions.RLock()
	defer promotions.RUnlock()
	promotion, ok := promotions.promotionsMap[id]
	if ok {
		response := PromotionResponse{
			ID: promotion.ID,
			Price: fmt.Sprintf("%.2f", promotion.Price),
			ExpirationDate: promotion.ExpirationDate.Format("2006-01-02 15:04:05"),
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func updateCsvData(filePath string) (map[string]Promotion, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	tempPromotionsMap := make(map[string]Promotion)

	const layout = "2006-01-02 15:04:05 -0700 MST"
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		price, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, err
		}
		expirationDate, err := time.Parse(layout, line[2])
		if err != nil {
			return nil, err
		}
		tempPromotionsMap[line[0]] = Promotion{
			ID:             line[0],
			Price:          price,
			ExpirationDate: expirationDate,
		}
	}
	return tempPromotionsMap, nil
}
