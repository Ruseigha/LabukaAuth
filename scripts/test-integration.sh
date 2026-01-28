#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Integration Test Runner${NC}"
echo ""

# Check if MongoDB is running
echo -e "${YELLOW}Checking MongoDB...${NC}"
if ! mongosh --eval "db.adminCommand('ping')" --quiet > /dev/null 2>&1; then
  echo -e "${RED}✗ MongoDB is not running!${NC}"
  echo ""
  echo "Please start MongoDB:"
  echo "  macOS:   brew services start mongodb-community@7.0"
  echo "  Linux:   sudo systemctl start mongod"
  echo "  Windows: net start MongoDB"
  exit 1
fi
echo -e "${GREEN}✓ MongoDB is running${NC}"
echo ""

# Clean test database
echo -e "${YELLOW}Cleaning test database...${NC}"
mongosh auth_service_test --eval "db.dropDatabase()" --quiet
echo -e "${GREEN}✓ Database cleaned${NC}"
echo ""

# Run tests
echo -e "${YELLOW}Running integration tests...${NC}"
go test -v -race -cover ./test/integration/...

# Exit with test result
exit $?