package service

import (
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/nooderg/autotest_templating/pkg/parsing"
)

type GenerateRequest struct {
	UserId string `json:"user_id"`
	File   string `json:"file"`
}

func Generate(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest

	err := yaml.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "cannot decode body", 400)
		return
	}

	file, err := b64.StdEncoding.DecodeString(req.File)
	if err != nil {
		http.Error(w, "cannot decode body from b64", 400)
	}

	jsonRes, err := parsing.ParseOpenApi(file)
	if err != nil {
		http.Error(w, "could not parse yaml file", 500)
		return
	}

	json.NewEncoder(w).Encode(jsonRes)
}

func UploadFileAndSyncDB(file *os.File, userid string) error {
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region:      aws.String("fr-par"),
	// 	Endpoint:    aws.String("https://autotests3.s3.fr-par.scw.cloud"),
	// 	Credentials: credentials.NewSharedCredentials("./.aws/credentials", "default"),
	// }))

	// uploader := s3manager.NewUploader(sess)

	// result, err := uploader.Upload(&s3manager.UploadInput{
	// 	Bucket: aws.String("https://autotests3.s3.fr-par.scw.cloud"),
	// 	Key:    aws.String("file.yaml"),
	// 	Body:   file,
	// })
	// if err != nil {
	// 	log.Println("failed to upload file, ", err)
	// }
	// loc := aws.StringValue(&result.Location)

	// log.Println(result.Location)

	// db := config.GetDBClient()

	// err := db.Model(&domain.User{}).Where("id = ?", userid).Update("file_url", "https://www.youtube.com/watch?v=dQw4w9WgXcQ").Error
	// if err != nil {
	// 	return err
	// }

	return nil
}
