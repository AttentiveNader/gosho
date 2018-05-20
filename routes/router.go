package routes

import (
	"../models"
	"../utils"
	"encoding/json"
	"github.com/gorilla/mux"
	hashids "github.com/speps/go-hashids"
	"log"
	"net/http"
	"time"
)

type DataDescending struct {
	r        *http.Request
	nextFunc func()
	url      Short
	w        http.ResponseWriter
	endExcu  bool
}

type Data struct {
	Err     bool
	ErrText string
}
type Short struct {
	Longurl  string `json:"longurl"`
	Id       string `json:"id"`
	Shorturl string `json:"shorturl"`
}

func InitRouter() *mux.Router {
	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("./static/"))
	r.HandleFunc("/", HomePage).Methods("GET")
	r.HandleFunc("/post", CreateUrlMang).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/r/{id}", RedirectToLongUrl).Methods("GET")
	return r
}
func RedirectToLongUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	longurl := models.Dbc.GetUrl(mux.Vars(r)["id"])
	if models.Dbc.Err != nil {
		log.Println(models.Dbc.Err, "RedirectTo")
		utils.ExecuteT(w, "index.gohtml", Data{Err: true, ErrText: models.Dbc.Err.Error()})
		return
	}
	http.Redirect(w, r, longurl, http.StatusSeeOther)
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	utils.ExecuteT(w, "index.gohtml", Data{})
}

func (data *DataDescending) DecodeBody() {
	decoder := json.NewDecoder(data.r.Body)
	var url Short
	err := decoder.Decode(&url)
	data.HandleErr(http.StatusInternalServerError, err, "Decondig Url err")
	data.url = url
	data.nextFunc = data.Checkexistence
}
func (data *DataDescending) Checkexistence() {
	shortMess := models.Dbc.GetShort(data.url.Longurl)
	if shortMess != "" {
		data.url.Shorturl = shortMess
		data.url.Id = shortMess[len(shortMess)-7:]
		data.nextFunc = data.EncodeJson
		return
	}
	data.HandleErr(http.StatusConflict, models.Dbc.Err, shortMess)
	data.nextFunc = data.Hashurl
}
func (data *DataDescending) Hashurl() {
	hd := hashids.NewData()
	hd.Salt = data.url.Longurl
	h, err := hashids.NewWithData(hd)
	data.HandleErr(http.StatusInternalServerError, err, "Hashing url Err")

	now := time.Now()
	id, err := h.Encode([]int{int(now.Unix())})
	data.url.Id = id
	data.url.Shorturl = "http://gosho.herokuapp.com/r/" + id
	data.HandleErr(http.StatusInternalServerError, err, "Hashing url err")
	data.nextFunc = data.InsertInto
}
func (data *DataDescending) InsertInto() {
	models.Dbc.InsertUrl(data.url.Id, data.url.Shorturl, data.url.Longurl)
	data.HandleErr(http.StatusInternalServerError, models.Dbc.Err, "sql Err")
	data.nextFunc = data.EncodeJson
}
func (data *DataDescending) EncodeJson() {
	data.w.WriteHeader(http.StatusOK)
	data.w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(data.w).Encode(data.url)
	data.HandleErr(http.StatusInternalServerError, err, "Encoding json err")
	data.endExcu = true
}
func (data *DataDescending) HandleErr(status int, err error, message string) {
	if err != nil {
		data.w.WriteHeader(status)
		data.w.Write([]byte(message + " " + err.Error()))
		data.endExcu = true
	}
}

func CreateUrlMang(w http.ResponseWriter, r *http.Request) {
	data := &DataDescending{
		w:       w,
		r:       r,
		endExcu: false,
		url:     Short{},
	}
	data.nextFunc = data.DecodeBody
	data.CreateUrlFun()
}
func (data *DataDescending) CreateUrlFun() { //thought it is better to use recursion to handle errors also I can skip function when I need
	data.nextFunc()
	if data.endExcu {
		return
	}
	data.CreateUrlFun()
}
