package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/gorilla/mux"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

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

	var authorized = false

	var apikeys, ok = r.URL.Query()["apikey"]

	if !ok || len(apikeys[0]) < 1 {
		fmt.Println("Url Param 'apikey' is missing")
		var api_key_header = r.Header.Get("X-API-Key")
		authorized = api_key_header == "123456"

	} else {
		authorized = apikeys[0] == "123456"
	}

	if !authorized {
		fmt.Println("Not authorized!")
		w.WriteHeader(401)
		return
	}

	file, filename, err := getFilefromRequest(r)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		w.WriteHeader(400)
		return
	}
	defer file.Close()

	file_uploaded := upload_to_rrf(filename, file)

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
		"done": file_uploaded,
	}

	json.NewEncoder(w).Encode(result)
}

func upload_to_rrf(filename string, file multipart.File) bool {
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://192.168.178.69/rr_upload?name=gcodes/%v", filename), file)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Uploaded File: %+v\n", content)

	var file_uploaded = true
	return file_uploaded
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/files/local", fileUpload).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

type myAppClass struct {
	app    fyne.App
	window fyne.Window
}

func createMyApp(title string) *myAppClass {
	fApp := app.New()
	//a.SetIcon(icon.Data)
	return &myAppClass{fApp, fApp.NewWindow(title)}
}

func (myApp *myAppClass) onWindowCloseHandler() func() {
	return func() {
		myApp.window.Hide()
	}
}

func (myApp *myAppClass) onSettingsClicked() {
	myApp.window.Show()
}

func main() {
	fmt.Println("OctoPrint to RepRap REST API Wrapper listening on '127.0.0.1:10000'")

	myApp := createMyApp("OctoPrint to RepRap REST API Wrapper")

	a := myApp.app

	w := myApp.window
	octoprint_api := widget.NewMultiLineEntry()
	octoprint_api.SetText("123456")
	reprap_hostname := widget.NewMultiLineEntry()

	w.SetContent(container.NewVBox(
		container.NewHBox(widget.NewLabel("Octoprint API"), octoprint_api),
		container.NewHBox(widget.NewLabel("RepRap Hostname"), reprap_hostname),
		widget.NewButton("Close", myApp.onWindowCloseHandler()),
	))

	w.Resize(fyne.NewSize(480, 360))

	w.SetCloseIntercept(myApp.onWindowCloseHandler())

	systray.Register(func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("OctoPrint to RepRap REST API Wrapper")
		systray.SetTooltip("OctoPrint to RepRap REST API Wrapper")
		about_menu := systray.AddMenuItem("Settings", "Settings")
		menu_quit := systray.AddMenuItem("Quit", "Quit the app")
		go func() {
			for {
				select {
				case <-about_menu.ClickedCh:
					myApp.onSettingsClicked()

				case <-menu_quit.ClickedCh:
					systray.Quit()
					os.Exit(0)
				}
			}
		}()
	}, nil)

	go func() {
		handleRequests()
	}()

	a.Run()
}
