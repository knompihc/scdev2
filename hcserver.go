package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/revel/revel"
)

var reqChan, respChan chan string

type Address struct {
	IP   string
	Port string
	//Uri Uri
}

//Contains necessary information to communicate with a light network
type NetAdaptor struct {
	Id       uint64
	Address  Address
	Protocol string
	Client   string
}

type Resp struct {
	Data interface{}
}
type ReqData struct {
	FirstByte  string
	SecondByte string
	DeviceId   string
}

const (
	PORT = "9001"
)

var (
	tclNetAdaptor *NetAdaptor
)

func init() {
	//lorawan specific change
	tclNetAdaptor = &NetAdaptor{
		Id:       2538,
		Address:  Address{"http://168.87.87.213:8080/davc/m2m/HPE_IoT/0004a30b001ba065/DownlinkPayload", "8080"},
		Protocol: "http",
		Client:   "tcl",
	}
}

var templates = template.Must(template.ParseFiles("hexcode.html", "response.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, resp interface{}) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	reqChan, respChan = make(chan string), make(chan string)

	http.HandleFunc("/", viewHex)
	http.HandleFunc("/resp", respHandler)
	fmt.Println("Started running at", PORT)
	http.ListenAndServe(":"+PORT, nil)
}

func viewHex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "hexcode", "")
}

func respHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	chkErr(err)

	firstByteStr := r.Form.Get("firstByte")
	secondByteStr := r.Form.Get("secondByte")
	deviceIdStr := r.Form.Get("deviceId")

	reqData := NewReqData(firstByteStr, secondByteStr, deviceIdStr)
	resp, err := reqData.makeResponse()
	if err != nil {
		log.Println("no response received")
		resp = []byte("no response received")
	}
	var response interface{}
	json.Unmarshal(resp, &response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
	//renderTemplate(w, "response", Resp{string(resp)})
}

func NewReqData(firstByte, secondByte, deviceId string) *ReqData {
	return &ReqData{
		FirstByte:  firstByte,
		SecondByte: secondByte,
		DeviceId:   deviceId,
	}
}

func (rd *ReqData) makeResponse() ([]byte, error) {
	log.Println("input received ", rd.FirstByte, rd.SecondByte, rd.DeviceId)
	data, err := rd.makeData()
	if err != nil {
		log.Println("error in creating data")
		return nil, err
	}
	resp, err := sendData(data)
	if err != nil {
		log.Println("error in sending data")
		return nil, err
	}
	return resp, nil
}

func (rd *ReqData) makeData() ([]byte, error) {
	firstHex, err := hex.DecodeString(rd.FirstByte)
	if err != nil {
		log.Println(err, "convert error")
		return nil, err
	}
	secondHex, err := hex.DecodeString(rd.SecondByte)
	if err != nil {
		log.Println(err, "convert error")
		return nil, err
	}

	log.Printf("hex code received : first byte: %x, second byte: %x\n",
		firstHex[0], secondHex[0])
	plData := base64.StdEncoding.EncodeToString([]byte{firstHex[0], secondHex[0]})
	log.Println("plData ", plData)
	data := `{
		"m2m:cin": {
			"ty":4,
			"cnf":"text/plain:0",
			"cs":300,
			"con":"{\"payload_dl\":{\"deveui\":\"` + rd.DeviceId +
		`\",\"port\":2,\"confirmed\":true,\"data\":\"` +
		plData + `\",\"on_busy\":\"fail\",\"tag\":\"98861544465w\"}}"
		}
	}`
	return []byte(data), nil
}

func sendData(data []byte) ([]byte, error) {
	revel.INFO.Println("SENDING DATA OVER HTTP PROTOCOL TO ", tclNetAdaptor.Client)
	req, err := http.NewRequest("POST", tclNetAdaptor.Address.IP, bytes.NewBuffer(data))

	if err != nil {
		log.Println("ERROR IN GETTING REQUEST")
		return nil, err
	}

	updateRequest(req, tclNetAdaptor.Client)

	client := &http.Client{Timeout: time.Second * 75}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("ERROR IN RESPONSE ")
		return nil, err
	}
	defer resp.Body.Close()

	var test interface{}
	json.NewDecoder(resp.Body).Decode(&test)
	blob, _ := json.Marshal(test)

	switch tclNetAdaptor.Client {
	case "tcl":
		//only for lorawan hp_iot
		log.Println("RESPONE RECEIVED FROM LORAWAN", test)
	case "devtech":
		log.Println("RESPONE RECEIVED FROM DEVTECH", test)
	}

	return blob, nil
}

//Handle for different data types
func updateRequest(req *http.Request, client string) {
	switch client {
	case "tcl":
		req.Header.Set("Content-Type", "application/vnd.onem2m-res+json;ty=4")
		req.Header.Set("X-M2M-Origin", "C5F414079-304954fa")
		req.Header.Set("X-M2M-RI", "12444328")
		req.Header.Set("Accept", "application/vnd.onem2m-res+json;")

		req.SetBasicAuth("C5F414079-304954fa", "test@123")
	case "devtech":
		req.Header.Set("Content-Type", "application/json")

		req.SetBasicAuth("chipmonk", "123")
	}
}

func chkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s \n", err.Error())
		os.Exit(1)
	}
}
