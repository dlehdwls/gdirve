package main

import (
	"fmt"
	"gdrive/drive"
	"log"
)

func main() {
	fname := ""
	fmt.Print("문서 이름을 입력해 주세요 :")
	_, err := fmt.Scan(&fname)
	if err != nil {
		log.Fatal(err)
	}

	err = drive.CopyDocs(fname)
	if err != nil {
		log.Fatal(err)
	}
}
