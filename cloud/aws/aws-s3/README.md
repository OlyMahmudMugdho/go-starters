# AWS S3-like API

This project is a Go-based RESTful API that provides AWS S3-like functionality for managing files and buckets. It allows users to upload, download, and delete files, as well as create, list, update, and delete buckets. The API is built using the AWS SDK for Go (v2) and Gorilla Mux, with a Bash automation script (`s3_operations.sh`) to interact with the API using `curl` and `jq`.

## Features

- **File Operations**:
  - Upload files to a specified bucket (`POST /s3/upload`).
  - Download files from a bucket (`GET /s3/download/{bucket}/{key}`).
  - Delete files from a bucket (`DELETE /s3/delete/{bucket}/{key}`).
- **Bucket Operations**:
  - Create a new bucket (`POST /s3/create`).
  - List all buckets (`GET /s3/list`).
  - Update a bucket's name and region (`PUT /s3/update`).
  - Delete a bucket (`DELETE /s3/delete/{name}`).
- **Automation Script**: A Bash script (`s3_operations.sh`) automates all API operations, with error handling and JSON response parsing.
- **RESTful Design**: Uses standard HTTP methods and JSON/multipart form data for requests and responses.

## Project Structure

```
├── bin
│   ├── random.txt
│   ├── testf.txt
│   └── test.txt
├── cmd
│   └── server
│       └── main.go        # Entry point for the API server
├── config
│   └── config.go         # Configuration for AWS SDK
├── go.mod
├── go.sum
├── internal
│   ├── aws
│   │   └── s3.go         # AWS S3 client setup
│   ├── handler
│   │   ├── s3_file_handler.go  # Handlers for file operations
│   │   └── s3_handler.go       # Handlers for bucket operations
│   ├── routes
│   │   └── router.go     # API route definitions
│   └── service
│       ├── s3_file.go    # Service logic for file operations
│       └── s3_service.go # Service logic for bucket operations
├── s3_operations.sh      # Bash automation script
└── README.md             # Project documentation
```

## Prerequisites

- **Go**: Version 1.16 or higher.
- **AWS SDK for Go v2**: Configured with valid AWS credentials.
- **curl**: For making HTTP requests in the Bash script.
- **jq**: For parsing JSON responses in the Bash script.
- **AWS Account**: Valid AWS credentials for S3 access (set in `~/.aws/credentials` or environment variables).

## Setup

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/OlyMahmudMugdho/go-starters.git
   cd go-starters/cloud/aws/aws-s3
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Configure AWS Credentials**:
   Ensure your AWS credentials are set up in `~/.aws/credentials` or via environment variables:
   ```bash
   export AWS_ACCESS_KEY_ID=<your-access-key>
   export AWS_SECRET_ACCESS_KEY=<your-secret-key>
   export AWS_REGION=eu-north-1
   ```

4. **Install `curl` and `jq`**:
   On Ubuntu:
   ```bash
   sudo apt-get update
   sudo apt-get install curl jq
   ```
   On macOS:
   ```bash
   brew install curl jq
   ```

## Running the API

1. **Start the Server**:
   ```bash
   go run cmd/server/main.go
   ```
   The API will be available at `http://localhost:8080`.

2. **Verify Server**:
   Test connectivity:
   ```bash
   curl http://localhost:8080/s3/list
   ```

## Using the Automation Script

The `s3_operations.sh` script automates all API operations. It creates a test file, performs file and bucket operations, and cleans up afterward.

1. **Make the Script Executable**:
   ```bash
   chmod +x s3_operations.sh
   ```

2. **Run the Script**:
   Ensure the server is running, then execute:
   ```bash
   ./s3_operations.sh
   ```

3. **Script Operations**:
   - Creates a bucket (`my-demo-bucket-<date>`).
   - Uploads `test.txt` to the bucket.
   - Downloads the file to `downloaded_test.txt`.
   - Deletes the file from the bucket.
   - Lists all buckets.
   - Updates the bucket name to `new-bucket-<date>`.
   - Deletes the bucket.
   - Cleans up `test.txt`.

4. **Example Output**:
   ```
   Checking if server is reachable at http://localhost:8080...
   Success: Server is reachable.
   Creating bucket: my-demo-bucket-2025-05-18
   Success: Create Bucket - HTTP 201
   {
     "message": "Bucket created successfully"
   }
   ...
   All operations completed successfully!
   ```

## API Endpoints

| Method | Endpoint                        | Description                     | Request Body/Params                     | Response                     |
|--------|---------------------------------|-------------------------------|-----------------------------------------|------------------------------|
| POST   | `/s3/create`                   | Create a bucket               | JSON: `{"name":"string","region":"string"}` | 201: `{"message":"Bucket created successfully"}` |
| GET    | `/s3/list`                     | List all buckets              | None                                    | 200: `["bucket1","bucket2"]` |
| DELETE | `/s3/delete/{name}`            | Delete a bucket               | URL param: `name`                       | 200: `{"message":"Bucket deleted successfully"}` |
| PUT    | `/s3/update`                   | Update a bucket               | JSON: `{"old_name":"string","new_name":"string","new_region":"string"}` | 200: `{"message":"Bucket updated successfully"}` |
| POST   | `/s3/upload`                   | Upload a file                 | Multipart: `bucket`, `key`, `file`      | 201: `File uploaded successfully` |
| GET    | `/s3/download/{bucket}/{key}`  | Download a file               | URL params: `bucket`, `key`             | 200: File stream or 404: `{"error":"..."}` |
| DELETE | `/s3/delete/{bucket}/{key}`    | Delete a file                 | URL params: `bucket`, `key`             | 200: `{"message":"File '...' deleted from bucket '...'"}` |

## Troubleshooting

- **HTTP 000 Error**:
  If the script reports `Failed: ... - No response from server (HTTP 000)`:
  1. Ensure the server is running: `go run cmd/server/main.go`.
  2. Verify the port (default: 8080) and update `BASE_URL` in `s3_operations.sh` if needed.
  3. Test connectivity: `curl -v http://localhost:8080/s3/list`.
  4. Check for firewall issues: `sudo ufw allow 8080`.

- **AWS Credential Errors**:
  Ensure AWS credentials are valid and the region matches `eu-north-1` (or update `REGION` in the script).

- **File Not Found**:
  The script creates `test.txt` automatically. If upload fails, ensure the file is readable.

## Contributing

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature-name`.
3. Commit changes: `git commit -m "Add feature"`.
4. Push to the branch: `git push origin feature-name`.
5. Open a pull request.
