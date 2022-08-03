package main

type Releases struct {
	AssetsURL  string   `json:"assets_url"`
	ID         int      `json:"id"`
	TagName    string   `json:"tag_name"`
	Assets     []Assets `json:"assets"`
	TarballURL string   `json:"tarball_url"`
	ZipballURL string   `json:"zipball_url"`
	Body       string   `json:"body"`
}

type Assets struct {
	Name               string `json:"name"`
	Size               int    `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type Repositories struct {
	Name    string `json:"name"`
	Private bool   `json:"private"`
	HTMLURL string `json:"html_url"`
}
