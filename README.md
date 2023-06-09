## Running the Program
To run the program, follow these steps:
1. Make sure you have Go installed (version 1.13 and above).
2. Clone the repository or copy the source code to a local directory.
3. Navigate to the root directory of the program.
4. Run the app using commands:
```shell
docker build -t promotions .
docker run -p 1312:1312 promotions
```
5. The server will be available at http://localhost:1312.

## Program Description
This program is an example server implemented in the Go programming language that provides an API for retrieving information about promotions.

# Functionality
The program provides the following functionality:
Retrieving information about a promotion based on its ID.
Automatic updating of promotion data from a CSV file every 30 minutes.
Returning data in JSON format.

# Implementation
The program is implemented using the following components and technologies:
Programming Language: Go.
Framework: Gorilla Mux for HTTP request routing.
Go standard packages for handling HTTP, JSON, CSV, and concurrency.
Use of a mutex (sync.RWMutex) to ensure safe access to promotion data during updates.

## API

The program provides the following API:

### GET /promotions/{id}

Returns information about a promotion based on the specified ID.

- Path Parameters:
  - `id` - The ID of the promotion.

- Responses:
  - 200 OK: Returns the promotion information in JSON format.
  - 404 Not Found: If a promotion with the specified ID is not found.

## Updating Data

Promotion data is automatically updated every 30 minutes from the `promotions.csv` CSV file. When updating the data, the mutex is locked, the data is read from the file, and the `promotionsMap` is updated with the promotion information. This ensures the freshness of the data when handling requests.


