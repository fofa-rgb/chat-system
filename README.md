Chat System

This repository contains the code for the Chat System API. The application is containerized using Docker Compose, which makes it easy to set up and run.

Prerequisites

Docker installed on your machine.

Docker Compose installed.

Setup Instructions

Clone the repository

git clone <repository-url>
cd chat-system

Set up the .env file

Copy the contents of the env-values.txt file and paste it into a new file named .env in the root of the project directory.

Example:

cp .env-values.txt .env

Start the application

Run the following command to build and start the application:

docker-compose up

This will build and run all the required services as defined in the docker-compose.yml file.

Access Logs

To view logs for the application, use the following command:

docker logs -f chat-system-container

Stopping the Application

To stop the running containers, press Ctrl+C or use the command:

docker-compose down

Environment Variables

The .env file must include the variables in .env-values.txt

Troubleshooting

If the application fails to start, check the logs using:

docker logs -f chat-system-container

Ensure that Docker and Docker Compose are installed and running.

Verify that the .env file exists and contains valid values.

Additional Notes

This setup is designed for local development. For production deployment, additional configurations may be required.

