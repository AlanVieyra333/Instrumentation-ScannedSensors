package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

var validPath = regexp.MustCompile("^/(json)*(/([a-zA-Z0-9]*))?$")
var validPathImage = regexp.MustCompile("^(https?://)?" + "([a-zA-Z0-9]*.?)+" + "(:[0-9]*)?" + "/uploaded/([a-zA-Z0-9]*.[a-zA-Z0-9]*)" + "$")

/*	---------------------------*	Functions	---------------------------*	*/

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

/*	--------------------------->	Handlers	--------------------------->	*/

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	indexTmpl := template.New("index.html").Delims("<<", ">>")
	indexTmpl, _ = indexTmpl.ParseFiles("index.html")
	indexTmpl.Execute(w, nil)
}

/*	--------------------------->	JSON	--------------------------->	*/

type ReqJSON struct {
	Operation int
}

// Temperatura C°, Sensor de luz con activa motor {0,1}, ritmo cardiaco BMP.
type ResJSON struct {
	Code    int
	Message string
	Degrees int
	Motor   int
	BPM     int
}

func jsonHandler(w http.ResponseWriter, req *http.Request) {
	//fmt.Printf("recive\n")
	confirm := &ResJSON{
		Code:    1,
		Message: "Datos obtenidos correctamente.",
		Degrees: 0,
		Motor:   0,
		BPM:     0,
	}
	var data ReqJSON

	if req.Body == nil {
		confirm.Code = 2
		confirm.Message = "Please send a request body."
	} else {
		bodyb, err := ioutil.ReadAll(req.Body)
		if err != nil {
			confirm.Code = 2
			confirm.Message = "Please send a request body."
		} else {
			//fmt.Println(string(bodyb))
			err = json.Unmarshal([]byte(string(bodyb)), &data)
			if err != nil {
				confirm.Code = 2
				confirm.Message = err.Error()
			} else if data.Operation == 1 { // Get data from micro.
				degrees, motor, bpm, err := GetData()

				if err != nil {
					confirm.Code = 2
					confirm.Message = err.Error()
				} else {
					confirm.Degrees = degrees
					confirm.Motor = motor
					confirm.BPM = bpm
				}
			}
		}
	}

	// Parse to JSON format.
	b, err := json.Marshal(confirm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

/*	--------------------------->	Server turn on	--------------------------->	*/
func ServerTurnOn() {
	http.HandleFunc("/", makeHandler(homeHandler))
	http.HandleFunc("/json/", makeHandler(jsonHandler))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css/"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./fonts/"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js/"))))

	port := flag.Int("port", 8080, "port to serve on")
	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	log.Printf("Servidor listo en: %s\n", addr)

	err := http.ListenAndServe(addr, nil)
	fmt.Println(err.Error())
}
