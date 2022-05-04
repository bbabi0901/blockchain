package explorer

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/bbabi0901/learngo/blockchain_3/blockchain"
)

const (
	port        string = ":4000"
	templateDir string = "explorer/templates/"
)

var templates *template.Template

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(rw, "Hello from home!") // Fprint() prints not to a console
	// tmpl, err := template.ParseFiles("templates/home.html")

	// tmpl := template.Must(template.ParseFiles("templates/home.gohtml"))
	// tmpl.Execute(rw, data)
	// 위의 line 대신에 main()에서 pages의 모든 파일을 load하기

	data := homeData{"Home", blockchain.GetBlockchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockdata")
		blockchain.GetBlockchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("Listening on http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
	// ListenAndServe() returns error if occurs. log.Fatal()은 error 받으면 실행되며, Exit(1)과 함께 종료 and error를 안받으면 실행X
}
