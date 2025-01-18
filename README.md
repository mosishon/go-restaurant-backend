# Go Shopping Project

This is a simple shopping project built with Go. It is designed to manage food menus, orders, and invoices using a MongoDB database and the Gin web framework.

## Features

- Manage food items
- Manage menus
- Handle orders
- Generate invoices
- User authentication

## Technologies Used

- **Go**: The main programming language used for this project.
- **Gin**: A web framework for building APIs.
- **MongoDB**: A NoSQL database for storing data.

## Getting Started

### Prerequisites

- Go 1.22.11 or later
- MongoDB

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/go-shopping-project.git
    cd go-shopping-project
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Set up MongoDB:
    - Ensure MongoDB is running on `mongodb://127.0.0.1:27017/rest`.

### Running the Application

1. Start the application:
    ```sh
    go run main.go
    ```

2. The application will run on port `8001` by default. You can access it at `http://localhost:8001`.

## License

This project is licensed under the MIT License.