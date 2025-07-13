package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func TokenFromFile(name string) (*oauth2.Token, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func GetGoogleConfig(scope ...string) *oauth2.Config {
	// 1. 인증 정보 읽기
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
		return nil
	}

	// 2. Drive API용 OAuth2 설정
	config, err := google.ConfigFromJSON(b, scope...)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil
	}

	return config
}

func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("다음 링크를 열어 인증 코드를 입력하세요: \n%v\n", authURL)

	var authCode string
	fmt.Print("인증 코드 입력: ")
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("인증 코드 입력 실패: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("토큰 교환 실패: %v", err)
	}
	return tok
}

func SaveToken(name string, token *oauth2.Token) {
	fmt.Printf("토큰 저장: %s\n", name)
	f, err := os.Create(name)
	if err != nil {
		log.Fatalf("토큰 저장 실패: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
