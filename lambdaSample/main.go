package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

//Area は地域
type Area struct {
	//zip は郵便番号
	Zip string `json:"zip"`
	//Pref は県
	Pref string `json:"pref"`
	//City は市
	City string `json:"city"`
	//Town は町村
	Town string `json:"town"`
	//PrefKana は県のカナ
	PrefKana string `json:"pref_kana"`
	//CityKana は市のカナ
	CityKana string `json:"city_kana"`
	//TownKana は町村のカナ
	TownKana string `json:"town_kana"`
}

func validCheck(zipCode string) error {

	if zipCode == "" {
		return fmt.Errorf("郵便番号が未入力です。")
	}

	regex := regexp.MustCompile(`^[0-9]{7}$`)
	if !regex.MatchString(zipCode) {
		return fmt.Errorf("半角数字7桁で入力してください。")
	}

	return nil
}

//loadAreaFromZip は地域を郵便番号をみる
func loadAreaFromZip(dbConn *gorm.DB, zipCode string) (Area, error) {

	area := Area{}
	if err := dbConn.Table("area").Where("zip = ?", zipCode).First(&area).Error; err != nil {
		return Area{}, fmt.Errorf("SQLの実行に失敗しました。 %w", err)
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
	log.Print("DBに接続しました。")

	return db, nil
}

//MyResponse はレスポンス
type MyResponse events.APIGatewayProxyResponse

//urlのパターンは　/area/274-0077(郵便番号)
func handler(ctx context.Context, request events.APIGatewayProxyRequest) MyResponse {

	zipCode := request.PathParameters["zipCode"]
	log.Printf("住所検索取得APIを開始します。リクエストパラメーター %s", zipCode)

	err := validCheck(zipCode)

	var response MyResponse
	response.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	if err != nil {
		errorMsg := fmt.Sprintf("入力値に不適切です。メッセージ　%s", err)
		log.Printf(errorMsg)
		response.StatusCode = 400
		response.Body = errorMsg
		return response
	}

	dbConn, err := Connect()
	if err != nil {
		errorMsg := fmt.Sprintf("DB接続時に失敗しました。メッセージ　%s", err)
		log.Printf(errorMsg)
		response.StatusCode = 500
		response.Body = errorMsg
		return response
	}
	defer dbConn.Close()

	area, err := loadAreaFromZip(dbConn, zipCode)
	if err != nil {
		errorMsg := fmt.Sprintf("郵便番号の取得に失敗しました。%s", err)
		log.Printf(errorMsg)
		response.StatusCode = 500
		response.Body = errorMsg
		return response
	}

	if area.Zip == "" {
		msg := fmt.Sprintf("住所が存在しません。")
		log.Printf(msg)
	}

	var buf bytes.Buffer
	body, err := json.Marshal(area)

	if err != nil {
		errorMsg := fmt.Sprintf("JSONのパースに失敗しました。%s", err)
		log.Printf(errorMsg)
		response.StatusCode = 500
		response.Body = errorMsg
		return response
	}

	json.HTMLEscape(&buf, body)

	response.StatusCode = 200
	response.Body = buf.String()

	return response
}

func main() {
	lambda.Start(handler)
}
