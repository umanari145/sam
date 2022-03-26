package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestValidCheck(t *testing.T) {
	//エラー
	err := validCheck("")
	if err == nil {
		t.Fatalf("空白")
	} else {
		fmt.Printf("空白 %s \n", err)
	}

	//エラー
	err2 := validCheck("aaaa")
	if err2 == nil {
		t.Fatalf("7桁NG")
	} else {
		fmt.Printf("7桁NG %s \n", err)
	}

	//正常系
	err3 := validCheck("2740077")
	if err3 != nil {
		t.Fatalf("7桁 %s", err)
	} else {
		fmt.Printf("7桁OK \n")
	}
}

func TestDbConnect(t *testing.T) {
	dbHost := os.Getenv("POSTGRES_DBHOST")

	//エラー
	os.Setenv("POSTGRES_DBHOST", "aaa")
	_, err := Connect()
	if err != nil {
		fmt.Printf("DB接続NG %s\n", err)
	} else {
		t.Fatalf("DB接続NG")
	}

	// 正常系
	os.Setenv("POSTGRES_DBHOST", dbHost)
	_, err2 := Connect()
	if err2 == nil {
		fmt.Printf("DB接続OK \n")
	} else {
		t.Fatalf("DB接続OK")
	}
}

func TestLoadAreaFromZip(t *testing.T) {

	dbConn, err := Connect()
	if err != nil {
		t.Fatalf("DB接続NG")
	}

	// 正常系
	area, err2 := loadAreaFromZip(dbConn, "2740077")
	if err2 != nil {
		t.Fatalf("zip取得 %s", err2)
	} else {
		fmt.Printf("zip取得 %s\n", area)
	}

}

//わざとSQLをこわして実行する。
/*func TestLoadAreaFromZipError(t *testing.T) {

	dbConn, err := Connect()
	if err != nil {
		t.Fatalf("DB接続NG")
	}

	// 正常系
	_, err2 := loadAreaFromZip(dbConn, "2740077")
	if err2 != nil {
		fmt.Printf("zip取得NG %s\n", err2)
	} else {
		t.Fatalf("zip取得NG")
	}
}*/

func TestParseJson(t *testing.T) {

	dbConn, err := Connect()
	if err != nil {
		t.Fatalf("DB接続NG")
	}

	// 正常系
	area, err2 := loadAreaFromZip(dbConn, "274007997")
	if err2 != nil {
		t.Fatalf("zip取得 %s", err2)
	} else {
		fmt.Printf("zip取得 %s\n", area)
	}

	var buf bytes.Buffer
	body, err := json.Marshal(area)
	if err != nil {
		t.Fatalf("JSONパース失敗")
	}

	json.HTMLEscape(&buf, body)
	fmt.Println(buf.String())
}
