#!/bin/bash
# Check if Docker is already installed
if [ -x "$(command -v docker)" ]; then
  echo "Docker is already installed."
  exit 0
fi
# Check if Homebrew is installed
if ! [ -x "$(command -v brew)" ]; then
  echo "Homebrew is not installed. Installing Homebrew..."
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
fi
# Install Docker using Homebrew
brew install --cask docker
# Start Docker Desktop
open /Applications/Docker.app
# Wait for Docker to start
echo "Waiting for Docker to start..."
while ! docker info >/dev/null 2>&1; do sleep 1; done
# Enable Docker CLI to connect to the Docker Desktop instance
docker context use default
echo "Docker is now installed and running on your Mac."