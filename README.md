# commenter 

A lightweight HTTP comment server written in Go, using SQLite for persistent storage.

The server provides a simple API for submitting and retrieving comments, alongside a static HTML frontend.

## Features

- Submit comments using HTTP PUT requests.
- Retrieve stored comments as JSON.
- SQLite-backed persistent storage.
- Configurable listening address, port, and database path.
- Built with Go's standard `net/http` package.
- SQLite driver powered by [`modernc.org/sqlite`](https://pkg.go.dev/modernc.org/sqlite).

## Requirements

- Go
- SQLite database file

## Installation

Clone the repository:

```sh
git clone https://github.com/<username>/<repository>.git
cd <repository>
```

Install dependencies:

```sh
go mod tidy
```

Build the server:

```sh
go build -o commenter .
```

## Usage

Start the server using the default configuration:

```sh
./commenter
```

The server listens on:

```text
http://localhost:8000
```

### Command-line Arguments

| Flag | Default | Description |
|------|---------|-------------|
| `-h` | `localhost` | Host or IP address to listen on |
| `-p` | `8000` | Listening port |
| `-d` | `./messages.sqlite` | SQLite database file path |

Example:

```sh
./commenter -h 0.0.0.0 -p 8080 -d ./data/messages.sqlite
```

The specified database directory must exist.

## API

### `PUT /comment`

Submit a new comment.

**Content-Type:** `application/x-www-form-urlencoded`

#### Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `username` | Yes | Comment author's username |
| `message` | Yes | Comment content |

#### Example

```sh
curl -X PUT http://localhost:8000/comment \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=guest&message=Hello%20world"
```

#### Response

```text
insert successful
```

HTTP Status: `201 Created`

### `GET /fetch`

Retrieve all stored comments.

#### Example

```sh
curl http://localhost:8000/fetch
```

#### Response

```json
[
    {
        "time": "2026-11-02T13:38:26.110186357+07:00"
        "username": "guest",
        "message": "Hello world"
    }
]
```

HTTP Status: `200 OK`

### `GET /`

Serves the static `index.html` frontend.

The HTML file must be present in the server's working directory.

## Database

Comments are stored in a SQLite database using the following schema:

```sql
CREATE TABLE messages (
    author TEXT,
    data TEXT
);
```

The database is expected to be initialized before the server starts.

## Project Structure

```text
.
├── main.go
├── index.html
├── messages.sqlite
├── go.mod
└── go.sum
```

## Missing
This comment server is missing:
- A HTTP body size limit, so now everyone can send a bomb to your server.
- Valid encoding, so now everyone can send garbled binaries and text.
- Literally every quality-of-life features.

## License

This project is licensed under the ISC License.

## FAQ

> **Was the README slopped?**
- Yes.

> **Ok... was the frontend slopped?**
- ~Yes.~
- NOT ANYMORE!!!

> **...Was the backend slopped?**
- Nah.

> **...Really?**
- Nah yeah!

> **What?**
- Use something else other than this 🤑.

