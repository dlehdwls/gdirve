package docs

import (
	"context"
	"fmt"
	"gdrive/drive/auth"
	"time"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"
)

const TOKEN_NAME = "docs-token.json"

func GetLeaveApplicationKeyMap(now time.Time, name, department, position, day string) map[string]string {
	return map[string]string{
		"{{name}}":          name,
		"{{department}}":    department,
		"{{position}}":      position,
		"{{period_of_use}}": fmt.Sprintf("%04d.  %02d.  %02d.   ( %s )요일   ( %s 일간)", now.Year(), now.Month(), now.Day(), now.Weekday().String(), day),
		"{{registe_date}}":  fmt.Sprintf("%04d년   %02d월   %02d일 ", now.Year(), now.Month(), now.Day()),
	}
}

// keyMap: "{{{키워드}}}" → "바꿀문구"
func ReplaceKeywordTexts(id string, keyMap map[string]string) error {
	config := auth.GetGoogleConfig(docs.DocumentsScope)
	token, err := auth.TokenFromFile(TOKEN_NAME)
	if err != nil {
		token = auth.GetTokenFromWeb(config)
		auth.SaveToken(TOKEN_NAME, token)
	}

	client := config.Client(context.Background(), token)
	serv, err := docs.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return err
	}

	var requests []*docs.Request
	for placeholder, repl := range keyMap {
		requests = append(requests, &docs.Request{
			ReplaceAllText: &docs.ReplaceAllTextRequest{
				ContainsText: &docs.SubstringMatchCriteria{
					Text:      placeholder,
					MatchCase: true,
				},
				ReplaceText: repl,
			},
		})
	}

	_, err = serv.Documents.BatchUpdate(id, &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}).Do()
	return err
}
