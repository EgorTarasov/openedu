# OpenEdu Project

OpenEdu is a simple web application built using Go and the Fiber framework. This project is designed to collect HTML data from iframes and send it to a server endpoint.

## Project Structure

```
openedu
├── cmd
│   └── main.go          # Entry point of the application
├── extension
│   └── sub.js           # JavaScript code for collecting iframe content
├── internal
│   ├── handlers
│   │   └── collect.go   # Handler for processing collected data
│   └── models
│       └── payload.go    # Data structure for incoming HTML data
├── go.mod                # Module definition and dependencies
├── go.sum                # Dependency checksums
└── README.md             # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd openedu
   ```

2. **Install Go dependencies:**
   ```
   go mod tidy
   ```

3. **Run the application:**
   ```
   go run cmd/main.go
   ```

## Usage

The application exposes a single route:

- **POST /collect**: This route receives HTML data from the `sub.js` script. The data should be sent in JSON format containing the HTML content and the URL.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.


## TODO:
- [ ] fix payload from extension to match one in server
- [ ] add metadata into payload such as title from a page or name of the test
- [ ] create pipeline for processing html into questions with answers and correct answers
- [ ] make list endpoint for all available tests and questions