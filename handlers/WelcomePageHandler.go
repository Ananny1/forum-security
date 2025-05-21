package handlerfuncitons

import (
	"fmt"
	"html/template"
	"net/http"
)

func Welcomepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/welcome" {
		filename := "Templates/welcomepage.html"
		t, err := template.ParseFiles(filename) // Parse the index.html template file
		if err != nil {
			fmt.Println("error parsing:", err)
			return
		}
		err = t.ExecuteTemplate(w, "welcomepage.html", nil) // Execute the template and write the output to the response writer // Redirect the user to the homepage after 7 seconds
		if err != nil {
			fmt.Println("error executing:", err)
			InternalServerError(w, r)
            return
		}
	}
}
