package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"projek/config"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// data User
type Users struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Project struct {
	Id           int
	Title        string
	Sdate        time.Time
	Edate        time.Time
	Duration     string
	Content      string
	Technologies []string
	Tnode        bool
	Treact       bool
	Tjs          bool
	Thtml        bool
	User         string
}

// data login
var Data = map[string]interface{}{
	"Title":     "Personal Web",
	"IsLogin":   true,
	"Id":        1,
	"UserName":  "Alza",
	"FlashData": "",
}

func main() {
	// menyiapkan routingan
	router := mux.NewRouter()

	// create connection database.go
	config.ConnectDB()

	// create static folder for public
	router.PathPrefix("/public").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	// routing pages
	router.HandleFunc("/", Home).Methods("GET")
	router.HandleFunc("/addProject", AddProject).Methods("GET")
	router.HandleFunc("/contact", Contact).Methods("GET")
	router.HandleFunc("/register", FormRegister).Methods("GET")
	router.HandleFunc("/login", FormLogin).Methods("GET")

	// routing actions
	router.HandleFunc("/add", Add).Methods("POST")
	router.HandleFunc("/update/{id}", Update).Methods("GET")
	router.HandleFunc("/upost/{id}", UpdatePost).Methods("POST")
	router.HandleFunc("/delete/{id}", Delete).Methods("GET")
	router.HandleFunc("/detail/{id}", Detail).Methods("GET")

	// routing auth and session
	router.HandleFunc("/register", Register).Methods("POST")
	router.HandleFunc("/login", Login).Methods("POST")

	// create port
	port := "5000"
	fmt.Println("server running on port", port)
	http.ListenAndServe("localhost:"+port, router)

	// fmt.Println("Server Running on port 5000")
	// http.ListenAndServe("localhost:5000", router)
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")
	var tmpl, err = template.ParseFiles("./views/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// session get SESSION_ID
	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	// conditional is login
	var readDT string
	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
		readDT = "SELECT id, title, start_date, end_date, description, technology, author_id FROM public.tb_project;"
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["UserName"] = session.Values["Name"].(string)
		readDT = "SELECT tb_project.id, title, start_date, end_date, description, technology, tb_user.name as user FROM tb_project LEFT JOIN tb_user ON tb_project.author_id = tb_user.id WHERE tb_user.name='" + Data["UserName"].(string) + "' ORDER BY id DESC"
	}

	rows, _ := config.ConnDB.Query(context.Background(), readDT)

	var result []Project
	for rows.Next() {
		var each = Project{}
		var err = rows.Scan(&each.Id, &each.Title, &each.Sdate, &each.Edate, &each.Content, &each.Technologies, &each.User)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// duration
		duration := each.Edate.Sub(each.Sdate)
		var distance string
		if duration.Hours()/24 < 7 {
			distance = strconv.FormatFloat(duration.Hours()/24, 'f', 0, 64) + " Days"
		} else if duration.Hours()/24/7 < 4 {
			distance = strconv.FormatFloat(duration.Hours()/24/7, 'f', 0, 64) + " Weeks"
		} else if duration.Hours()/24/30 < 12 {
			distance = strconv.FormatFloat(duration.Hours()/24/30, 'f', 0, 64) + " Months"
		} else {
			distance = strconv.FormatFloat(duration.Hours()/24/30/12, 'f', 0, 64) + " Years"
		}
		each.Duration = distance

		result = append(result, each)
	}

	resp := map[string]interface{}{
		"Data":     Data,
		"Projects": result,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func FormRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("./views/pageRegister.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var data = map[string]interface{}{
		"title":   "Register | Alza",
		"isLogin": true,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("pass")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	postUSER := "INSERT INTO public.tb_user(name, email, password) VALUES ($1, $2, $3);"
	_, err = config.ConnDB.Exec(context.Background(), postUSER, name, email, passwordHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("pass")
	user := Users{}
	selectUSER := "SELECT id, name, email, password FROM tb_user WHERE email=$1"
	rows := config.ConnDB.QueryRow(context.Background(), selectUSER, email)
	err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	session.Values["IsLogin"] = true
	session.Values["Name"] = user.Name
	session.Values["Id"] = user.Id
	session.Options.MaxAge = 10800

	// flash login
	session.AddFlash("Login Success", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func FormLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("./views/pageLogin.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var data = map[string]interface{}{
		"title":   "Login | Alza",
		"isLogin": true,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, data)
}

func Add(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("pname")
	content := r.PostForm.Get(("desc"))
	SD := r.PostForm.Get("sdate")
	ED := r.PostForm.Get("edate")
	tech := r.Form["check"]

	addID := "INSERT INTO tb_project(title, start_date, end_date, description, technology) VALUES ($1, $2, $3, $4, $5)"
	config.ConnDB.Exec(context.Background(), addID, title, SD, ED, content, tech)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")

	var tmpl, err = template.ParseFiles("./views/projectUpdate.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	selectID := "SELECT id, title, start_date, end_date, description, technology FROM tb_project WHERE id=$1"
	rows := config.ConnDB.QueryRow(context.Background(), selectID, id)
	var getID Project
	err = rows.Scan(&getID.Id, &getID.Title, &getID.Sdate, &getID.Edate, &getID.Content, &getID.Technologies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	result := []Project{getID}

	var SD = getID.Sdate.Format("2006-01-02")
	var ED = getID.Edate.Format("2006-01-02")

	var data = map[string]interface{}{
		"title":   "Edit Project",
		"isLogin": true,
	}

	resp := map[string]interface{}{
		"Data":      data,
		"GetUpdate": result,
		"SD":        SD,
		"ED":        ED,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("pname")
	content := r.PostForm.Get(("desc"))
	SD := r.PostForm.Get("sdate")
	ED := r.PostForm.Get("edate")
	// tech := r.Form["check"]
	// image := r.Context().Value("dataFile")
	// img := image.(string)

	// sDateFormat := SD.Format("02 January 2006")
	// eDateFormat := ED.Format("02 January 2006")

	update := "UPDATE public.tb_project SET title=$1, start_date=$2, end_date=$3, description=$4 WHERE id=$5"
	_, err = config.ConnDB.Exec(context.Background(), update, title, SD, ED, content, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	deleteID := "DELETE FROM tb_project WHERE id=$1"
	_, err := config.ConnDB.Exec(context.Background(), deleteID, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Detail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	// parsing template html
	var tmpl, err = template.ParseFiles("./views/projectDetail.html")
	// error handling
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	}

	// manangkap id (id, _ (tanda _ tidak ingin menampilkan eror))
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	selectID := "SELECT id, title, start_date, end_date, description, technology FROM tb_project WHERE id=$1"
	detail := Project{}
	rows := config.ConnDB.QueryRow(context.Background(), selectID, id)
	err = rows.Scan(&detail.Id, &detail.Title, &detail.Sdate, &detail.Edate, &detail.Content, &detail.Technologies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// convert date ke string
	sDateFormat := detail.Sdate.Format("02 January 2006")
	eDateFormat := detail.Edate.Format("02 January 2006")

	// duration
	duration := detail.Edate.Sub(detail.Sdate)
	var distance string
	if duration.Hours()/24 < 7 {
		distance = strconv.FormatFloat(duration.Hours()/24, 'f', 0, 64) + " Days"
	} else if duration.Hours()/24/7 < 4 {
		distance = strconv.FormatFloat(duration.Hours()/24/7, 'f', 0, 64) + " Weeks"
	} else if duration.Hours()/24/30 < 12 {
		distance = strconv.FormatFloat(duration.Hours()/24/30, 'f', 0, 64) + " Months"
	} else {
		distance = strconv.FormatFloat(duration.Hours()/24/30/12, 'f', 0, 64) + " Years"
	}

	// technology
	node, react, js, html := false, false, false, false
	tech := detail.Technologies
	for _, i := range tech {
		switch i {
		case "node":
			node = true
		case "react":
			react = true
		case "js":
			js = true
		case "html5":
			html = true
		}
	}

	resp := map[string]interface{}{
		"Data":     Data,
		"Projects": detail,
		"Duration": distance,
		"T1":       node,
		"T2":       react,
		"T3":       js,
		"T4":       html,
		"SD":       sDateFormat,
		"ED":       eDateFormat,
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, resp)
}

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

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")
	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
		Data["Id"] = session.Values["Id"]
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["UserName"] = session.Values["Name"].(string)
		Data["Id"] = session.Values["Id"].(int)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}

func Contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset-utf-8")

	var tmpl, err = template.ParseFiles("./views/contact-me.html")
	// error handling
	if err != nil {
		panic(err)
	}

	var store = sessions.NewCookieStore([]byte("SESSION_ID"))
	session, _ := store.Get(r, "SESSION_ID")

	if session.Values["IsLogin"] != true {
		Data["IsLogin"] = false
		Data["Id"] = session.Values["Id"]
	} else {
		Data["IsLogin"] = session.Values["IsLogin"].(bool)
		Data["UserName"] = session.Values["Name"].(string)
		Data["Id"] = session.Values["Id"].(int)
	}

	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, Data)
}
