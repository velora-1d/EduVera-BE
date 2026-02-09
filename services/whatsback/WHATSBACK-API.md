# Combined API Documentation

This document provides detailed information about various API endpoints for managing health checks, cron jobs, WhatsApp messaging, job scheduling, command management, contacts, and groups.

Each section below outlines the endpoints, methods, parameters, responses, and example requests.

---

## Table of Contents

- [Health Check](#health-check)
- [Cron Next Runs](#cron-next-runs)
- [WhatsApp Message API](#whatsapp-message-api)
  - [Send Message to User](#send-message-to-user)
  - [Send Message to Group](#send-message-to-group)
- [Job API](#job-api)
  - [Get Jobs](#get-jobs)
  - [Get Job by ID](#get-job-by-id)
  - [Get Jobs by Status](#get-jobs-by-status)
  - [Create Job](#create-job)
  - [Update Job](#update-job)
  - [Delete Job (Soft Delete)](#delete-job-soft-delete)
  - [Force Delete Job](#force-delete-job)
- [Command API](#command-api)
  - [Base URL](#command-api-base-url)
  - [Get All Commands](#get-all-commands)
  - [Create a Command](#create-a-command)
  - [Update a Command](#update-a-command)
  - [Delete a Command](#delete-a-command)
- [Contact API](#contact-api)
  - [Base URL](#contact-api-base-url)
  - [Get Paginated Contacts](#get-paginated-contacts)
- [Group API](#group-api)
  - [Base URL](#group-api-base-url)
  - [Get Paginated Groups](#get-paginated-groups)

---

## Health Check

- **Endpoint:** `GET /health`
- **Description:**  
  Checks the health of the API service.
- **Response:**
  - **Status Code:** `200 OK`
  - **Body:**  
    ```
    OK
    ```

---

## Cron Next Runs

- **Endpoint:** `POST /api/cron-next-runs`
- **Description:**  
  Processes a cron expression and returns its human-readable description along with the next run times.
- **Request Headers:**
  - `Content-Type: application/json`
- **Request Body:**
  - `exp` (string): A valid cron expression.
- **Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "description": "Every minute",
      "nexRuns": [
        "2023-09-22T12:00:00+07:00",
        "2023-09-22T12:01:00+07:00",
        "2023-09-22T12:02:00+07:00"
      ]
    }
    ```
- **Notes:**  
  The cron expression is processed using the timezone `Asia/Jakarta`.

---

## WhatsApp Message API

### Send Message to User

- **Endpoint:** `POST /api/message/send-message`
- **Description:**  
  Sends a direct message to a specified WhatsApp user.
- **Request Headers:**
  - `Content-Type: application/json`
- **Request Body:**
  - `number` (string): The recipient's phone number.
  - `message` (string): The text message to send.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "status": true,
      "message": "Message sent"
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - <error details>"
    }
    ```
- **Example Request:**
  ```bash
  curl -X POST http://localhost:5001/api/message/send-message \
       -H "Content-Type: application/json" \
       -d '{"number": "+628123456789", "message": "Hello from WhatsApp"}'
  ```

---

### Send Message to Group

- **Endpoint:** `POST /api/message/send-group-message`
- **Description:**  
  Sends a message to a specified WhatsApp group.
- **Request Headers:**
  - `Content-Type: application/json`
- **Request Body:**
  - `groupId` (string): The WhatsApp group ID (must end with `@g.us`).
  - `message` (string): The text message to send.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "status": true,
      "message": "Message sent to group"
    }
    ```
- **Validation Error Response:**  
  If the `groupId` does not end with `@g.us`:
  - **Status Code:** `422 Unprocessable Entity`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "Message sent to group"
    }
    ```
  *Note:* An error is logged indicating that group IDs must end with `@g.us`.
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - <error details>"
    }
    ```
- **Example Request:**
  ```bash
  curl -X POST http://localhost:5001/api/message/send-group-message \
       -H "Content-Type: application/json" \
       -d '{"groupId": "<group-id>@g.us", "message": "Hello Group"}'
  ```

---

## Job API

### Get Jobs

- **Endpoint:** `GET /api/jobs`
- **Description:**  
  Retrieves a paginated list of jobs with an optional search term.
- **Query Parameters:**
  - `search` (string, optional): Term to filter jobs.
  - `perPage` (number, optional): Number of jobs per page (default is 10).
  - `page` (number, optional): Page number (default is 1).
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "success": true,
      "data": [
        {
          "id": "1",
          "job_name": "Sample Job",
          "job_trigger": "send_message",
          "target_contact_or_group": "+628123456789",
          "message": "Hello",
          "job_cron_expression": "* * * * *"
        }
      ],
      "total": 1
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - <error details>"
    }
    ```

---

### Get Job by ID

- **Endpoint:** `GET /api/jobs/:id`
- **Description:**  
  Retrieves details of a specific job using its ID.
- **URL Parameter:**
  - `id` (string): The job identifier.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "success": true,
      "data": {
        "id": "1",
        "job_name": "Sample Job",
        "job_trigger": "send_message",
        "target_contact_or_group": "+628123456789",
        "message": "Hello",
        "job_cron_expression": "* * * * *"
      }
    }
    ```
- **Not Found Response:**
  - **Status Code:** `404 Not Found`
  - **Body Example:**
    ```json
    {
      "success": false,
      "message": "Job with ID 1 not found"
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - <error details>"
    }
    ```

---

### Get Jobs by Status

- **Endpoint:** `GET /api/jobs/status/:status`
- **Description:**  
  Retrieves jobs filtered by their status.
- **URL Parameter:**
  - `status` (number or string): The status of the jobs to retrieve (e.g., `0` for disabled, `1` for enabled).
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "success": true,
      "data": [
        {
          "id": "1",
          "job_name": "Sample Job",
          "job_trigger": "send_message",
          "target_contact_or_group": "+628123456789",
          "message": "Hello",
          "job_cron_expression": "* * * * *"
        }
      ]
    }
    ```
- **Not Found Response:**
  - **Status Code:** `404 Not Found`
  - **Body Example:**
    ```json
    {
      "success": false,
      "message": "Job with status 1 not found"
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - <error details>"
    }
    ```

---

### Create Job

- **Endpoint:** `POST /api/jobs`
- **Description:**  
  Creates a new job.
- **Request Headers:**
  - `Content-Type: application/json`
- **Request Body:**
  - `job_name` (string): The name of the job.
  - `job_trigger` (string): The trigger type (e.g., `"send_message"`).
  - `target_contact_or_group` (string): The target contact or group.
  - `message` (string): The message content.
  - `job_cron_expression` (string): The cron expression for scheduling the job.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "success": true,
      "data": {
        "id": "1",
        "job_name": "Sample Job",
        "job_trigger": "send_message",
        "target_contact_or_group": "+628123456789",
        "message": "Hello",
        "job_cron_expression": "* * * * *"
      }
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - <error details>"
    }
    ```

---

### Update Job

- **Endpoint:** `PUT /api/jobs/:id`
- **Description:**  
  Updates an existing job by its ID.
- **URL Parameter:**
  - `id` (string): The job identifier.
- **Request Headers:**
  - `Content-Type: application/json`
- **Request Body:**  
  An object containing the job fields to be updated.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "success": true,
      "data": {
        "id": "1",
        "job_name": "Updated Job",
        "job_trigger": "send_message",
        "target_contact_or_group": "+628123456789",
        "message": "Hello updated",
        "job_cron_expression": "* * * * *"
      }
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - <error details>"
    }
    ```

---

### Delete Job (Soft Delete)

- **Endpoint:** `DELETE /api/jobs/:id`
- **Description:**  
  Performs a soft delete of a job by its ID.
- **URL Parameter:**
  - `id` (string): The job identifier.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "success": true,
      "data": {
        "id": "1",
        "deleted": true
      }
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - <error details>"
    }
    ```

---

### Force Delete Job

- **Endpoint:** `DELETE /api/jobs/force/:id`
- **Description:**  
  Permanently deletes a job by its ID.
- **URL Parameter:**
  - `id` (string): The job identifier.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "success": true,
      "data": {
        "id": "1",
        "deleted": true
      }
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - <error details>"
    }
    ```

---

## Command API

This API allows you to manage commands and their associated responses.

### Command API Base URL

- **Base URL:**  
  `http://[your-host]/api/command`

---

### Get All Commands

- **Method:** `GET`
- **URL:** `/`
- **Description:**  
  Retrieve a list of all available commands.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "status": true,
      "data": {
        "commands": [
          { "command": "!greet", "response": "Hello!" },
          { "command": "!help", "response": "How can I assist?" }
        ]
      }
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - [error details]"
    }
    ```

---

### Create a Command

- **Method:** `POST`
- **URL:** `/`
- **Description:**  
  Create a new command.
- **Request Headers:**
  - `Content-Type: application/json`
- **Request Body (JSON):**
  ```json
  {
    "command": "string",
    "response": "string"
  }
  ```
- **Validation:**  
  - `command` and `response` must be strings.
- **Responses:**
  - **Status Code:** `200 OK`  
    Command saved successfully.
    ```json
    {
      "status": true,
      "message": "Command saved successfully"
    }
    ```
  - **Status Code:** `400 Bad Request`  
    Invalid request body (e.g., missing fields or invalid types).
  - **Status Code:** `500 Internal Server Error`
    ```json
    {
      "status": false,
      "message": "API Error - [error details]"
    }
    ```
- **Example Request:**
  ```bash
  curl -X POST http://localhost:3000/api/command \
    -H "Content-Type: application/json" \
    -d '{"command": "!greet", "response": "Hello!"}'
  ```

---

### Update a Command

- **Method:** `PUT`
- **URL:** `/:command_id`
- **Description:**  
  Update an existing command by ID.
- **Path Parameter:**
  - `command_id` (integer): ID of the command to update.
- **Request Headers:**
  - `Content-Type: application/json`
- **Request Body (JSON):**
  ```json
  {
    "command": "string",
    "response": "string"
  }
  ```
- **Validation:**  
  Same as for creating a command.
- **Responses:**
  - **Status Code:** `200 OK`  
    Command updated successfully.
    ```json
    {
      "status": true,
      "message": "Command updated successfully"
    }
    ```
  - **Status Code:** `400 Bad Request`  
    Invalid request body.
  - **Status Code:** `500 Internal Server Error`
    ```json
    {
      "status": false,
      "message": "API Error - [error details]"
    }
    ```
- **Example Request:**
  ```bash
  curl -X PUT http://localhost:3000/api/command/1 \
    -H "Content-Type: application/json" \
    -d '{"command": "!hello", "response": "Hi there!"}'
  ```

---

### Delete a Command

- **Method:** `DELETE`
- **URL:** `/:command_id`
- **Description:**  
  Delete a command by ID.
- **Path Parameter:**
  - `command_id` (integer): ID of the command to delete.
- **Responses:**
  - **Status Code:** `200 OK`  
    Command deleted successfully.
    ```json
    {
      "status": true,
      "message": "Command deleted successfully"
    }
    ```
  - **Status Code:** `500 Internal Server Error`
    ```json
    {
      "status": false,
      "message": "API Error - [error details]"
    }
    ```
- **Example Request:**
  ```bash
  curl -X DELETE http://localhost:3000/api/command/1
  ```

---

## Contact API

This API allows you to fetch paginated contact data with optional search capabilities.

### Contact API Base URL

- **Base URL:**  
  `http://[your-host]/api/contacts`

---

### Get Paginated Contacts

- **Method:** `GET`
- **URL:** `/`
- **Description:**  
  Retrieve a paginated list of contacts with optional search filtering.
- **Query Parameters:**
  - `search` (string, optional): Search term to filter contacts by name. Default: `""`.
  - `perPage` (number, optional): Number of contacts per page. Default: `10`.
  - `page` (number, optional): Page number to fetch. Default: `1`.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "status": true,
      "message": "Contacts fetched successfully",
      "data": {
        "contacts": [
          { "id": 1, "name": "John Doe", "email": "john@example.com" },
          { "id": 2, "name": "Jane Smith", "email": "jane@example.com" }
        ],
        "totalContacts": 100
      }
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "Internal Server Error"
    }
    ```
- **Example Request:**
  ```bash
  curl -X GET "http://localhost:3000/api/contacts?search=John&perPage=5&page=2"
  ```

---

## Group API

This API allows you to fetch paginated group data with optional search capabilities.

### Group API Base URL

- **Base URL:**  
  `http://[your-host]/api/groups`

---

### Get Paginated Groups

- **Method:** `GET`
- **URL:** `/`
- **Description:**  
  Retrieve a paginated list of groups with optional search filtering.
- **Query Parameters:**
  - `search` (string, optional): Search term to filter groups by name. Default: `""`.
  - `perPage` (number, optional): Number of groups per page. Default: `10`.
  - `page` (number, optional): Page number to fetch. Default: `1`.
- **Success Response:**
  - **Status Code:** `200 OK`
  - **Body Example:**
    ```json
    {
      "status": true,
      "message": "Fetch group page 1",
      "data": [
        { "id": 1, "name": "Admins", "description": "Administrators group" },
        { "id": 2, "name": "Users", "description": "Regular users group" }
      ]
    }
    ```
- **Error Response:**
  - **Status Code:** `500 Internal Server Error`
  - **Body Example:**
    ```json
    {
      "status": false,
      "message": "API Error - [error details]"
    }
    ```
- **Example Request:**
  ```bash
  curl -X GET "http://localhost:3000/api/groups?search=Admins&perPage=5&page=1"
  ```

---

## Common Error Response Structure

For all endpoints, error responses follow this structure:

```json
{
  "status": false,
  "message": "Error description"
}
```
