package file

import (
	"bytes"
	"chat/handlers/ws"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var Client *s3.Client

type UrlFile struct {
	Url      string
	FileName string
}

func SetupClient() error {
	c, err := MakeAWSClient()
	if err != nil {
		return err
	}
	Client = c
	return nil
}
func FileUpload(w http.ResponseWriter, r *http.Request, hub *ws.Hub) {
	log.Println("got a file")
	r.ParseMultipartForm(100 << 20)
	file, fileHeader, err := r.FormFile("upload")
	userName := r.FormValue("userName")
	chatClient, ok := hub.ClientMap[userName]

	if !ok {
		log.Printf("user name %s not found in Hub ", userName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("socket connection dead"))
		return

	}

	if err != nil {
		log.Println("error occured while getting form file in file upload", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while access form file"))

		return
	}

	fileName := fileHeader.Filename

	data, err := io.ReadAll(file)

	if err != nil {
		log.Println("error occured while reading form file", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while reading form file"))
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Println("error occured while getting wd path", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while reading form file"))
		return
	}

	reader := bytes.NewReader(data)

	tempFile, err := os.CreateTemp(dir, fileName)
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
		log.Println("temp file removed ")
	}()
	log.Println("file anme is ", fileName, dir, tempFile.Name())

	if err != nil {
		log.Println("error ocuured while creating temp file", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while reading form file"))
		return
	}

	io.Copy(tempFile, reader)
	tempFile.Close()
	tempFile, err = os.Open(tempFile.Name())
	if err != nil {
		log.Println("error ocuured while opening temp file", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while reading form file"))
		return
	}
	err = UploadFile(Client, tempFile, fileName)

	if err != nil {
		log.Println("error occured while uplaoding file ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while uploading form file"))
		return
	}

	url, err := GetFileDownlaodUrl(Client, fileName)

	if err != nil {
		log.Println("error occured while getting file url", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while uploading form file"))
		return
	}
	// urlFile := types.UrlFile{
	// 	Url:url,
	// 	FileName: fileName,
	// }
	log.Println("printting ", chatClient)
	msg := ws.SendMessage{Sender: chatClient.Name, Message: url}

	bm, err := json.Marshal(msg)
	if err != nil {
		log.Println("erro occured why conv send file data to byte from json", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error occured while uploading form file"))
		return

	}
	chatClient.Hub.Broadcast <- bm
	log.Println("file uplaod suuccess ", err)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("file uploaded success"))

}
