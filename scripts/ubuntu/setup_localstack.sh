#!/bin/bash
# Function to check if a command is available
command_exists() {
    command -v "$1" &> /dev/null
}
# Check and install Python and Pip if not already installed
if ! command_exists python || ! command_exists pip; then
    echo "Python and/or Pip is not installed. Installing Python and Pip..."
    sudo apt update
    sudo apt install -y python3 python3-pip
    echo "Python and Pip have been installed."
else
    echo "Python and Pip are already installed."
fi
# Check and install AWS CLI if not already installed
if ! command_exists aws; then
    echo "AWS CLI is not installed. Installing AWS CLI..."
    sudo pip install awscli
    echo "AWS CLI has been installed."
else
    echo "AWS CLI is already installed."
fi
# Set up LocalStack
echo "Setting up LocalStack..."
sudo pip install localstack
localstack start -d

# Configure AWS CLI to use LocalStack
aws configure set aws_access_key_id test
aws configure set aws_secret_access_key test
aws configure set default.region us-east-1
aws configure set default.output json
aws configure set endpoint_url http://localhost:4566
​
# Script completed
echo "Setup complete."