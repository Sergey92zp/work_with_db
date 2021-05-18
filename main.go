package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Posts struct {
	UsrId      int    `json:"userId"`
	Id_post    int    `json:"id"`
	Title_post string `json:"title"`
	Body_post  string `json:"body"`
}

type Comments struct {
	PostId  int    `json:"postId"`
	Id_comm int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Body    string `json:"body"`
}

func main() {
	fmt.Println("[main] main() started")
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts?userId=7")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var post []Posts
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Fatal(err)
	}

	intCh := make(chan int)

	for i := range post {
		go writeToDB_posts(post[i].UsrId, post[i].Id_post, post[i].Title_post, post[i].Body_post)
		go getComments(intCh)
		intCh <- post[i].Id_post
		time.Sleep(5 * time.Second)
		//fmt.Println(strconv.Itoa(post[i].Id_post) + "-" + post[i].Title_post + "\n")
	}
	fmt.Println("[main] main() stopped")
}

func getComments(ch <-chan int) { //idPost int

	idPost := <-ch
	path := "https://jsonplaceholder.typicode.com/comments?postId=" + strconv.Itoa(idPost)
	//fmt.Println(path)
	resp, err := http.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var comments []Comments
	err = json.Unmarshal(body, &comments)
	if err != nil {
		log.Fatal(err)
	}

	for i := range comments {
		//fmt.Println(strconv.Itoa(comments[i].Id_comm) + "-" + comments[i].Name + "\n")
		go writeToDB_comm(comments[i].PostId, comments[i].Id_comm, comments[i].Name, comments[i].Email, comments[i].Body)
		time.Sleep(1 * time.Second)
	}

	fmt.Scanln()
}

func writeToDB_comm(post_id, comm_id int, name, email, body string) {

	db, err := sql.Open("mysql", "root:LangGo21@/postgres")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//result, err := db.Exec("Insert into postgres.comments values (?, ?, ?, ?, ?, ?)", 1, 1, "1", "1", "1") //post_id, comm_id, name, email, body)
	result, err := db.Exec("insert into postgres.comments (post_id, id, name, email, body) values (" + strconv.Itoa(post_id) + "," + strconv.Itoa(comm_id) + ",'" + name + "','" + email + "','" + body + "')")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())

	fmt.Scanln()
}

func writeToDB_posts(user_id, post_id int, title, body string) {

	db, err := sql.Open("mysql", "root:LangGo21@/postgres")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into postgres.posts (user_id, id, title, body) values (" + strconv.Itoa(user_id) + "," + strconv.Itoa(post_id) + ",'" + title + "','" + body + "')")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())

	fmt.Scanln()
}
