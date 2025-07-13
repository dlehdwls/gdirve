package drive

import (
	"context"
	"errors"
	"fmt"
	"gdrive/drive/auth"
	"gdrive/drive/docs"
	"log"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

const TOKEN_NAME = "drive-token.json"

func NewService(ctx context.Context) (*drive.Service, error) {
	config := auth.GetGoogleConfig(drive.DriveScope)
	if config == nil {
		return nil, errors.New("Unable to parse client secret file to config")
	}

	tok, err := auth.TokenFromFile(TOKEN_NAME)
	if err != nil {
		tok = auth.GetTokenFromWeb(config)
		auth.SaveToken(TOKEN_NAME, tok)
	}

	client := config.Client(context.Background(), tok)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil
}

func SearchFileRecursivelyFromRoot(srv *drive.Service, targetName string) (*drive.File, error) {
	return SearchFileRecursively(srv, "root", targetName)
}

func SearchFileRecursively(srv *drive.Service, parentID string, targetName string) (*drive.File, error) {
	query := fmt.Sprintf("'%s' in parents", parentID)
	r, err := srv.Files.List().
		Q(query).
		Fields("files(id, name, mimeType, parents)").
		Do()
	if err != nil {
		return nil, err
	}

	for _, f := range r.Files {
		if f.MimeType == "application/vnd.google-apps.folder" {
			found, err := SearchFileRecursively(srv, f.Id, targetName)
			if err != nil {
				return nil, err
			}
			if found != nil {
				return found, nil
			}
		} else {
			if f.Name == targetName {
				return f, nil
			}
		}
	}

	return nil, errors.New("file not found") // 파일 없음
}

func CopyDocs(fname string) error {
	ctx := context.TODO()
	srv, err := NewService(ctx)
	if err != nil {
		return err
	}

	f, err := SearchFileRecursivelyFromRoot(srv, fname)
	if err != nil {
		return err
	}

	now := time.Now()
	fname = fmt.Sprintf("%s-%s", now.Format(time.DateOnly), fname)
	to := &drive.File{
		Name: fname,
	}

	copyf, err := srv.Files.Copy(f.Id, to).Do()
	if err != nil {
		return err
	}

	err = docs.ReplaceKeywordTexts(copyf.Id, docs.GetLeaveApplicationKeyMap(now, "홍길동", "개발팀", "개발자", "월요일"))
	if err != nil {
		log.Fatalf("fail replace %v", err)
	}

	fmt.Printf("✅ 복사된 문서 ID: %s\n", copyf.Id)
	fmt.Printf("🔗 링크: https://docs.google.com/document/d/%s\n", copyf.Id)
	return nil
}
