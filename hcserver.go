package main

import (
	"net/http"
	"fmt"
	"html/template"
	"net"
	"os"
)

var reqChan, respChan chan string

type Resp struct {
	Data string
}

func tcpServerSetup(){
	service := ":13004"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	chkErr(err)
	fmt.Println("tcp server ", tcpAddr)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	chkErr(err)

	for {
		conn, err := listener.Accept()
		chkErr(err)
		fmt.Println("conn is ", conn)
		go handleConn(conn)
	}
}


func handleConn(conn net.Conn){
	defer conn.Close()

	for {
		req := <-reqChan
		len, err := conn.Write([]byte(req))
		chkErr(err)
		fmt.Println("Req writen ", req, "len :", len)
		buf := make([]byte, 1000)
		rlen, err := conn.Read(buf)
		chkErr(err)
		fmt.Println("Respons recvd: ", string(buf), "rlen :", rlen)
		respChan <- string(buf)
	}
}

var templates = template.Must(template.ParseFiles("hexcode.html", "response.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, resp interface{}){
	if err := templates.ExecuteTemplate(w, tmpl+".html", resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func viewHex(w http.ResponseWriter, r *http.Request){
	renderTemplate(w, "hexcode", "")
}

func respHandler(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	chkErr(err)

	req := r.Form.Get("hexcode")
	fmt.Println("hex code received :", req)
	reqChan <- req

	resp := <- respChan
	renderTemplate(w, "response", Resp{resp})
}

func main(){
	reqChan, respChan = make(chan string), make(chan string)

	go tcpServerSetup()

	http.HandleFunc("/", viewHex)
	http.HandleFunc("/resp", respHandler)
	fmt.Println("Started running at 8004")
	http.ListenAndServe(":8004", nil)
}

func chkErr(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %s \n", err.Error())
		os.Exit(1)
	}
}
