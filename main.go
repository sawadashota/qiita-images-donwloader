package main

import (
	"download-qiitateam-image/html-images"
	"flag"
	"github.com/mitchellh/go-homedir"
	"github.com/sawadashota/qiita-posts-go"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	qiitaStatusCode int
	qiitaPosts      []qiita.Post
	qiitaPost       qiita.Post
	key             int
	processID       int
	shouldRetry     bool
)

const QiitaRequestInterval = 3.6

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	qiitaTeamName := fs.String("q", "", "Qiita::Team Name")
	qiitaToken := fs.String("qToken", "", "Qiita::Team Access Token")
	downloadDir := fs.String("dir", "", "Download directory")
	restartFromStr := fs.String("restart-from", "1", "Restart from process ID")
	fs.Parse(os.Args[1:])

	if *qiitaTeamName == "" {
		panic("Please type -q (Qiita::Team Name)")
	}

	if *qiitaToken == "" {
		panic("Please type -qToken (Qiita::Team Access Token)")
	}

	if *downloadDir == "" {
		*downloadDir, _ = homedir.Expand("~/Downloads/qiita-team-images")
	}

	restartFrom, err := strconv.Atoi(*restartFromStr)
	if err != nil {
		panic("-restart-from should be integer.")
	}

	startingQiitaPage := postPage(restartFrom, qiita.PagePerPost)

	for i := startingQiitaPage; ; i++ {
		qiitaStatusCode, qiitaPosts = qiita.Posts(i, *qiitaTeamName, *qiitaToken).Get()

		if qiitaStatusCode != 200 {
			println("-----------------------------------------------")
			println("Qiita Status Code: " + strconv.Itoa(qiitaStatusCode))
			println("-----------------------------------------------")

			// 502だったら再挑戦
			if qiitaStatusCode != http.StatusBadGateway && shouldRetry {
				println("Retrying...")
				time.Sleep(QiitaRequestInterval * 1000 * time.Millisecond)
				shouldRetry = false
				i--
				continue
			}
			break
		}

		shouldRetry = true

		for key, qiitaPost = range qiitaPosts {
			processID = (qiita.PagePerPost * (i - 1)) + (key + 1)

			if processID < restartFrom {
				continue
			}

			println("-----------------------------------")
			println(strconv.Itoa(processID) + ". Processing: " + qiitaPost.Title)
			images.Images(*downloadDir, qiitaPost.Title, qiitaPost.RenderedBody).Download(*qiitaToken)
		}

		// ループ終了条件
		if len(qiitaPosts) < qiita.PagePerPost {
			println("-----------------------------------------------")
			println(strconv.Itoa(processID-restartFrom) + " posts processed!")
			println("-----------------------------------------------")
			break
		}

	}

}

// Qiitaは何ページ目から読み込むか
func postPage(processId int, pagePerPost int) int {
	return int(math.Floor(float64((processId-1)/pagePerPost))) + 1
}
