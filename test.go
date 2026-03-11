package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/wkrdev/gopkg/crypto"
	"github.com/wkrdev/gopkg/dbms"
	"github.com/wkrdev/gopkg/restapi"
)

type JwtClaims struct {
	ISS string `json:"iss"`
	IAT int64  `json:"iat"`
	EXP int64  `json:"exp"`
	UID uint32 `json:"uid"`
	ORG uint16 `json:"org"`
}

var secret string = "123456"

func init() {
	envFile := flag.String("e", ".env", "env file")
	flag.Parse()
	err := godotenv.Load(*envFile)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbmsTest()
	cryptoTest()
	// restapi
	fmt.Println("- restapi -")
	res, err := restapi.Call(&restapi.ApiPack{
		Url:    "https://wkrh.info/api/v1",
		Method: "GET",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(res))
}

func dbmsTest() {
	db, err := dbms.Init(&dbms.Config{
		RWDB: &dbms.DBCON{
			HOST: os.Getenv("RWDB_HOST"),
			USER: os.Getenv("RWDB_USER"),
			PASS: os.Getenv("RWDB_PASS"),
		},
		RODB: &dbms.DBCON{
			HOST: os.Getenv("RODB_HOST"),
			USER: os.Getenv("RODB_USER"),
			PASS: os.Getenv("RODB_PASS"),
		},
		REDIS: &dbms.RDBCON{
			DB:   1,
			HOST: os.Getenv("REDIS_HOST"),
			PASS: os.Getenv("REDIS_PASS"),
		},
	})
	if err != nil {
		fmt.Println("dbms init", err)
		return
	}
	// dbms mysql
	fmt.Println("- dbms mysql -")
	rows, err := db.RO.Query("SHOW DATABASES")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			fmt.Println(err)
		}
		fmt.Println(dbName)
	}
	// dbms cache
	fmt.Println("- dbms cache -")
	var inCache, outCache JwtClaims
	inCache.ISS = "cache"
	inCache.IAT = 3600
	// set
	err = db.SetCache("test", &inCache, time.Minute*1)
	if err != nil {
		fmt.Println(err)
		return
	}
	// get
	err = db.GetCache("test", &outCache)
	if err != nil {
		fmt.Println(err)
		return
	}
	// del
	db.DelCache("test")
	// get keys
	keys, err := db.GetCacheKeys()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("get: %+v\n", outCache)
	fmt.Printf("key: %+v\n", keys)
}

func cryptoTest() {
	fmt.Println("- crypto -")
	claims := &JwtClaims{
		ISS: "iamsvc",
		IAT: time.Now().Unix(),
		EXP: time.Now().Add(time.Hour).Unix(),
		UID: 55,
		ORG: 88,
	}
	token, err := crypto.JwtEncode(claims, secret)
	if err != nil {
		fmt.Println(err)
		return
	}
	jwt, err := crypto.JwtDecode[JwtClaims](token, secret)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", jwt)
	fmt.Printf("%+v\n", crypto.JwtParse[JwtClaims](token))
	fmt.Println("STR:", crypto.StringHash(32))
	fmt.Println("MD5:", crypto.MD5Hash("123456"))
	fmt.Println("SHA:", crypto.SHA256Hash("123456"))
}
