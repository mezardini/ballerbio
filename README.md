## Introduction

ballerbio is a Go-based application designed to create and manage public profiles for amateur football players. This project addresses the need for a centralized platform where players can showcase their skills, stats, and contact information, making it easier for coaches, scouts, and other players to discover them.

ballerbio offers several key benefits: a simple and intuitive user interface for profile creation and editing, a robust data model for storing player information, and a public API for data retrieval and integration with other platforms. This allows for easy sharing and discovery of player profiles.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

*   Create and manage your public football player profile.
*   Customize your profile with personal information and football-related stats.
    *   Specify your position, preferred foot, and playing experience.
    *   Input key statistics such as goals, assists, and games played.
*   Share your profile with a unique, accessible URL.
*   View other players' profiles.
*   Search for players by name, position, or other criteria.

## Tech Stack

This project leverages the following technologies:

*   **Backend:**
    *   Go (version specified in `go.mod`)
    *   Gin (Go web framework)
*   **Database:**
    *   PostgreSQL (version unspecified)

## Prerequisites

To successfully build and run this application, ensure the following prerequisites are met:

**Required:**

*   **Go:** Version 1.20 or later. Download and install Go from the official website: [https://go.dev/dl/](https://go.dev/dl/)
*   **Git:** Version 2.28 or later. Ensure Git is installed and accessible in your system's PATH.
*   **Operating System:** A Unix-like operating system (Linux, macOS) or Windows with WSL2.

**Optional:**

*   **Database:** A PostgreSQL database instance. Configure the database connection string in the application's configuration file.
*   **Text Editor/IDE:** A code editor or IDE with Go support (e.g., VS Code with the Go extension, GoLand).

## Installation

To install and configure the ballerbio application, follow these steps:

1.  **Clone the Repository:** Clone the ballerbio repository to your local machine using the provided Git URL.

    ```bash
    git clone https://github.com/mezardini/ballerbio.git
    ```

2.  **Navigate to the Project Directory:** Change your current directory to the newly cloned `ballerbio` directory.

    ```bash
    cd ballerbio
    ```

3.  **Install Go Dependencies:** Utilize the Go package manager to install all required dependencies.

    ```bash
    go mod tidy
    ```

4.  **Database Setup:** Set up the database by running the database migrations. Ensure you have a PostgreSQL database instance running and the necessary credentials configured.

    ```bash
    go run ./cmd/migrate/main.go up
    ```

5.  **Environment Variable Configuration:** Configure the necessary environment variables. Create a `.env` file in the project root and populate it with the required values.  Example:

    ```
    DATABASE_URL=postgres://user:password@host:port/database_name
    PORT=8081
    ```

```markdown
## Usage

To run the application, execute the following command from the project's root directory:

```bash
go run main.go
```

This command compiles and runs the `main.go` file, starting the Gin web server. The server will listen on `localhost:8081`.

## Contributing

Thank you for your interest in contributing to `ballerbio`. Your contributions are highly valued. Please review the following guidelines before submitting any issues or pull requests.

## License

This project is not licensed. You are not granted any rights to use, copy, modify, or distribute this software.
