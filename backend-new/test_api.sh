#!/bin/bash

# Base URL
API_URL="http://localhost:8080"
COOKIE_JAR="cookies.txt"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}Starting API Tests...${NC}"

# 1. Register
echo -e "\n${GREEN}[1] Testing Registration...${NC}"
TIMESTAMP=$(date +%s)
EMAIL="user${TIMESTAMP}@test.com"
USERNAME="user${TIMESTAMP}"
PASSWORD="password123"

# Note: Register expects Multipart Form Data because it handles avatar uploads
curl -s -X POST "$API_URL/register" \
     -F "username=$USERNAME" \
     -F "email=$EMAIL" \
     -F "password=$PASSWORD" \
     -F "confirmpassword=$PASSWORD" \
     -F "first_name=Test" \
     -F "last_name=User" \
     -F "date_of_birth=1990-01-01" \
     -F "gender=male" \
     -F "nickname=tester" \
     -F "about=I am a test user" \
     -F "avatar=@/dev/null;filename=" \
     | grep "successfully" > /dev/null

# We don't check strict success here because sometimes it redirects or returns different JSON.
# But let's check if we can login next.
echo -e "${GREEN}Registration Request Sent ($USERNAME)${NC}"

# 2. Login
echo -e "\n${GREEN}[2] Testing Login...${NC}"
# Login expects JSON
LOGIN_PAYLOAD="{\"username\": \"$USERNAME\", \"password\": \"$PASSWORD\"}"

curl -c $COOKIE_JAR -s -X POST "$API_URL/login" \
     -H "Content-Type: application/json" \
     -d "$LOGIN_PAYLOAD" > login_response.txt

if grep -q "ok" login_response.txt; then
    echo -e "${GREEN}Login Successful${NC}"
else
    echo -e "${RED}Login Failed${NC}"
    cat login_response.txt
    exit 1
fi

# 3. Check Auth State
echo -e "\n${GREEN}[3] Testing Check Auth State...${NC}"
curl -b $COOKIE_JAR -s -X GET "$API_URL/checkstate" > checkstate_response.txt

if grep -q "true" checkstate_response.txt; then
    echo -e "${GREEN}Auth Check Passed${NC}"
else
    echo -e "${RED}Auth Check Failed${NC}"
    cat checkstate_response.txt
fi

# 4. Create Post
echo -e "\n${GREEN}[4] Testing Create Post...${NC}"
# Create Post expects Multipart Form Data
curl -b $COOKIE_JAR -s -X POST "$API_URL/Createpost" \
     -F "title=Test Post $TIMESTAMP" \
     -F "content=This is a test post content created at $TIMESTAMP" \
     -F "privacy=public" > createpost_response.txt

if grep -q "id" createpost_response.txt; then
    echo -e "${GREEN}Create Post Successful${NC}"
else
    echo -e "${RED}Create Post Failed${NC}"
    cat createpost_response.txt
fi

# 5. Get Posts
echo -e "\n${GREEN}[5] Testing Get Posts...${NC}"
curl -b $COOKIE_JAR -s -X GET "$API_URL/getposts" > getposts_response.txt

if grep -q "Test Post $TIMESTAMP" getposts_response.txt; then
    echo -e "${GREEN}Get Posts Successful (Found new post)${NC}"
else
    echo -e "${RED}Get Posts Failed or Post not found${NC}"
fi

# 6. Get Profile
# We need the user ID. Login response gave us "user_id".
USER_ID=$(grep -o '"user_id":[0-9]*' login_response.txt | cut -d: -f2)
if [ -n "$USER_ID" ]; then
    echo -e "\n${GREEN}[6] Testing Get Profile (ID: $USER_ID)...${NC}"
    curl -b $COOKIE_JAR -s -X GET "$API_URL/profile/?id=$USER_ID" > profile_response.txt
    
    if grep -q "Test" profile_response.txt; then
        echo -e "${GREEN}Get Profile Successful${NC}"
    else
        echo -e "${RED}Get Profile Failed${NC}"
        # cat profile_response.txt
    fi
fi

echo -e "\n${GREEN}Tests Completed.${NC}"
rm $COOKIE_JAR login_response.txt checkstate_response.txt createpost_response.txt getposts_response.txt profile_response.txt
