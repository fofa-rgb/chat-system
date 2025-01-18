# Chat System

This repository contains the code for the Chat System API. The application is containerized using Docker Compose, which makes it easy to set up and run.

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

## Additional Notes
- This setup is designed for local development. For production deployment, additional configurations may be required.

