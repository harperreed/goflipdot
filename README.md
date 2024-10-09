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

## Testing

The `goflipdot` project includes a comprehensive test suite located in the `test/` directory. These tests cover various aspects of the library's functionality, including controller operations, packet handling, and sign management.

### Running Tests

To run the entire test suite, use the following command:

```sh
go test ./...
```

For more verbose output, add the `-v` flag:

```sh
go test -v ./...
```

### Test Structure

The test suite is organized into three main files:

1. `test/controller_test.go`: Tests for controller functionality
   - Creating a new controller
   - Adding signs
   - Starting and stopping test signs
   - Drawing images
   - Error handling for network issues

2. `test/packet_test.go`: Tests for packet handling
   - TestSignsStartPacket and TestSignsStopPacket
   - ImagePacket with various image sizes and patterns
   - Checksum calculation and verification

3. `test/sign_test.go`: Tests for sign-related functionality
   - Creating signs with different configurations
   - Creating and validating images
   - Image flipping

### Test Coverage

To view the test coverage, run:

```sh
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

This will generate an HTML report of the test coverage.

### Hardware Testing

While the provided tests use mocks and simulations, it's crucial to test the library with actual Hanover flipdot hardware to ensure proper functionality. When testing with real hardware:

1. Ensure proper serial port configuration.
2. Verify that packets are being sent in the correct format.
3. Check for any responses from the display and handle them appropriately.
4. Test various display sizes and configurations to ensure compatibility.

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
    - Comprehensive unit tests for the controller logic, including error scenarios.
- **test/packet_test.go**
    - Extensive unit tests for packet structures and behaviors, including various image sizes and patterns.
- **test/sign_test.go**
    - Thorough unit tests for sign-related logic, covering different configurations and image manipulations.

We hope this README has provided you with the needed information to get started with our project and understand our testing procedures. Happy coding! 💻
