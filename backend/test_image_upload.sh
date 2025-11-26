#!/bin/bash

# Test script for image upload functionality
# This script tests the complete flow: register -> login -> create house -> upload image
# Requirements: curl, jq (for JSON parsing), and optionally ImageMagick (for creating test images)

BASE_URL="${BASE_URL:-http://localhost:8080}"

# Check if server is running
echo "🔍 Checking if server is running..."
HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" "${BASE_URL}/health" || echo -e "\n000")
HTTP_CODE=$(echo "$HEALTH_RESPONSE" | tail -n1)

if [ "$HTTP_CODE" != "200" ]; then
  echo "❌ Server is not running or not accessible at ${BASE_URL}"
  echo "   Please start the server first:"
  echo "   cd $(dirname "$0") && go run main.go"
  exit 1
fi

echo "✅ Server is running"
echo ""

BASE_URL="${BASE_URL}/api/v1"

# Check if jq is available
if ! command -v jq &> /dev/null; then
  echo "⚠️  Warning: jq is not installed. JSON parsing will be limited."
  echo "   Install jq for better output: sudo apt-get install jq (or brew install jq)"
  USE_JQ=false
else
  USE_JQ=true
fi
TIMESTAMP=$(date +%s)
TEST_EMAIL="test_landlord_${TIMESTAMP}@test.com"
TEST_PASSWORD="testpass123"
TEST_NAME="Test Landlord"
TEST_PHONE="+260123456789"

echo "🧪 Testing BondiHub Image Upload Functionality"
echo "=============================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Register a new landlord user
echo "📝 Step 1: Registering test landlord user..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"${TEST_EMAIL}\",
    \"password\": \"${TEST_PASSWORD}\",
    \"full_name\": \"${TEST_NAME}\",
    \"phone\": \"${TEST_PHONE}\",
    \"role\": \"landlord\"
  }")

# Check if registration was successful or if user already exists
if [ "$USE_JQ" = true ]; then
  TOKEN=$(echo $REGISTER_RESPONSE | jq -r '.data.token // empty')
else
  TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4 || echo "")
fi

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo -e "${YELLOW}⚠️  User might already exist, trying to login...${NC}"
  
  # Try to login instead
  LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
      \"email\": \"${TEST_EMAIL}\",
      \"password\": \"${TEST_PASSWORD}\"
    }")
  
  if [ "$USE_JQ" = true ]; then
    TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.data.token // empty')
    USER_ID=$(echo $LOGIN_RESPONSE | jq -r '.data.user.id // empty')
  else
    TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | cut -d'"' -f4 || echo "")
    USER_ID=$(echo $LOGIN_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4 | head -1 || echo "")
  fi
  
  if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
    echo -e "${RED}❌ Failed to register or login${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
  fi
  
  echo -e "${GREEN}✅ Logged in successfully${NC}"
else
  echo -e "${GREEN}✅ Registered successfully${NC}"
  if [ "$USE_JQ" = true ]; then
    USER_ID=$(echo $REGISTER_RESPONSE | jq -r '.data.user.id // empty')
  else
    USER_ID=$(echo $REGISTER_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4 | head -1 || echo "")
  fi
fi

echo "Token: ${TOKEN:0:20}..."
echo ""

# Step 2: Create a test house
echo "🏠 Step 2: Creating a test house..."
HOUSE_RESPONSE=$(curl -s -X POST "${BASE_URL}/houses" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ${TOKEN}" \
  -d "{
    \"title\": \"Test House for Image Upload\",
    \"description\": \"This is a test house created to test image upload functionality\",
    \"address\": \"123 Test Street, Lusaka, Zambia\",
    \"monthly_rent\": 5000.00,
    \"house_type\": \"apartment\",
    \"bedrooms\": 2,
    \"bathrooms\": 1,
    \"area\": 75.5,
    \"latitude\": -15.3875,
    \"longitude\": 28.3228
  }")

if [ "$USE_JQ" = true ]; then
  HOUSE_ID=$(echo $HOUSE_RESPONSE | jq -r '.data.house.id // empty')
else
  HOUSE_ID=$(echo $HOUSE_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4 | head -1 || echo "")
fi

if [ -z "$HOUSE_ID" ] || [ "$HOUSE_ID" == "null" ]; then
  echo -e "${RED}❌ Failed to create house${NC}"
  echo "Response: $HOUSE_RESPONSE"
  
  # Try to get existing houses
  echo -e "${YELLOW}⚠️  Trying to get existing houses...${NC}"
  HOUSES_RESPONSE=$(curl -s -X GET "${BASE_URL}/houses?limit=1" \
    -H "Authorization: Bearer ${TOKEN}")
  
  if [ "$USE_JQ" = true ]; then
    HOUSE_ID=$(echo $HOUSES_RESPONSE | jq -r '.data.houses[0].id // empty')
  else
    HOUSE_ID=$(echo $HOUSES_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4 | head -1 || echo "")
  fi
  
  if [ -z "$HOUSE_ID" ] || [ "$HOUSE_ID" == "null" ]; then
    echo -e "${RED}❌ No houses found and failed to create one${NC}"
    exit 1
  fi
  
  echo -e "${YELLOW}⚠️  Using existing house: ${HOUSE_ID}${NC}"
else
  echo -e "${GREEN}✅ House created successfully${NC}"
  echo "House ID: $HOUSE_ID"
fi

echo ""

# Step 3: Create a test image file
echo "🖼️  Step 3: Creating test image file..."
# Create a simple 1x1 pixel PNG image using ImageMagick or use a simple test image
TEST_IMAGE="test_image.png"

# Try to create a simple test image
# If ImageMagick is available, use it; otherwise download a small test image
if command -v convert &> /dev/null; then
  convert -size 200x200 xc:blue -pointsize 20 -fill white -gravity center -annotate +0+0 "TEST" "${TEST_IMAGE}"
  echo -e "${GREEN}✅ Created test image using ImageMagick${NC}"
else
  # Download a small test image (1x1 pixel PNG)
  echo -e "${YELLOW}⚠️  ImageMagick not found, downloading test image...${NC}"
  curl -s -o "${TEST_IMAGE}" "https://via.placeholder.com/200x200.png?text=TEST"
  if [ ! -f "${TEST_IMAGE}" ]; then
    echo -e "${RED}❌ Failed to create/download test image${NC}"
    echo "Please ensure you have an image file named 'test_image.png' in the current directory"
    exit 1
  fi
  echo -e "${GREEN}✅ Downloaded test image${NC}"
fi

echo ""

# Step 4: Upload the image
echo "📤 Step 4: Uploading image to Cloudinary..."
UPLOAD_RESPONSE=$(curl -s -X POST "${BASE_URL}/houses/${HOUSE_ID}/images" \
  -H "Authorization: Bearer ${TOKEN}" \
  -F "image=@${TEST_IMAGE}")

# Check if upload was successful
if [ "$USE_JQ" = true ]; then
  IMAGE_URL=$(echo $UPLOAD_RESPONSE | jq -r '.data.image.image_url // .data.image.ImageURL // empty')
  IMAGE_ID=$(echo $UPLOAD_RESPONSE | jq -r '.data.image.id // empty')
else
  IMAGE_URL=$(echo $UPLOAD_RESPONSE | grep -o '"image_url":"[^"]*' | cut -d'"' -f4 || echo $UPLOAD_RESPONSE | grep -o '"ImageURL":"[^"]*' | cut -d'"' -f4 || echo "")
  IMAGE_ID=$(echo $UPLOAD_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4 | head -1 || echo "")
fi

if [ -z "$IMAGE_URL" ] || [ "$IMAGE_URL" == "null" ]; then
  echo -e "${RED}❌ Failed to upload image${NC}"
  echo "Response: $UPLOAD_RESPONSE"
  # Clean up test image
  rm -f "${TEST_IMAGE}"
  exit 1
fi

echo -e "${GREEN}✅ Image uploaded successfully!${NC}"
echo ""
echo "📊 Upload Results:"
echo "=================="
echo "Image ID: $IMAGE_ID"
echo "Image URL: $IMAGE_URL"
echo ""
echo "Full Response:"
if [ "$USE_JQ" = true ]; then
  echo "$UPLOAD_RESPONSE" | jq '.'
else
  echo "$UPLOAD_RESPONSE"
fi
echo ""

# Clean up test image
rm -f "${TEST_IMAGE}"
echo -e "${GREEN}🧹 Cleaned up test image file${NC}"
echo ""
echo -e "${GREEN}✅ All tests passed! Image upload is working correctly.${NC}"

