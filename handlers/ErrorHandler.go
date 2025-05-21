package handlerfuncitons

import (
	"fmt"
	"html/template"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	t, err := template.ParseFiles("Templates/404.html") // Parse the internal server error template file
	if err != nil {
		fmt.Println("Error Parsing Files: ", err)
		http.Error(w, "Internal Server Error \n Error in parsing 500 HTML template..", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil) // Execute the template and write the output to the response writer
	if err != nil {
		fmt.Println("Error executing: ", err)
		http.Error(w, "Internal Server Error \n Error in executing 404 HTML template..", http.StatusInternalServerError)
		return
	}
	
}

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	//http.ServeFile(w, r, "500.html")
	t, err := template.ParseFiles("Templates/500.html") // Parse the internal server error template file
	if err != nil {
		fmt.Println("Error Parsing Files: ", err)
		http.Error(w, "Internal Server Error \n Error in parsing 500 HTML template..", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, nil) // Execute the template and write the output to the response writer
	if err != nil {
		fmt.Println("Error executing: ", err)
		http.Error(w, "Internal Server Error \n Error in executing 500 HTML template..", http.StatusInternalServerError)
		return
	}
	
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	//http.ServeFile(w, r, "405.html")

	t, err := template.ParseFiles("Templates/405.html") // Parse the error template file
	if err != nil {
		fmt.Println("Error Parsing Files: ", err)
		InternalServerError(w,r)
		return
	}
	err = t.Execute(w, nil) // Execute the template and write the output to the response writer
	if err != nil {
		fmt.Println("Error executing: ", err)
		InternalServerError(w,r)
		return
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	//http.ServeFile(w, r, "400.html")
	t, err := template.ParseFiles("Templates/400.html") // Parse the error template file
	if err != nil {
		fmt.Println("Error Parsing Files: ", err)
		InternalServerError(w,r)
		return
	}
	err = t.Execute(w, nil) // Execute the template and write the output to the response writer
	if err != nil {
		fmt.Println("Error executing: ", err)
		InternalServerError(w,r)
		return
	}
}
