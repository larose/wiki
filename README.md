# Personal Wiki

This project is a personal wiki written in Go.

- Markup language: Markdown
- Storage: Git


# Compiling

With `go get`:

`$ go get github.com/larose/wiki`

or from the source directory:

```
$ go get
$ make wiki
```


# Running

Execute:

`$ ./wiki`

and open your browser to:

`http://127.0.0.1:8000`

Options:

- addr: TCP address to listen on. Default: 127.0.0.1:8000
- data-dir: Data directory. Default: data

Example:

`$ ./wiki --addr 127.0.0.1:8888 --data-dir /path/to/git/repo`


# Deploying

1. Install `git` on your server.
2. Copy the binary to your server and run it.

Notes:

- The css, js and template files are embedded in the binary.
- You should run the wiki behind a reverse proxy with authentication.


# Keyboard Shortcuts

## Editing

| Action  | Shortcut        |
| :-----: |:---------------:|
| diff    | `alt+shift + d` |
| edit    | `alt+shift + e` |
| preview | `alt+shift + p` |
| save    | `alt+shift + s` |

# Contributing

Looking to contribute? Here are some ideas:

- Add tests.

- Search in title: the search functionality looks only in the body of
  pages. It should also look in the title of the pages.

- Add file uploads.

- Improve visual diff: it's currently the output of `git diff`.

- Add diff between two revisions of a page.

- Add help page: should be written in markdown and located at `/_/help`.

- Add logs.

Send me an email if you have any questions.


# Author

Mathieu Larose <mathieu@mathieularose.com>


# License

See the LICENSE file for details.
