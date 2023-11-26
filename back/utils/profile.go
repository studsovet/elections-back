package utils

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSaveStudentData(c *gin.Context, token string) error {
	email, err := ExtractTokenEmail(c)
	// TODO check if saved
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.hseapp.ru/v2/dump/email/"+email, nil)
	req.Header.Set("Accept-Encoding", "br;q=1.0, gzip;q=0.9, deflate;q=0.8")
	req.Header.Set("User-Agent", "HSE App X/1.17.1 (ru.hse.HSEAppX; build:7379; iOS 15.3.1) Alamofire/5.5.0")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-ios-build", "7379")
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("Got incorrect answer") // А как понятнее написать?))
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Got incorrect answer") // А как понятнее написать?))
	}
	client_data := string(bodyBytes)
	print(client_data)
	_ = client_data // TODO save client data
	return nil
}
