// REST API; go data를 json으로 변환 후(using Marshal and Unmarshal) json을 사용자에게 보낸다.
// documentation -> using route("/") -> / 에서 API에서 할 수 있는 일들의 목록을 볼 수 있다.

package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/bbabi0901/blockchain/blockchain"
	"github.com/bbabi0901/blockchain/utils"
	"github.com/gorilla/mux"
)

var port string

type url string

type urlDescription struct {
	URL         url    `json:"url"`    // struct field tag; struct가 json일 경우 ""안의 형식대로 표기. 사용시 백틱(``)으로 감싸고.
	Method      string `json:"method"` // if you want to ignore the Field, use "-".
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"` // omitempty를 추가시 값이 비어있을 경우, 표기하지 않는다.
}

type addBlockBody struct {
	Message string
}

type errResponse struct {
	ErrMessage string `json:"errorMessage"`
}

// String()은 fmt에 있는 모듈. 근데 그걸 URLDescription에 한하여 재정립.
// URLDescription이 출력될 때 마다 String()은 fmt에서 자동으로 호출되서 function의 내용대로 출력된다.
// 그니까, 특정 type에 한하여 기존의 function이 실행될 때 끼어들어서 우리 입맛대로 실행하게끔 하는 것.
// func (u URLDescription) String() string {
// 	return "Printing URLDescription in my way"
// }

// MarshalText()는 josn.NewEncode()에서 Marshal할 때 실행되는 함수(source code 확인시 볼 수 있다)
// URL에 MarshalText()가 실행될 때, 얘가 별도로 작용. 그러면 이제 URL은 아래의 data대로 실행되는게 아니라 정의한 내용을 출력.
// 이 implementing은 함수명, 리턴값까지 완벽히 같아야 한다.
func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"), // bc implement is defined, URL("/") = fmt.Sprintf("http://localhost%s%s", port, "/")
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data: string",
		},
		{
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See A Block",
		},
	}
	// rw.Header().Add("Content-Type", "application/json") // sending json response
	// version 1; go data -> json
	// b, err := json.Marshal(data)
	// utils.HandleErr(err)
	// fmt.Fprintf(rw, "%s", b)
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockchain().AllBlocks())
	case "POST":
		var addBlockBody addBlockBody                                  // addBlockBody라는 AddBlockBody의 type의 빈 변수를 만들어서
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody)) // decode한 request의 body를 addBlockBody에 넣는다. 이 때, 꼭 pointer를 인자로!
		blockchain.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
		// default:
		// 	rw.WriteHeader(http.StatusMethodNotAllowed) no needed if using gorilla mux
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["height"])
	utils.HandleErr(err)
	block, err := blockchain.GetBlockchain().Block(id)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}

}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func Start(aPort int) {
	// handler := http.NewServeMux() // ServeMux takes url(ex. "/blocks") and make link to func blocks().
	port = fmt.Sprintf(":%d", aPort)

	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)                // middleWare; adding "application/json" to every request header
	router.HandleFunc("/", documentation).Methods("GET") // 어떤 Method를 처리할지 특정할 수 있다 with gorilla mux
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")

	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
