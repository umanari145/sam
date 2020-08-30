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
func loadAreaFromZip(zipCode string) Area {
	dbConn, err := Connect()

	if err != nil {
		log.Printf("loadAreaFromZip Cannnot Connect DB %s", err.Error())
		return Area{}
	}
	defer dbConn.Close()

	area := Area{}

	if zipCode == "" {
		log.Println("loadAreaFromZip ZipCode is not exist")
		return Area{}
	}

	dbConn.Table("area").Where("zip = ?", zipCode).First(&area)

	return area

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

//urlのパターンは　/area/274-0077(郵便番号)
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (string, error) {
	zipCode := request.PathParameters["zipCode"]
	loadAreaFromZip(zipCode)

	return "", nil
}

func main() {
	lambda.Start(handler)
}
