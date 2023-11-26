package utils

import (
	"elections-back/db"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Education struct {
	FacultytTitle string `json:"faculty_title"`
}
type Response struct {
	FullName  string      `json:"full_id"`
	Email     string      `json:"email"`
	Education []Education `json:"education"`
}

func createElector(c *gin.Context, resp Response) (db.Elector, error) {
	var elector db.Elector
	id, err := ExtractTokenID(c)
	if err != nil {
		return db.Elector{}, err
	}
	elector.ID = id
	var FacultyIds []string = []string{}
	for _, education := range resp.Education {
		FacultyIds = append(FacultyIds, education.FacultytTitle)
	}
	elector.FullName = resp.FullName
	elector.Email = resp.Email
	elector.FacultyIds = FacultyIds
	return elector, nil
}

func GetSaveElectorData(c *gin.Context, token string) error {
	email, err := ExtractTokenEmail(c)
	if err != nil {
		return err
	}
	id, err := ExtractTokenID(c)
	if err != nil {
		return err
	}
	saved, err := db.IsElectorSaved(id)
	if err != nil {
		return err
	}
	if saved {
		return nil
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
		return errors.New("Got incorrect answer")
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("Got incorrect answer")
	}
	var response Response
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return errors.New("Got incorrect answer")
	}
	elector, err := createElector(c, response)
	if err != nil {
		return errors.New("Got incorrect answer")
	}
	_, err = elector.Save()
	if err != nil {
		return err
	}
	return nil
}
