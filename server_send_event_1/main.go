package main

import (
	"bytes"
	"encoding/json"
	"math/rand"

	"fmt"
	"log"
	"net/http"
)

type (
	Currency float64

	Item struct {
		// omitempty: json 객체를 생성할 때 해당 필드에 값이 없으면 필드를 변환하지 않고 건너뛴다
		// `json:"-"`로 설정 시 해당 필드는 무조건 변환되지 않는다
		Name		string		`json:"name,omitempty""`
		Quantity	int		`json:"quantity,omitempty"`
		Price		Currency	`json:"price,omitempty"`
	}

	Store struct {
		Items	map[string]Item	`json:"items,omitempty"`

	}

	DashBoard struct {
		Users		uint		`json:"users,omitempty"`
		UsersLoggedIn	uint		`json:"users_logged_in,omitempty"`
		Inventory	*Store		`json:"inventory,omitempty"`
		ChartOne	[]int		 `json:"chart_one,omitempty"`
		ChartTwo	[]Currency	 `json:"chart_two,omitempty"`
	}
)

var dashboard chan *DashBoard

func main() {

	// chanel은 Queue로 생각하면 된다(FIFO: First In First Out)
	// pop(queue에 넣기): channel <- data
	// push(queue에서 빼기): var <- channel
	dashboard = make(chan *DashBoard)
	go updateDashboard()

	// register static files handle '/index.html -> client/index.html'
	http.Handle("/", http.FileServer(http.Dir("client")))
	// register RESTful handler for '/sse/dashboard'
	// HandleFunc: Handle을 function으로 사용한다
	http.HandleFunc("/sse/dashboard", dashbaordHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
// http.ResponseWriter는 HTTP Response에 무언가를 쓸 수 있게 한다
func dashbaordHandler(w http.ResponseWriter, r *http.Request) {
	// sse를 위한 설정
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	var buf bytes.Buffer
	// encoder를 생성하고 앞으로 입력될 데이터를 &buf로 연결하는 stream을 갖는다
	enc := json.NewEncoder(&buf)
	// json.Encoder 타입은 stream 기반으로 문자열을 만든다 
	// Queue(channel) dashboard의 값을 JSON으로 변환하여 &buf로 전달한다 
	enc.Encode(<-dashboard)

	// w(http.ResponseWriter)에 buf(bytes)를 string으로 변환하여 전달
	fmt.Fprintf(w, "data: %v\n\n", buf.String())
	// command 창에 buf(bytes)를 string으로 변환하여 전달
	fmt.Printf("data: %v\n", buf.String())
}

func updateDashboard() {
	for {
		inv := updateInventory()
		db := &DashBoard {
			Users:		uint(rand.Uint32()),
			UsersLoggedIn:	uint(rand.Uint32() % 200),
			Inventory:	inv,
			ChartOne:	[]int{4, 22, 523, 66, 23454},
			ChartTwo:	[]Currency{1.23, 6.54, 366.34},
		}
		dashboard <- db
	}
}

func updateInventory() *Store {
	inv := &Store{}
	// string을 key로 갖고 Item을 value로 갖는 map을 생성
	inv.Items = make(map[string]Item)

	a := Item{Name: "Books", Price: 33.59, Quantity: int(rand.Int31() % 53)}
	// book을 key로 갖는 map을 생성하여 a 값을 value로 저장
	inv.Items["book"] = a

	a = Item{Name: "Bicycles", Price: 190.89, Quantity: int(rand.Int31() % 232)}
        inv.Items["bicycle"] = a

	a = Item{Name: "RC Car", Price: 83.19, Quantity: int(rand.Int31() % 73)}
        inv.Items["rccar"] = a

	return inv

}
