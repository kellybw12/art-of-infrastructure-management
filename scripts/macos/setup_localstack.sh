#!/bin/bash
​
# Check if Homebrew is installed, and install it if not
if ! command -v brew &> /dev/null; then
    echo "Homebrew is not installed. Installing..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
fi
​
# Install Python 3.x
brew install python@3.9  # You can specify a specific version if desired
​
# Check if Python 3 is installed correctly
python3 --version
​
# Install pip (Python package manager)
python3 -m ensurepip --default-pip
​
# Check if pip is already installed
if command -v pip3 &> /dev/null; then
    echo "pip is already installed. Checking for updates..."
    pip3 install --upgrade pip
    echo "pip has been updated to the latest version."
else
    echo "pip is not installed. Installing pip..."
    curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
    python3 get-pip.py
    rm get-pip.py
    echo "pip has been installed."
fi
​
# Verify the pip installation
pip3 --version
​
​
# Check if awscli is already installed
if command -v awscli &> /dev/null; then
    echo "awscli is already installed. Checking for updates..."
    pip3 install --upgrade awscli
    echo "awscli has been updated to the latest version."
else
    echo "awscli is not installed. Installing awscli..."
    pip3 install awscli
    # Add the installation paths to your PATH
    python_lib_bin_path=$(ls -d ~/Library/Python/*/bin)
    echo "export PATH=$PATH:$python_lib_bin_path" >> ~/.bashrc  # Adjust for your shell if needed
    source ~/.bashrc
​
    echo "awscli has been installed."
fi
​
​
# Install LocalStack using pip
pip3 install localstack
​
# Start LocalStack using Docker
localstack start -d
​
# Configure AWS CLI to use LocalStack
aws configure set aws_access_key_id test
aws configure set aws_secret_access_key test
aws configure set default.region us-east-1
aws configure set default.output json
aws configure set endpoint_url http://localhost:4566
​
echo "Python, pip, AWS CLI, and LocalStack setup complete."