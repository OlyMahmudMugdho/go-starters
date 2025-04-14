#!/bin/bash

BASE_URL="http://localhost:8080"

# Colors
green='\033[0;32m'
red='\033[0;31m'
nc='\033[0m' # No Color

echo -e "${green}🔐 Registering users...${nc}"
curl -s -X POST $BASE_URL/register -H "Content-Type: application/json" -d '{"username":"user1","password":"pass123","role":"user"}'
echo -e "\n"

curl -s -X POST $BASE_URL/register -H "Content-Type: application/json" -d '{"username":"admin1","password":"adminpass","role":"admin"}'
echo -e "\n"

echo -e "${green}🔑 Logging in...${nc}"

USER_TOKEN=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"pass123"}' | jq -r '.token')

ADMIN_TOKEN=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin1","password":"adminpass"}' | jq -r '.token')

echo -e "🧑‍💻 User Token: ${green}${USER_TOKEN}${nc}"
echo -e "👑 Admin Token: ${green}${ADMIN_TOKEN}${nc}"

echo -e "\n🌐 ${green}Testing public route...${nc}"
curl -s $BASE_URL/public
echo -e "\n"

echo -e "\n🔒 ${green}Testing /user route as user...${nc}"
curl -s $BASE_URL/user -H "Authorization: Bearer $USER_TOKEN"
echo -e "\n"

echo -e "\n🔒 ${green}Testing /admin route as user (should fail)...${nc}"
curl -s -w "\nHTTP Status: %{http_code}\n" $BASE_URL/admin -H "Authorization: Bearer $USER_TOKEN"
echo -e "\n"

echo -e "\n🔒 ${green}Testing /admin route as admin...${nc}"
curl -s $BASE_URL/admin -H "Authorization: Bearer $ADMIN_TOKEN"
echo -e "\n"

echo -e "\n🌍 ${nc}Visit my website https://olymahmud.vercel.app${nc}"

echo -e "\n"