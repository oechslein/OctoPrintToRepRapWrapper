package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/mux"
)

func getFilefromRequest(r *http.Request) (multipart.File, string, error) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, "", err
	}

	return file, handler.Filename, nil
}

func fileUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: fileUpload")

	var rrf_server = r.Header.Get("X-API-Key")

	file, filename, err := getFilefromRequest(r)
	if err != nil {
		fmt.Println("Error retrieving the file")
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading the file")
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	err = upload_to_rrf(filename, bytes.NewReader(fileContent), rrf_server)
	if err != nil {
		fmt.Println("Error uploading the file")
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}

	var result = map[string]interface{}{
		"files": map[string]interface{}{
			"local": map[string]interface{}{
				"name":   filename,
				"origin": "local",
				"refs": map[string]interface{}{
					"resource": fmt.Sprintf("http://%v/api/files/local/%v", r.Host, filename),
					"download": fmt.Sprintf("http://%v/downloads/files/local/%v", r.Host, filename),
				},
			},
		},
		"done": true,
	}

	json.NewEncoder(w).Encode(result)
}

func upload_to_rrf(filename string, data *bytes.Reader, rrf_server string) error {
	client := &http.Client{}
	url := fmt.Sprintf("http://%v/rr_upload?name=gcodes/%v", rrf_server, filename)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("Uploaded File: %+v\n", content)

	return nil
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/files/local", fileUpload).Methods("POST")
	log.Fatal(http.ListenAndServe(":80", myRouter))
}

func main() {
	fmt.Println("OctoPrint to RepRap REST API Wrapper listening on '127.0.0.1'")
	handleRequests()
}
