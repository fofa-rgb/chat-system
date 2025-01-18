# Chat System

This repository contains the code for the Chat System API. The application is containerized using Docker Compose.

## Prerequisites
- Docker installed on your machine.
- Docker Compose installed.

## Setup Instructions

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd chat-system
   ```

2. **Set up the `.env` file**
   - Copy the contents of the `env-values.txt` file and paste it into a new file named `.env` in the root of the project directory.
   - Example:
     ```bash
     cp .env-values.txt .env
     ```

3. **Start the application**
   - Run the following command to build and start the application:
     ```bash
     docker-compose up
     ```

   - This will build and run all the required services as defined in the `docker-compose.yml` file.

4. **Access Logs**
   - To view logs for the application, use the following command:
     ```bash
     docker logs -f chat-system-container
     ```

5. **Stopping the Application**
   - To stop the running containers, press `Ctrl+C` or use the command:
     ```bash
     docker-compose down
     ```

## Environment Variables
The `.env` file must include the variables in the `.env-values.txt` file for the application to run:

## Troubleshooting
- If the application fails to start, check the logs using:
  ```bash
  docker logs -f chat-system-container
  ```

- Ensure that Docker and Docker Compose are installed and running.

- Verify that the `.env` file exists and contains valid values.

## Database Schema
The database schema can be found in the following file:  
`migrations/V1__create_tables.sql`

## API Controllers (Handlers)
The API controllers are implemented in the following files:  
- `api/handlers/applications.go`
- `api/handlers/chats.go`
- `api/handlers/messages.go`

## Workers and Tasks
The workers and tasks (cron jobs) are implemented in:  
`api/cron/cron.go`

## Routes and Parameters
The routes are defined in:  
`cmd/main/main.go`

### Path Parameters
Most user inputs are included as path parameters, visible directly in the routes. Below are specific details for additional required inputs:

#### Additional Parameters
1. **POST `/applications`**  
   - **Body**: `{"name": "string"}`
2. **POST `/applications/:token/chats`**  
   - **Body**: `{"subject": "string"}`
3. **POST `/applications/:token/chats/:chat_number/messages`**  
   - **Body**: `{"body": "string"}`
4. **GET `/applications/:token/chats/:chat_number/messages/search`**  
   - **Query**: `{"query": "string"}`

The structure of requests and responses is detailed in:  
`api/handlers/requestResponseStructure.go`

## Additional Notes
- This setup is designed for local development. For production deployment, additional configurations may be required.

