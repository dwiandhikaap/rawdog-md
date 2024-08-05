---
Template: docs
Title: API Documentation
Section:
    Name: Development
    Index: 2
---

# API Documentation

The TODO App provides a RESTful API that allows clients to interact with the application. Below are the details of the available endpoints.

## Base URL

```
http://localhost:3000/api
```

## Authentication

Some endpoints require authentication. Use the following method to obtain a token:

- **Login**: 
  - **Endpoint**: `/auth/login`
  - **Method**: POST
  - **Request Body**:
    ```json
    {
      "username": "your_username",
      "password": "your_password"
    }
    ```
  - **Response**:
    ```json
    {
      "token": "your_jwt_token"
    }
    ```

## Endpoints

### Tasks

#### Get All Tasks

- **Endpoint**: `/tasks`
- **Method**: GET
- **Headers**: 
  - `Authorization: Bearer your_jwt_token`
- **Response**:
  ```json
  [
    {
      "id": "1",
      "title": "Task 1",
      "description": "Description for Task 1",
      "dueDate": "2023-10-01",
      "completed": false
    },
    ...
  ]
  ```

#### Create a New Task

- **Endpoint**: `/tasks`
- **Method**: POST
- **Headers**: 
  - `Authorization: Bearer your_jwt_token`
- **Request Body**:
  ```json
  {
    "title": "New Task",
    "description": "Task description",
    "dueDate": "2023-10-01"
  }
  ```
- **Response**:
  ```json
  {
    "id": "2",
    "title": "New Task",
    "description": "Task description",
    "dueDate": "2023-10-01",
    "completed": false
  }
  ```

#### Update a Task

- **Endpoint**: `/tasks/:id`
- **Method**: PUT
- **Headers**: 
  - `Authorization: Bearer your_jwt_token`
- **Request Body**:
  ```json
  {
    "title": "Updated Task",
    "description": "Updated description",
    "completed": true
  }
  ```
- **Response**:
  ```json
  {
    "id": "2",
    "title": "Updated Task",
    "description": "Updated description",
    "dueDate": "2023-10-01",
    "completed": true
  }
  ```

#### Delete a Task

- **Endpoint**: `/tasks/:id`
- **Method**: DELETE
- **Headers**: 
  - `Authorization: Bearer your_jwt_token`
- **Response**:
  ```json
  {
    "message": "Task deleted successfully"
  }
  ```

## Conclusion

This API documentation provides the essential endpoints and their usage for interacting with the TODO App. For further details, please refer to the codebase or contact the development team.