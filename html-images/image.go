package images

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

const needTokenHost = ".qiita.com"

type PostImages struct {
	DirName          string
	ImageDownloaders []ImageDownloader
}

type ImageDownloader struct {
	Url      *url.URL
	FileName string
}

func Images(dir string, title string, renderedBody string) PostImages {
	dirName := dir + "/" + title + "/"
	return PostImages{DirName: dirName, ImageDownloaders: setImageDownloaders(renderedBody, dirName)}
}

func (p PostImages) Download(token string) (int, int) {
	success, fail := 0, 0

	if len(p.ImageDownloaders) > 0 {
		if err := os.MkdirAll(p.DirName, 0777); err != nil {
			println(err.Error())
			panic(err)
		}
	}

	for _, imageDownloader := range p.ImageDownloaders {
		if imageDownloader.save(token) {
			success++
		} else {
			fail++
		}
	}

	return success, fail
}

func (i ImageDownloader) save(token string) bool {
	println("Downloading: " + i.Url.String())

	if i.Url.Host == "" {
		println("Wrong host...")
		return false
	}

	req, _ := http.NewRequest("GET", i.Url.String(), nil)

	if i.needToken(needTokenHost) {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	httpClient := new(http.Client)
	resp, err := httpClient.Do(req)

	if err != nil {
		println(err.Error())
		return false
	}

	defer resp.Body.Close()

	if err != nil {
		println(err.Error())
		return false
	}

	if resp.StatusCode != http.StatusOK {
		println("Status Code: " + resp.Status)
		return false
	}

	file, err := os.Create(i.FileName)
	defer file.Close()

	if err != nil {
		println(err.Error())
		return false
	}

	if _, err := io.Copy(file, resp.Body); err != nil {
		println(err.Error())
		return false
	}

	println("Saved " + i.FileName)

	return true
}

func setImageDownloaders(renderedBody string, dirName string) []ImageDownloader {
	var imageDownloaders []ImageDownloader

	re := regexp.MustCompile(`<img.+?src="(.+?)".*?>`)
	imageURLs := re.FindAllStringSubmatch(renderedBody, -1)

	for _, imageURL := range imageURLs {
		parsedURL, err := url.Parse(imageURL[1])
		if err != nil {
			panic(err)
		}

		if parsedURL.Scheme == "" {
			parsedURL.Scheme = "http"
		}

		imageDownloader := ImageDownloader{Url: parsedURL, FileName: fileName(parsedURL.String(), dirName)}
		imageDownloaders = append(imageDownloaders, imageDownloader)
	}

	return imageDownloaders
}

func fileName(url string, dirName string) string {
	_, filename := path.Split(url)

	return dirName + filename
}

func (i ImageDownloader) needToken(host string) bool {
	return strings.Contains(i.Url.Host, host)
}
