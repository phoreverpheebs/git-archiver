<h1>git-archiver</h1>

A simple repository archiver written in Go

* Saves all assets in every release of a repository
* Saves the release notes in release-notes.md
* Saves the source code tarball
* Does all the above for all repositories of a user, if a user is provided

<h2>Usage</h2>

`git-archiver <repo-url or user-url>`

repo-url: `https://github.com/{user}/{repo}`<br>
user-url: `https://github.com/{user}`

<h2>Installation</h2>

Go to releases of this github (ironic I know) and download the executable for your operating system

Build from source:<br>
```
git clone https://github.com/phoreverpheebs/git-archiver
cd git-archiver
go build .
```