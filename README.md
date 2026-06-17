# Go & MySQL Todo API

A lightweight, containerized Todo REST API built with Go and MySQL. The entire application environment is managed using Docker Compose, allowing for seamless local development without needing to manually install Go or MySQL on your host machine.

## 🛠️ Architecture & Tech Stack

- **Backend:** Go (Golang)
- **Database:** MySQL 8.0
- **Containerization:** Docker & Docker Compose
- **Version Control:** Git (with atomic feature-driven commits)

## 🚀 Getting Started

### Prerequisites
Make sure you have [Docker](https://www.docker.com/products/docker-desktop/) and Docker Compose installed on your machine.

### Installation & Running
You can spin up the entire stack (the Go backend and the MySQL database) with a single command:

```bash
docker-compose up --build
```
#### The server will start up and listen for requests on http://localhost:8000