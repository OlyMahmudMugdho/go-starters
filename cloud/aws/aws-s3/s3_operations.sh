#!/bin/bash

# Configuration
BASE_URL="http://localhost:8080"
BUCKET="my-demo-bucket-$(date +%Y-%m-%d)"
FILE_KEY="test.txt"
FILE_PATH="./test.txt" # Local file for upload
DOWNLOAD_PATH="./downloaded_test.txt" # Path to save downloaded file
NEW_BUCKET="new-bucket-$(date +%Y-%m-%d)"
REGION="eu-north-1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check for curl and jq
command -v curl >/dev/null 2>&1 || { echo -e "${RED}curl is required but not installed.${NC}"; exit 1; }
command -v jq >/dev/null 2>&1 || { echo -e "${RED}jq is required but not installed.${NC}"; exit 1; }

# create the file
echo "Hello, S3!" > test.txt

# Function to check if server is reachable
check_server() {
    echo -e "${YELLOW}Checking if server is reachable at $BASE_URL...${NC}"
    response=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/s3/list" -o /dev/null)
    if [ "$response" -eq 000 ]; then
        echo -e "${RED}Error: Cannot connect to $BASE_URL. Is the server running?${NC}"
        echo -e "${YELLOW}Troubleshooting steps:${NC}"
        echo "1. Ensure the Go server is running: 'go run cmd/server/main.go'"
        echo "2. Verify the server is listening on port 8080."
        echo "3. Check if $BASE_URL is correct (modify BASE_URL in the script if needed)."
        echo "4. Test connectivity: 'curl -v $BASE_URL/s3/list'"
        exit 1
    fi
    echo -e "${GREEN}Server is reachable.${NC}"
}

# Function to check HTTP status code and parse response
check_status() {
    local status=$1
    local body=$2
    local operation=$3
    if [ "$status" -eq 000 ]; then
        echo -e "${RED}Failed: $operation - No response from server (HTTP 000)${NC}"
        echo -e "${YELLOW}Possible causes:${NC}"
        echo "1. Server is not running."
        echo "2. Incorrect BASE_URL ($BASE_URL)."
        echo "3. Network issue."
        echo -e "${YELLOW}Raw response:${NC} $body"
        exit 1
    elif [ "$status" -eq 200 ] || [ "$status" -eq 201 ]; then
        echo -e "${GREEN}Success: $operation - HTTP $status${NC}"
        if [[ "$body" == *"{("* || "$body" == *"{"* || "$body" == *"["* ]]; then
            echo "$body" | jq . 2>/dev/null || echo -e "${YELLOW}Warning: Response is not valid JSON:${NC} $body"
        else
            echo "$body"
        fi
    else
        echo -e "${RED}Failed: $operation - HTTP $status${NC}"
        echo "$body" | jq . 2>/dev/null || echo -e "${YELLOW}Error: Response is not valid JSON:${NC} $body"
        exit 1
    fi
}

# Check server connectivity
check_server

# Create Bucket
echo "Creating bucket: $BUCKET"
response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/s3/create" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"$BUCKET\",\"region\":\"$REGION\"}")
status=${response: -3}
body=${response%???}
check_status "$status" "$body" "Create Bucket"

# Upload File
echo "Uploading file: $FILE_PATH to bucket: $BUCKET with key: $FILE_KEY"
if [ ! -f "$FILE_PATH" ]; then
    echo -e "${RED}Error: File $FILE_PATH does not exist.${NC}"
    exit 1
fi
response=$(curl -s -w "%{http_code}" -X POST "$BASE_URL/s3/upload" \
    -F "bucket=$BUCKET" \
    -F "key=$FILE_KEY" \
    -F "file=@$FILE_PATH")
status=${response: -3}
body=${response%???}
check_status "$status" "$body" "Upload File"

# Download File
echo "Downloading file from bucket: $BUCKET with key: $FILE_KEY to: $DOWNLOAD_PATH"
response=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/s3/download/$BUCKET/$FILE_KEY" \
    -o "$DOWNLOAD_PATH" 2>&1)
status=${response: -3}
body=$(cat "$DOWNLOAD_PATH" 2>/dev/null || echo "Download failed")
if [ "$status" -eq 200 ]; then
    echo -e "${GREEN}Success: File downloaded to $DOWNLOAD_PATH${NC}"
else
    echo -e "${RED}Failed: Download File - HTTP $status${NC}"
    echo -e "${YELLOW}Error:${NC} $body"
    exit 1
fi

# Delete File
echo "Deleting file from bucket: $BUCKET with key: $FILE_KEY"
response=$(curl -s -w "%{http_code}" -X DELETE "$BASE_URL/s3/delete/$BUCKET/$FILE_KEY")
status=${response: -3}
body=${response%???}
check_status "$status" "$body" "Delete File"

# List Buckets
echo "Listing buckets"
response=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/s3/list")
status=${response: -3}
body=${response%???}
check_status "$status" "$body" "List Buckets"

# Update Bucket
echo "Updating bucket from: $BUCKET to: $NEW_BUCKET"
response=$(curl -s -w "%{http_code}" -X PUT "$BASE_URL/s3/update" \
    -H "Content-Type: application/json" \
    -d "{\"old_name\":\"$BUCKET\",\"new_name\":\"$NEW_BUCKET\",\"new_region\":\"$REGION\"}")
status=${response: -3}
body=${response%???}
check_status "$status" "$body" "Update Bucket"

# Delete Bucket
echo "Deleting bucket: $NEW_BUCKET"
response=$(curl -s -w "%{http_code}" -X DELETE "$BASE_URL/s3/delete/$NEW_BUCKET")
status=${response: -3}
body=${response%???}
check_status "$status" "$body" "Delete Bucket"

echo -e "${GREEN}All operations completed successfully!${NC}"

# clean up
rm test.txt downloaded_test.txt
