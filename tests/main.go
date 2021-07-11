package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL    = "http://localhost:8082"
	sm9TestTag = "今年レンコンコマンダー常盤"
)

func main() {
	// Can we login?
	client := authenticate("admin", "admin")
	log.Println("Logged in successfully")
	// Can we try to archive something?
	makeArchiveRequest(client, "bilibili", "tag", "sm35952346")
	makeArchiveRequest(client, "niconico", "tag", "今年レンコンコマンダー常盤")
	makeArchiveRequest(client, "youtube", "channel", "UCF43Xa8ZNQqKs1jrhxlntlw") // Some random channel I found with short videos. Good enough!

	time.Sleep(time.Minute * 30)

	// Are videos being downloaded and transcoded correctly?
	pageHasVideos(client, "今年レンコンコマンダー常盤", 1)
	pageHasVideos(client, "sm35952346", 1)
	pageHasVideos(client, "しゅんなな", 13)
	log.Println("All videos downloaded and transcoded successfully")
}

func pageHasVideos(client *http.Client, tag string, count int) {
	response, _ := client.Get(baseURL + fmt.Sprintf("/?search=%s&category=upload_date", tag))
	cont, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	c := strings.Count(string(cont), "href=\"/videos/")
	if c != count {
		log.Panicf("page does not contain the right number of videos for %s. Found: %d", tag, c)
	}
}

func makeArchiveRequest(client *http.Client, website, contentType, contentValue string) {
	response, _ := client.PostForm(baseURL+"/archiverequests", url.Values{
		"website":      {website},
		"contentType":  {contentType},
		"contentValue": {contentValue},
	})

	if response.StatusCode != 301 {
		log.Fatalf("bad archival request response status: %d", response.StatusCode)
	}

	return
}

var redirectErr error = errors.New("don't redirect")

func authenticate(username, password string) *http.Client {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return redirectErr
		},
		Jar: jar,
	}

	response, _ := client.PostForm(baseURL+"/login", url.Values{
		"username": {username},
		"password": {password},
	})
	// lol how do i check for an error here? :thinking:
	// if err != nil && err != redirectErr {
	// 	log.Panicf("failed to post with err: %s", err)
	// }

	if response.StatusCode != 301 {
		log.Panicf("bad auth status code: %d", response.StatusCode)
	}

	jwt := ""
	for _, cookie := range response.Cookies() {
		if cookie.Name == "jwt" {
			jwt = cookie.Value
		}
	}

	if jwt == "" {
		log.Panicf("jwt cookie not set")
	}

	return client
}
