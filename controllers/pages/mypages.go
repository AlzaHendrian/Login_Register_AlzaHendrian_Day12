package pages

import (
	"html/template"
	"net/http"
)

func AddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("./views/projectAdd.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var data = map[string]interface{}{
		"title" : "Add Project | Marcel",
		"isLogin" : true,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}



func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("./views/contact-me.html")
	// error handling
	if err != nil {
		panic(err)
	}

	var data = map[string]interface{}{
		"title" : "Contact Me | Marcel",
		"isLogin" : true,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}


func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("./views/pageLogin.html")
	// error handling
	if err != nil {
		panic(err)
	}

	var data = map[string]interface{}{
		"title" : "Login | Marcel",
		"isLogin" : true,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}



func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("./views/pageRegister.html")
	// error handling
	if err != nil {
		panic(err)
	}

	var data = map[string]interface{}{
		"title" : "Register | Marcel",
		"isLogin" : true,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}