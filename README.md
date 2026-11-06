# commenter 

A lightweight HTTP comment server written in Go, using SQLite for persistent storage.

The server provides a simple API for submitting and retrieving comments, alongside a static HTML frontend.

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

## API Documentation

Base URL: `http://localhost:8000`

All API request and response bodies use JSON unless otherwise specified.

### Endpoints

#### `POST /comment`

Create a new comment.

**Request body:**

```json
{
  "username": "user",
  "message": "Hello, world!"
}
```

**Responses:**

| Status | Description |
|---|---|
| `201 Created` | Comment successfully created |
| `400 Bad Request` | Invalid JSON or missing username/message |
| `405 Method Not Allowed` | Request method is not POST |
| `413 Content Too Large` | Request body exceeds the configured size limit |
| `500 Internal Server Error` | Database insertion failed |

The server automatically assigns an ID and Unix timestamp.

#### `GET /fetch`

Retrieve all stored comments.

**Response body:**

```json
[
  {
    "id": 1,
    "time": 1793935201,
    "username": "user",
    "message": "Hello, world!"
  }
]
```

**Responses:**

| Status | Description |
|---|---|
| `200 OK` | Comments successfully retrieved |
| `405 Method Not Allowed` | Request method is not GET |
| `500 Internal Server Error` | Database retrieval failed |

Returns an empty JSON array when no comments exist.

#### `DELETE /delete`

Delete a comment by its ID.

**Request body:**

```json
{
  "id": 1
}
```

**Responses:**

| Status | Description |
|---|---|
| `204 No Content` | Comment successfully deleted |
| `400 Bad Request` | Invalid JSON |
| `404 Not Found` | Comment does not exist or affected-row lookup failed |
| `405 Method Not Allowed` | Request method is not DELETE |
| `413 Content Too Large` | Request body exceeds the configured size limit |
| `500 Internal Server Error` | Database deletion failed |

## Database

Comments are stored in a SQLite database using the following schema:

```sql
CREATE TABLE messages (
    time TIMESTAMP,
    author TEXT,
    data TEXT
);
```

The database is expected to be initialized before the server starts.

## FAQ

> **Was the README slopped?**
- Yes.

> **Ok... was the example frontend slopped?**
- ~Yes.~
- NOT ANYMORE!!!

> **...Was the backend slopped?**
- Nah.

> **...Really?**
- Nah yeah!

> **What?**
- Use something else other than this 🤑.

## License

This project is licensed under the ISC License.
