#!/bin/bash
set -e  # Exit on any error

# Colors for better output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get environment from argument (default to dev)
ENV=${1:-dev}

# Validate environment
if [[ "$ENV" != "dev" && "$ENV" != "prod" ]]; then
    print_error "Invalid environment: $ENV. Must be 'dev' or 'prod'"
    exit 1
fi

print_status "Deploying to $ENV environment"

# Get the repo root relative to this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Determine stack names
if [ "$ENV" == "prod" ]; then
    STACK_NAME="SidekickFrontend"
    CERT_STACK_NAME="SidekickCertificate"
else
    STACK_NAME="SidekickFrontend-${ENV}"
    CERT_STACK_NAME="SidekickCertificate-${ENV}"
fi

# Build frontend
print_status "Building frontend..."
cd "$REPO_ROOT/apps/web"
npm ci  # Use ci for faster, more reliable installs
npx vite build --mode $ENV && npx tsc
print_success "Frontend build complete"

# Build and deploy CDK
print_status "Building CDK..."
cd "$REPO_ROOT/cdk"
npm ci
print_success "CDK dependencies installed"

print_status "Synthesizing CDK stack..."
npx cdk synth --context env=$ENV

print_status "Deploying CDK stack: $STACK_NAME..."
npx cdk deploy $CERT_STACK_NAME $STACK_NAME --context env=$ENV --require-approval never

print_success "Deployment complete!"
