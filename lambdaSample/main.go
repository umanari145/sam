package main

import (
	"context"
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
	//ID はプライマリーキー
	ID string
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

	fmt.Printf("DB_CONNECT:%s", CONNECT)

	db, err := gorm.Open(DBMS, CONNECT)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	log.Print("DB connect")

	return db, nil
}

//MyResponse はレスポンス
type MyResponse struct {
	//HTTPStatusCode はHTTPステータスコード
	HTTPStatusCode int
	//Body は任意のレスポンス
	Body interface{}
}

//urlのパターンは　/area/274-0077(郵便番号)
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (MyResponse, error) {
	zipCode := request.PathParameters["zipCode"]
	area, error := loadAreaFromZip(zipCode)

	var response MyResponse
	if error != nil {
		response.HTTPStatusCode = 400
		response.Body = "郵便番号が存在しません。"
		return response, error
	}

	response.Body = area
	return response, nil
}

func main() {
	lambda.Start(handler)
}
