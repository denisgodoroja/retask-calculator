# Pack Calculator Service

A full-stack web application designed to manage configurable pack sizes and calculate the optimal number of packs required to fulfill orders of a specific quantity.

## Features

* **Pack Size Management:** defined available pack sizes (e.g., 250, 500, 1000) via a web interface.

* **Order Calculator:** Input a required number of items to receive the optimal combination of packs.

* **Web Interface:** Clean, responsive UI built with Bootstrap and vanilla JavaScript.

* **Dockerized:** Complete setup with Nginx acting as a reverse proxy for the Go backend and static frontend.

## Tech Stack

* **Backend:** Go (Golang)

* **Frontend:** HTML5, JavaScript, Bootstrap 5

* **Proxy/Web Server:** Nginx

* **Infrastructure:** Docker, Docker Compose

* **Database:** MySQL (Required for persistence)

## Prerequisites

* [Docker](https://www.docker.com/get-started) and [Docker Compose](https://docs.docker.com/compose/install/) installed on your machine.

* (Optional for local dev) Go 1.25+ installed.

## Getting Started (Docker)

The easiest way to run the application is using Docker Compose. This will spin up the Go backend and the Nginx proxy.

1. **Clone the repository:**

   ```
   git clone <your-repo-url>
   cd <your-repo-directory>
   ```

2. **Build and run:**

   ```
   docker-compose up --build
   ```

3. **Access the application:**
   Open your browser and navigate to [http://localhost:8080](http://localhost:8080).

   * *Note: Port 8080 is defined in your `docker-compose.yml`.*

## Development

### Running Locally (without Docker)

If you prefer to run the Go application directly on your machine:

1. Ensure you have a MySQL database running and accessible.

2. Set necessary environment variables required by your `main.go` application:

   * `DB_HOST` - the database host name to connect to (e.g. `localhost`).

   * `DB_USER` - the username to connect with.

   * `DB_PASSWORD` - the password to connect with.

   * `DB_DATABASE` - the database name to connect to.

3. Run the Go server:

   ```
   go run .
   ```

   The Go app will default to port `8080`, but can be changed by setting the environment variable `PORT` to any other port number.

4. You can open `ui/index.html` directly in a browser, but you may need to adjust the `ENDPOINTS` in the `<script>` tag to point to `http://localhost:8080` instead of relative paths if not serving through Nginx.

### Unit Testing

To run the unit tests for the Go backend:

```
go test ./... -v
```

## API Reference

The application exposes the following RESTful endpoints (proxied via Nginx at standard paths):

### 1. Get Pack Sizes

Retrieves the currently configured pack sizes.

* **URL:** `/pack/get-sizes`

* **Method:** `GET`

* **Success Response:**

  ```
  {
    "sizes": [250, 500, 1000]
  }
  ```

### 2. Set Pack Sizes

Updates the list of available pack sizes.

* **URL:** `/pack/set-sizes`

* **Method:** `POST`

* **Body:**

  ```
  {
    "sizes": [250, 500, 1000]
  }
  ```

### 3. Calculate Order

Calculates the required packs for a given number of items.

* **URL:** `/order`

* **Method:** `POST`

* **Body:**

  ```
  {
    "amount": 123
  }
  ```

* **Success Response:**

  ```
  {
    "packs": {
      "250": 1,
      "500": 2
    }
  }
  ```
