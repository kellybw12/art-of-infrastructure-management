#!/bin/bash
# Check if Docker is already installed
if ! [ -x "$(command -v docker)" ]; then
    echo "Docker is not installed. Installing Docker..."
    # Update package list and install prerequisites
    sudo apt update
    sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
    # Add Docker repository
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    # Update package list again and install Docker
    sudo apt update
    sudo apt install -y docker-ce docker-ce-cli containerd.io
    echo "Docker has been installed."
else
    echo "Docker is already installed."
fi