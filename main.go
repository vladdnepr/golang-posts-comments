package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

import "database/sql"
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

func savePostToDb(post Post, db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO posts (id, user_id, title, body) VALUES (?, ?, ?, ?)",
		post.Id,
		post.UserId,
		post.Title,
		post.Body,
	)

	return err
}

func saveCommentToDb(comment Comment, db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO comments (id, post_id, name, email, body) VALUES (?, ?, ?, ?, ?)",
		comment.Id,
		comment.PostId,
		comment.Name,
		comment.Email,
		comment.Body,
	)

	return err
}

func main() {
	db, err := sql.Open("mysql", "root:root@/golang")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	// Clear DB
	_, _ = db.Exec("TRUNCATE posts")
	_, _ = db.Exec("TRUNCATE comments")

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

				err = savePostToDb(post, db)

				if err != nil {
					fmt.Println(err)
					return
				}

				// Go parallel fetching comments by post
				go func(postId int) {
					comments, err := getCommentsByPostId(postId)

					if err != nil {
						fmt.Println(err)
						return
					}

					for _, comment := range comments {
						err = saveCommentToDb(comment, db)

						if err != nil {
							fmt.Println(err)
							return
						}
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
