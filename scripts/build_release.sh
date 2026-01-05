#!/bin/bash
set -e

# Colors
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}Starting Siberia Release Build...${NC}"

# Ensure we are in project root
if [ ! -d "apps/backend" ]; then
    echo "Error: Must run from project root"
    exit 1
fi

# Clean previous builds
echo "Cleaning old builds..."
rm -rf apps/backend/build/bin

# Build Backend/Frontend via Wails
# We cd into apps/backend because wails.json is there
cd apps/backend

echo -e "${GREEN}Building macOS Universal Binary...${NC}"
wails build -clean -platform darwin/universal -production -o siberia

# Check if build succeeded
if [ -f "build/bin/siberia.app/Contents/MacOS/siberia" ]; then
    echo -e "${GREEN}Build Successful! Artifact: apps/backend/build/bin/siberia.app${NC}"
else
    echo "Error: Build artifact not found"
    exit 1
fi

# Optional: Add Windows cross-compile if needed
# echo "Building Windows..."
# wails build -platform windows/amd64 -production -o siberia.exe

echo -e "${GREEN}Done.${NC}"
