# goflipdot 🚀

Welcome to the **goflipdot** repository! This project is designed to manage and control Hanover flipdot displays via a Go-based interface.

## Summary of Project 📜

The **goflipdot** repository provides an easy-to-use interface for connecting, configuring, and controlling Hanover flipdot signs. It includes:

- A command-line example to run and test the flipdot display.
- Internal packages to handle different aspects of the flipdot communication and control.
- Unit tests to ensure the reliability and correctness of the code.

## History

This is a port of `pyflipdot` to Go. The original `pyflipdot` project can be found [here](https://github.com/briggySmalls/pyflipdot).

`pyflipdot` was built based off of John Whittington's [blog post](https://engineer.john-whittington.co.uk/2017/11/adventures-flippy-flip-dot-display/) and his node.js driver [node-flipdot](https://github.com/tuna-f1sh/node-flipdot)

## How to Use 📚

### Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/harperreed/goflipdot.git
    cd goflipdot
    ```

2. Build the project:
    ```sh
    make build
    ```

3. (Optional) Run tests to make sure everything is working:
    ```sh
    make test
    ```

### Running the Example

To run the example command-line application:
```sh
make run-example
```

This will start the test sequence on the connected flipdot sign and draw a checkerboard pattern.

## Tech Info ⚙️

- This project is written in Go, so make sure you have [Go installed](https://golang.org/doc/install).
- The code is organized into several packages:
  - `cmd`: Contains the example command-line application.
  - `internal`: Houses the internal logic for controllers, packets, and signs.
  - `pkg`: Provides the main public interface for controlling flipdot signs.
  - `test`: Includes unit tests for the different components.

### Directory/File Tree

```
goflipdot/
├── Makefile
├── cmd
│   └── example
│       └── main.go
├── go.mod
├── internal
│   ├── controller
│   │   └── controller.go
│   ├── packet
│   │   └── packet.go
│   └── sign
│       └── sign.go
├── pkg
│   └── goflipdot
│       └── goflipdot.go
└── test
    ├── controller_test.go
    ├── packet_test.go
    └── sign_test.go
```

### File Content Overview

- **Makefile**
    - Provides basic commands for building, testing, cleaning, and formatting the project.
- **cmd/example/main.go**
    - Example command-line application to interact with the flipdot signs.
- **go.mod**
    - The Go module file for dependency management.
- **internal/controller/controller.go**
    - Contains the logic for managing the communications with the flipdot signs.
- **internal/packet/packet.go**
    - Defines packets used to communicate with the flipdot signs.
- **internal/sign/sign.go**
    - Includes the logic for creating and manipulating flipdot sign properties.
- **pkg/goflipdot/goflipdot.go**
    - Main public interface providing higher-level functions for use by other applications.
- **test/controller_test.go**
    - Unit tests for the controller logic.
- **test/packet_test.go**
    - Unit tests for the packet structures and behaviors.
- **test/sign_test.go**
    - Unit tests for the sign-related logic.

We hope this README has provided you with the needed information to get started with our project. Happy coding! 💻
