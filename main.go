package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"tutorials/database"
)

import _ "github.com/go-sql-driver/mysql"

type Post struct {
	UserId int32  `json:"userId"`
	Id     int32  `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

type Comment struct {
	PostId int32  `json:"postId"`
	Id     int32  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func getPostsByUserId(id int) ([]Post, error) {
	var posts []Post

	result, err := http.Get("https://jsonplaceholder.typicode.com/posts?userId=" + strconv.Itoa(id))

	if err != nil {
		return posts, err
	}

	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		return posts, err
	}

	err = json.Unmarshal(body, &posts)

	if err != nil {
		return posts, err
	}

	return posts, nil
}

func getCommentsByPostId(id int) ([]Comment, error) {
	var comments []Comment

	result, err := http.Get("https://jsonplaceholder.typicode.com/comments?postId=" + strconv.Itoa(id))

	if err != nil {
		return comments, err
	}

	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		return comments, err
	}

	err = json.Unmarshal(body, &comments)

	if err != nil {
		return comments, err
	}

	return comments, nil
}

func main() {

	db, err := database.New()

	if err != nil {
		fmt.Println(err)
		return
	}

	err = db.AutoMigrate()

	if err != nil {
		fmt.Println(err)
		return
	}

	db.Truncate()

	var wg sync.WaitGroup

	// First 10 users
	for userId := 1; userId <= 10; userId++ {
		// Go parallel fetching posts
		go func(userId int) {
			posts, err := getPostsByUserId(userId)

			if err != nil {
				fmt.Println(err)
				return
			}

			for _, post := range posts {
				fmt.Println("Post ", post.Id)

				db.Db.Create(&post)

				// Go parallel fetching comments by post
				go func(postId int) {
					comments, err := getCommentsByPostId(postId)

					if err != nil {
						fmt.Println(err)
						return
					}

					for _, comment := range comments {
						db.Db.Create(&comment)
					}

					wg.Done()
				}(int(post.Id))

				wg.Add(1)
			}

			wg.Done()
		}(userId)

		wg.Add(1)
	}

	wg.Wait()
}
