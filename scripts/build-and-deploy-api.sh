#!/bin/bash
set -e  # Exit on any error

# Colors for better output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Build and deploy CDK
print_status "Building CDK..."
cd ../cdk
npm ci
print_success "CDK dependencies installed"

print_status "Synthesizing CDK stack..."
npm run cdk:synth

print_status "Deploying CDK stack..."
npm run cdk:deploy SimplyClimbingApi

print_success "Deployment complete!"
