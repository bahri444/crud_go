package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Mahasiswa struct {
	Id   int
	Nim  string
	Nama string
	Jk   string
}

var tpl *template.Template
var db *sql.DB

func main() {
	tpl, _ = template.ParseGlob("template/*.html")
	var err error
	db, err = sql.Open("mysql", "root@tcp(127.0.0.1:3306/mhs)")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	http.HandleFunc("/read", read)
	http.HandleFunc("/insert", insert)
	http.HandleFunc("/update", update)
	http.HandleFunc("/updateRes", updateRes)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/", home)
	http.ListenAndServe("localhost:8080", nil)
}
func read(w http.ResponseWriter, r *http.Request) {
	fmt.Println("reading data mahasiswa")
	Querys := "select *from mahasiswa"
	rows, err := db.Query(Querys)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var each []Mahasiswa
	for rows.Next() {
		var m Mahasiswa
		err = rows.Scan(&m.Id, &m.Nim, &m.Nama, &m.Jk)
		if err != nil {
			fmt.Println(err)
		}
		each = append(each, m)
	}
	tpl.ExecuteTemplate(w, "read.html", each)
}
func insert(w http.ResponseWriter, r *http.Request) {
	fmt.Println("insert data mahasiwa")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "insert.html", nil)
		return
	}
	r.ParseForm()
	nim := r.FormValue("nim")
	nama := r.FormValue("nama")
	jk := r.FormValue("jk")
	var err error
	if nim == "" || nama == "" || jk == "" {
		fmt.Println("error inserting data row", err)
		tpl.ExecuteTemplate(w, "insert.html", "error inserting data please check all field and try again")
		return
	}
	var ins *sql.Stmt
	ins, err = db.Prepare("insert into `mahasiswa`(`nim`,`nama`,`jk`)values(?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	defer ins.Close()
	res, err := ins.Exec(nim, nama, jk)
	rowsAffect, _ := res.RowsAffected()
	if err != nil || rowsAffect != 1 {
		fmt.Println("error inserting data row", err)
		tpl.ExecuteTemplate(w, "insert.html", "error inserting data please check all field and try again")
		return
	}
	lastInserted, _ := res.LastInsertId()
	rowsAffected, _ := res.RowsAffected()
	fmt.Println("id succes insert", lastInserted)
	fmt.Println("all rows succes insert", rowsAffected)
	tpl.ExecuteTemplate(w, "insert.html", "data succesfully inserted")
}
func update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("update data mahasiswa")
	r.ParseForm()
	id := r.FormValue("id mahasiswa")
	rows := db.QueryRow("select*from mahasiswa where id=?", id)
	var m Mahasiswa
	err := rows.Scan(&m.Id, &m.Nim, &m.Nama, &m.Jk)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/read", 307)
		return
	}
	tpl.ExecuteTemplate(w, "update.html", m)
}
func updateRes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("updateResult data mahasiswa")
	r.ParseForm()
	id := r.FormValue("id")
	nim := r.FormValue("nim")
	nama := r.FormValue("nama")
	jk := r.FormValue("jk")
	upQuerys := "update `mahasiswa` set `nim=?`,`nama`=?,`jk`=? where (`id`=?)"
	stmt, err := db.Prepare(upQuerys)
	if err != nil {
		fmt.Println("error preparing query update")
		fmt.Println(err.Error())
	}
	defer stmt.Close()
	var res sql.Result
	res, err = stmt.Exec(nim, nama, jk, id)
	rowsAff, _ := res.RowsAffected()
	if err != nil || rowsAff != 1 {
		fmt.Println(err)
		tpl.ExecuteTemplate(w, "update.html", "problem is update result mahasiswa")
		return
	}
	tpl.ExecuteTemplate(w, "update.html", "update res")

}
func delete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete data mahasiswa")
	r.ParseForm()
	id := r.FormValue("id mahasiswa")
	del, err := db.Prepare("delete from `mahasiswa` where (`id`=?)")
	if err != nil {
		fmt.Println(err)
	}
	defer del.Close()
	var res sql.Result
	res, err = del.Exec(id)
	rowsAff, _ := res.RowsAffected()
	fmt.Println("rowsAff", rowsAff)

	if err != nil || rowsAff != 1 {
		fmt.Fprintln(w, "error deleted")
		return
	}
	fmt.Println("error del:", err)
	tpl.ExecuteTemplate(w, "delete.html", "data mahasiswa succesfully deleted")
}
func home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/read", 307)
}
