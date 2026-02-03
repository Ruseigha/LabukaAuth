#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Generating Protocol Buffer code...${NC}"

# Create output directory
mkdir -p internal/delivery/grpc/proto

# Generate Go code
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/auth.proto

echo -e "${GREEN}✓ Protocol Buffer code generated${NC}"