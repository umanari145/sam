package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

//Area は地域
type Area struct {
	//zip は郵便番号
	Zip string
	//Pref は県
	Pref string
	//City は市
	City string
	//Town は町村
	Town string
	//PrefKana は県のカナ
	PrefKana string
	//CityKana は市のカナ
	CityKana string
	//TownKana は町村のカナ
	TownKana string
}

//loadAreaFromZip は地域を郵便番号をみる
func loadAreaFromZip(zipCode string) (Area, error) {
	dbConn, err := Connect()

	if err != nil {
		return Area{}, fmt.Errorf("loadAreaFromZip Cannnot Connect DB %s", err.Error())
	}
	defer dbConn.Close()

	area := Area{}

	if zipCode == "" {
		return Area{}, fmt.Errorf("loadAreaFromZip ZipCode is not exist")
	}

	dbConn.Table("area").Where("zip = ?", zipCode).First(&area)

	if area.Zip == "" {
		return Area{}, fmt.Errorf("地域が存在しません。")
	}

	return area, nil
}

//Connect はDBへの接続
func Connect() (*gorm.DB, error) {
	DBMS := os.Getenv("DB_TYPE")
	USER := os.Getenv("POSTGRES_DBUSER")
	PASS := os.Getenv("POSTGRES_PASSWORD")
	DBNAME := os.Getenv("POSTGRES_DBNAME")

	CONNECT := "host=" + os.Getenv("POSTGRES_DBHOST") +
		" port=5432" +
		" user=" + USER +
		" dbname=" + DBNAME +
		" password=" + PASS +
		" sslmode=disable"

	log.Printf("DB_CONNECT:%s", CONNECT)

	db, err := gorm.Open(DBMS, CONNECT)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	log.Print("DB connect")

	return db, nil
}

//MyResponse はレスポンス
type MyResponse events.APIGatewayProxyResponse

//urlのパターンは　/area/274-0077(郵便番号)
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (MyResponse, error) {
	log.Print("API start")
	zipCode := request.PathParameters["zipCode"]
	log.Printf("zipCode %s", zipCode)

	area, error := loadAreaFromZip(zipCode)

	var response MyResponse
	response.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	if error != nil {
		log.Printf("Invalid argument error occur %s", error)
		response.StatusCode = 400
		return response, error
	}
	var buf bytes.Buffer

	body, error := json.Marshal(area)

	if error != nil {
		log.Printf("JSON parse error %s", error)
		response.StatusCode = 500
		return response, error
	}

	json.HTMLEscape(&buf, body)

	response.StatusCode = 200
	response.Body = buf.String()

	return response, nil
}

func main() {
	lambda.Start(handler)
}
