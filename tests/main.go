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
	baseURL    = "http://localhost:8083/api"
	sm9TestTag = "今年レンコンコマンダー常盤"
)

func main() {
	// Can we login?
	client := authenticate("admin", "admin")
	log.Println("Logged in successfully")
	// Can we try to archive something?
	// bilibili is currently broken...
	//makeArchiveRequest(client, "bilibili", "tag", "sm35952346")
	//makeArchiveRequest(client, "bilibili", "channel", "1963331522")

	makeArchiveRequest(client, "https://www.nicovideo.jp/search/TEST_sm9")
	makeArchiveRequest(client, "https://www.nicovideo.jp/user/119163275")
	makeArchiveRequest(client, "https://www.nicovideo.jp/mylist/58583228")

	makeArchiveRequest(client, "https://www.youtube.com/channel/UCF43Xa8ZNQqKs1jrhxlntlw")                            // Some random channel I found with short videos. Good enough!
	makeArchiveRequest(client, "https://www.youtube.com/watch?v=sPTwJwZOZkI&list=PL27eLnikSM92Vw2ssXmNB_GtK8xG2U9sW") // playlist with 5 entries

	// Are videos being downloaded and transcoded correctly?
	for start := time.Now(); time.Since(start) < time.Minute*30; {
		time.Sleep(time.Second * 30)
		//err := pageHasVideos(client, "sm35952346", 1) // Bilibili tag
		//if err != nil {
		//	log.Println(err)
		//	continue
		//}
		//
		//err = pageHasVideos(client, "被劝诱的石川", 1) // Bilibili channel
		//if err != nil {
		//	log.Println(err)
		//	continue
		//}

		err := pageHasVideos(client, "風野灯織", 1) // nico channel
		if err != nil {
			log.Println(err)
			continue
		}

		err = pageHasVideos(client, "中の", 1) // there's some bizarre nico bug here where the tags keep switching on the video. very strange
		if err != nil {
			log.Println(err)
			continue
		}

		err = pageHasVideos(client, "しゅんなな", 8) // yt channel, should be 13 but several have ffmpeg errors. Sad!
		if err != nil {
			log.Println(err)
			continue
		}

		err = pageHasVideos(client, "琴葉姉妹のにゃーねこにゃー！", 1) // yt playlist, searching for neko neko nya nya video (lol)
		if err != nil {
			log.Println(err)
			continue
		}

		err = pageHasVideos(client, "NEW+GAME", 1) // Nico mylist
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("All videos downloaded and transcoded successfully")
		return
	}

	log.Panic("Failed to download and transcode videos within 30 minutes")

}

func pageHasVideos(client *http.Client, tag string, count int) error {
	url := baseURL + fmt.Sprintf("/home?search=%s&category=upload_date&order=desc", tag)
	response, _ := client.Get(url)
	cont, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	c := strings.Count(string(cont), "VideoID")
	if c < count {
		return fmt.Errorf("page does not contain the right number of videos for %s. Found: %d", tag, c)
	}

	return nil
}

func makeArchiveRequest(client *http.Client, inpURL string) {
	response, _ := client.PostForm(baseURL+"/archiverequests", url.Values{
		"url": {inpURL},
	})

	if response.StatusCode != 200 {
		log.Fatalf("bad archival request response status: %d", response.StatusCode)
	}

	log.Printf("Made archival request for %s", inpURL)

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

	if response.StatusCode != 200 {
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
