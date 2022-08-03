package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
)

const USER_AGENT = "Mozilla/5.0"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing url to repo: e.g. https://github.com/yt-dlp/yt-dlp")
		os.Exit(1)
	}

	url_ := os.Args[1]
	if strings.Count(url_, "/") < 4 {
		repos := GetRepositories(url_)
		fmt.Printf("Got %d repositories!\n\n", len(repos))

		for _, v := range repos {
			fmt.Println("Repository: " + v.Name + "\n")
			GetReleases(v.HTMLURL)
		}
	} else {
		GetReleases(url_)
	}

}

func init() {
	os.MkdirAll("downloads", 0700)
}

func GetRepositories(owner string) []Repositories {
	var repositories []Repositories

	repositories_url := fmt.Sprintf("%s/repos", strings.Replace(owner, "https://github.com", "https://api.github.com/users", -1))

	req, err := http.NewRequest(http.MethodGet, repositories_url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", USER_AGENT)
	p := url.Values{}

	p.Add("per_page", "100")
	p.Add("page", "1")
	p.Add("type", "owner")
	page := 0

	for {
		req.URL.RawQuery = p.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		var obj []Repositories
		err = json.NewDecoder(resp.Body).Decode(&obj)
		if err != nil {
			panic(err)
		}

		if len(obj) == 0 {
			break
		}

		repositories = append(repositories, obj...)

		resp.Body.Close()

		page++
		p.Set("page", fmt.Sprint(page))
	}

	return repositories
}

func GetReleases(repo string) []Releases {
	var releases []Releases

	releases_url := fmt.Sprintf("%s/releases", strings.Replace(repo, "https://github.com", "https://api.github.com/repos", -1))

	req, err := http.NewRequest(http.MethodGet, releases_url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", USER_AGENT)
	p := url.Values{}

	p.Add("per_page", "100")
	p.Add("page", "1")
	page := 0

	for {
		req.URL.RawQuery = p.Encode()

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		var obj []Releases
		err = json.NewDecoder(resp.Body).Decode(&obj)
		if err != nil {
			panic(err)
		}

		if len(obj) == 0 {
			break
		}

		releases = append(releases, obj...)

		resp.Body.Close()

		page++
		p.Set("page", fmt.Sprint(page))
	}

	fmt.Printf("Found %d releases!\n", len(releases))

	var owner, repository string
	fmt.Sscanf(repo, "https://github.com/%s/%s", &owner, &repository)

	fmt.Println()

	for _, release := range releases {
		fmt.Println("Release tag: " + release.TagName)
		directory := fmt.Sprintf("downloads/%s/%s/%s [%d]", owner, repository, release.TagName, release.ID)

		os.MkdirAll(directory, 0700)

		taresp, err := http.Get(release.TarballURL)
		if err != nil {
			fmt.Printf("Couln't get tarball: %s", err.Error())
		} else {
			f, err := os.OpenFile(directory+"/"+release.TagName+".tar.gz", os.O_CREATE|os.O_WRONLY, 0700)
			if err != nil {
				panic(err)
			}

			pb := progressbar.NewOptions64(taresp.ContentLength,
				progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionShowBytes(true),
				progressbar.OptionShowCount(),
				progressbar.OptionThrottle(100*time.Millisecond),
				progressbar.OptionSetWidth(50),
				progressbar.OptionOnCompletion(func() {
					color.Success.Println(" File downloaded! ")
				}),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        color.HiGreen.Sprint("="),
					SaucerHead:    color.HiGreen.Sprint(">"),
					SaucerPadding: color.HiYellow.Sprint("-"),
					BarStart:      "[",
					BarEnd:        "]",
				}))

			io.Copy(io.MultiWriter(f, pb), taresp.Body)

			taresp.Body.Close()
			f.Close()
		}

		notes, err := os.Create(directory + "/release-notes.md")
		if err != nil {
			panic(err)
		}
		notes.WriteString(release.Body)
		notes.Close()

		fmt.Println()

		for _, asset := range release.Assets {
			fmt.Println("Download asset: " + asset.Name)
			destination := directory + "/" + asset.Name

			resp, err := http.Get(asset.BrowserDownloadURL)
			if err != nil {
				fmt.Printf("Error on asset %s: %s\n", asset.Name, err.Error())
				continue
			}

			f, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY, 0700)
			if err != nil {
				panic(err)
			}

			pb := progressbar.NewOptions64(resp.ContentLength,
				progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionShowBytes(true),
				progressbar.OptionShowCount(),
				progressbar.OptionThrottle(100*time.Millisecond),
				progressbar.OptionSetWidth(50),
				progressbar.OptionOnCompletion(func() {
					color.Success.Println(" File downloaded! ")
				}),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        color.HiGreen.Sprint("="),
					SaucerHead:    color.HiGreen.Sprint(">"),
					SaucerPadding: color.HiYellow.Sprint("-"),
					BarStart:      "[",
					BarEnd:        "]",
				}))

			io.Copy(io.MultiWriter(f, pb), resp.Body)

			resp.Body.Close()
			f.Close()
		}
		fmt.Print("\n\n")
	}

	return releases
}
