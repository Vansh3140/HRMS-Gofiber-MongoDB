# Fiber-MongoDB CRUD Application

This is a simple CRUD (Create, Read, Update, Delete) application built using the [Fiber](https://gofiber.io/) web framework and [MongoDB](https://www.mongodb.com/) as the database.

## Features
- Retrieve all employees
- Add a new employee
- Update an existing employee by ID
- Delete an employee by ID

## Prerequisites

1. **Go**: Ensure you have Go installed on your system. [Download Go](https://golang.org/dl/).
2. **MongoDB**: Ensure you have MongoDB installed and running. [Install MongoDB](https://docs.mongodb.com/manual/installation/).
3. **Go Modules**: Initialize the project with Go modules if not already done.

## Installation

1. Clone the repository:
   ```sh
   git clone <repository_url>
   cd <repository_folder>
   ```

2. Install dependencies:
   ```sh
   go mod tidy
   ```

3. Start MongoDB server if not running:
   ```sh
   mongod
   ```

## Configuration

The application connects to MongoDB using the following default configuration:
- **Database Name**: `fiber-hrms`
- **Mongo URI**: `mongodb://localhost:27017/fiber-hrms`

Ensure that your MongoDB instance is running locally, or update the `mongoURI` constant in the code to point to your MongoDB instance.

## Running the Application

1. Start the application:
   ```sh
   go run main.go
   ```

2. The server will start on `http://localhost:3000`.

## API Endpoints

### 1. Get All Employees
**GET** `/api/v1/employee`

- Fetches all employees from the database.
- **Response**: JSON array of employee objects.

### 2. Add a New Employee
**POST** `/api/v1/employee`

- Adds a new employee to the database.
- **Request Body** (JSON):
  ```json
  {
    "name": "John Doe",
    "salary": 50000,
    "age": 30
  }
  ```
- **Response**: The created employee object.

### 3. Update an Employee by ID
**PUT** `/api/v1/employee/:id`

- Updates an existing employee by their ID.
- **Request Parameters**:
  - `id`: The ID of the employee to update.
- **Request Body** (JSON):
  ```json
  {
    "name": "Jane Doe",
    "salary": 60000,
    "age": 32
  }
  ```
- **Response**: The updated employee object.

### 4. Delete an Employee by ID
**DELETE** `/api/v1/employee/:id`

- Deletes an employee by their ID.
- **Request Parameters**:
  - `id`: The ID of the employee to delete.
- **Response**: A success message if the deletion is successful.

## Project Structure

- **main.go**: The main application file containing all routes and logic.

## Technologies Used
- [Go](https://golang.org/)
- [Fiber](https://gofiber.io/)
- [MongoDB](https://www.mongodb.com/)

## Error Handling
- The application uses HTTP status codes to indicate success or failure.
  - `200`: Success
  - `400`: Bad Request (e.g., invalid ID format)
  - `404`: Not Found (e.g., employee not found)
  - `500`: Internal Server Error

## License
This project is open-source and available under the [MIT License](LICENSE)
