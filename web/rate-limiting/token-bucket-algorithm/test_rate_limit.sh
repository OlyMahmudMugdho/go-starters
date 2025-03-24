#!/bin/bash

URL="http://localhost:8080/"
NUM_REQUESTS=10  # Number of requests to send
DELAY=0.005         # Delay between requests in seconds

echo "Testing rate limiting with $NUM_REQUESTS requests..."

for ((i=1; i<=NUM_REQUESTS; i++)); do
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" $URL)
    
    if [ "$RESPONSE" -eq 200 ]; then
        echo "✅ Request $i: Allowed (200 OK)"
    elif [ "$RESPONSE" -eq 429 ]; then
        echo "❌ Request $i: Blocked (429 Too Many Requests)"
    else
        echo "⚠️ Request $i: Unexpected response ($RESPONSE)"
    fi

    sleep $DELAY
done

echo "Test completed."
