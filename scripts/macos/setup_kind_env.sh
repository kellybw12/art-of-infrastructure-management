#!/bin/bash
# Check if Homebrew is installed
if ! [ -x "$(command -v brew)" ]; then
  echo "Homebrew is not installed. Installing Homebrew..."
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
fi

# Install Golang using Homebrew
echo "Install Golang using Homebrew"
brew install go

# Set up Golang environment variables
echo "Set up Golang environment variables"
echo 'export GOPATH=$HOME/go' >> ~/.bash_profile
echo 'export GOBIN=$GOPATH/bin' >> ~/.bash_profile
echo 'export PATH=$PATH:$GOBIN' >> ~/.bash_profile
source ~/.bash_profile

# Install Kind and kubectl using Homebrew
echo "Install Kind and kubectl using Homebrew"
brew install kind kubernetes-cli

# Create a Kind configuration file
echo "Create a Kind configuration file"
cat << EOF > kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
EOF

# Create the Kind cluster
echo "Create the Kind cluster"
kind create cluster --config kind-config.yaml

# Set the environment variables to use Kind's Kubernetes cluster
export KUBECONFIG=$(kind get kubeconfig-path --name="kind")
# Clone Kubernetes source code (if you haven't already)
mkdir -p $GOPATH/src/k8s.io
cd $GOPATH/src/k8s.io
git clone https://github.com/kubernetes/kubernetes.git
cd kubernetes

# Verify the nodes in the cluster
kubectl get nodes
# Verify other resources (pods, deployments, services, etc.)
kubectl get pods,deployments,services --all-namespaces

# download kubebuilder and install locally.
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/latest/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/

echo "Kind cluster setup with Kubernetes API machinery develop env completed successfully!"
