# Petal 🌼

Petal is a free, open-source progressive Ib app (PWA) designed for plant enthusiasts. With Petal, you can manage and track your plants, participate in plant-related discussions, and share images of your plants. My intent during developing this is just to keep it simple and fast. I hope you like it and find it useful. If you have any suggestions or feedback, please feel free to share it with me. I would love to hear from anyone, and welcome any contributions to this project via a branch and a PR.

## Features
- **User Authentication:** Secure login and registration using JWT tokens.
- **Dashboard:** Easily manage and track your plants.
- **Image Uploads:** Upload and view images of your plants.
- **Discussion Forum:** Join discussions, create posts (roots), and add comments.
- **Progressive Ib App (PWA):** Installable on devices for a native app-like experience.

## Tech Stack
### Frontend:
- HTML5
- Bootstrap 5
- HTMX Ajax Requests for ssr
- Some vanilla JS for frontend interactions

### Backend:
- Go
- Chi router
- SQLite

### Deployment:
- Docker
- Caddy

## Installation

### Prerequisites
- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Caddy](https://caddyserver.com/docs/install)

### Steps
1. **Clone the Repository:**
    ```sh
    git clone https://github.com/smithd36/petal.git
    cd petal
    ```

2. **Environment Variables:** Create a `.env` file in the project root with the following content:
    ```env
    JWT_KEY=your_secret_key
    ```

3. **Build and Run with Docker:**
    ```sh
    docker-compose up --build
    ```

4. **Access the App:** Open your browser and go to `http://localhost:8080`.

## Usage
### User Authentication
- **Register:** Go to `/register` to create a new account.
- **Login:** Go to `/login` to log into your account.

### Dashboard
- **Manage Plants:** Add and track your plants on the dashboard.
- **Upload Images:** Upload images of your plants to your personal gallery.

### Forum
- **Create Roots:** Create new discussion posts (roots) on the forum.
- **Add Comments:** Comment on existing posts.

## Contributing
I wecome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) for more details.

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Thank you for using Petal! I hope it helps you in your plant management journey. 🌱
