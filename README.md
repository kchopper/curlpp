# 🚀 curlpp – A Modern Alternative to `curl`, Built with Go

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/github/license/kchopper/curlpp)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/kchopper/curlpp?style=social)](https://github.com/kchopper/curlpp/stargazers)

**curlpp** is a fast, user-friendly, and modern alternative to `curl`, designed to simplify HTTP requests while offering powerful features like automatic retries, response formatting, and parallel requests—all in a standalone Go binary.

---

## 🚀 Features

- **Simple & Intuitive** – No more overwhelming flags! Sensible defaults make requests easy.
- **Beautiful Output** – JSON pretty-printing, syntax highlighting, and HTML parsing.
- **Parallel Requests** – Fetch multiple URLs simultaneously with `-p`.
- **Retry Mechanism** – Automatic retries with exponential backoff.
- **Authentication Built-In** – Supports API keys, OAuth, and JWT.
- **Flexible Response Handling** – Extract specific data using CSS selectors or JSONPath.

---

## 📦 Installation

~~~sh
go install github.com/kchopper/curlpp@latest
~~~

Or download a prebuilt binary from the [releases](https://github.com/kchopper/curlpp/releases) page.

---

## ⚡ Quick Start

### 1️⃣ Basic GET Request
~~~sh
curlpp https://jsonplaceholder.typicode.com/todos/1
~~~

### 2️⃣ POST Request with JSON Body
~~~sh
curlpp -X POST https://api.example.com/data -d '{"name": "curlpp"}' -H "Content-Type: application/json"
~~~

### 3️⃣ Parallel Requests
~~~sh
curlpp -p https://site1.com https://site2.com
~~~

### 4️⃣ Extracting HTML Data
~~~sh
curlpp https://news.ycombinator.com --selector ".title a"
~~~

### 5️⃣ Automatic Retries
~~~sh
curlpp --retry 3 https://unstable.api.com
~~~

---

## ⚙️ Configuration

Create a `~/.curlpprc` file to store default headers, auth tokens, and profiles.

~~~json
{
  "default_headers": {
    "User-Agent": "curlpp/1.0"
  },
  "profiles": {
    "github": {
      "base_url": "https://api.github.com",
      "headers": {
        "Authorization": "Bearer YOUR_TOKEN"
      }
    }
  }
}
~~~

Use it like this:
~~~sh
curlpp --profile github /users/kchopper
~~~

---

## 🛠️ Roadmap

- [ ] **WebSocket Support**
- [ ] **GraphQL Queries**
- [ ] **Plugin System**

Have a feature request? Open an [issue](https://github.com/kchopper/curlpp/issues) or contribute!

---

## 🏆 Contributing

1. **Fork the repo**
2. **Create a feature branch**
3. **Commit and push your changes**
4. **Open a PR**

Check out our [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

## 📜 License

Licensed under the [MIT License](LICENSE).

---

💙 **Built with Go. Inspired by `curl`. Designed for modern developers.**
