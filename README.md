# Setting Up PostgreSQL with Docker

This guide explains how to set up a PostgreSQL database using Docker. If you don't already have PostgreSQL installed locally, this is a simple and efficient way to get started.

## Prerequisites

- [Docker](https://www.docker.com/get-started) installed on your system.
- A `.env` file in the root directory to configure database credentials.

## Steps to Set Up PostgreSQL with Docker

1. **Create a `.env` file:**
   In the root directory of your project, create a file named `.env` with the following content:

   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=mercadolibre
   ```

   Replace `your_username`, `your_password`, and `your_database_name` with your desired PostgreSQL credentials.

2. **Start the PostgreSQL container:**
   Run the following command to start the PostgreSQL container:

   ```bash
   docker-compose up -d
   ```

3. **Verify the setup:**

   - Use the `docker ps` command to confirm that the PostgreSQL container is running.
   - Connect to the database using your favorite PostgreSQL client or the `psql` command-line tool:

     ```bash
     psql -h localhost -U your_username -d your_database_name
     ```

     Use the username and password defined in your `.env` file.

## Notes

- Ensure the `.env` file is included in your `.gitignore` file to avoid committing sensitive credentials.
- Modify the port in the `.env` file if `5432` is already in use on your machine.

That's it! Your PostgreSQL database should now be up and running with Docker.
