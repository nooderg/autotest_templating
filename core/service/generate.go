package service

import (
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/nooderg/autotest_templating/config"
	"github.com/nooderg/autotest_templating/core/domain"
	"github.com/nooderg/autotest_templating/pkg/parsing"
	"github.com/nooderg/autotest_templating/pkg/templating"
)

type GenerateRequest struct {
	UserId string `json:"user_id"`
	File   string `json:"file"`
}

func Generate(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "cannot decode body", 400)
		return
	}

	file, err := b64.StdEncoding.DecodeString(req.File)
	if err != nil {
		http.Error(w, "cannot decode body from b64", 400)
	}

	dataArr, err := parsing.ParseOpenApi(string(file))
	if err != nil {
		http.Error(w, "could not parse yaml file", 500)
		return
	}
	bfile, _ := templating.TemplateFile(dataArr)

	localfile, _ := os.CreateTemp(".", "autotest_tavern_template")
	defer os.Remove(localfile.Name())
	localfile.Write(bfile.Bytes())

	go UploadFileAndSyncDB(localfile, req.UserId)

	return
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

	db := config.GetDBClient()

	err := db.Model(&domain.User{}).Where("id = ?", userid).Update("file_url", "https://www.youtube.com/watch?v=dQw4w9WgXcQ").Error
	if err != nil {
		return err
	}

	return nil
}
