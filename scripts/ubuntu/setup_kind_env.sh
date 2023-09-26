#!/bin/bash
# Check if Go is installed
if ! command -v go &>/dev/null; then
    echo "Go is not installed. Installing Go..."
    sudo apt update
    wget https://go.dev/dl/go1.21.1.linux-amd64.tar.gz
    sudo tar -xvf go1.21.1.linux-amd64.tar.gz
    sudo mv go /usr/local
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
    source ~/.profile
    # Validate go version
    go version
fi
# Check if kubectl is installed
if ! command -v kubectl &>/dev/null; then
    echo "kubectl is not installed. Installing kubectl..."
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x kubectl
    sudo mv kubectl /usr/local/bin/
fi
# Check if kind is installed
if ! command -v kind &>/dev/null; then
    echo "kind is not installed. Installing kind..."
    curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.11.1/kind-linux-amd64
    chmod +x ./kind
    sudo mv ./kind /usr/local/bin/kind
fi
# Check if necessary Go environment variables are set
if [ -z "$GOPATH" ]; then
    echo "GOPATH environment variable is not set. Setting GOPATH..."
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    source ~/.bashrc
fi
# Check if kubernetes client-go library is installed
if ! go list -f '{{ .Name }}' -m k8s.io/client-go &>/dev/null; then
    echo "kubernetes client-go library is not installed. Installing..."
    go get -u k8s.io/client-go/...
fi
# Check if kubebuilder is installed
if ! command -v kubebuilder &>/dev/null; then
    echo "kubebuilder is not installed. Installing kubebuilder..."
    OS=$(go env GOOS)
    ARCH=$(go env GOARCH)
    curl -L "https://go.kubebuilder.io/dl/2.3.2/${OS}/${ARCH}" | tar -xz -C /tmp/
    sudo mv /tmp/kubebuilder_2.3.2_${OS}_${ARCH} /usr/local/kubebuilder
    export PATH=$PATH:/usr/local/kubebuilder/bin
    echo 'export PATH=$PATH:/usr/local/kubebuilder/bin' >> ~/.bashrc
fi
# Create a new Kind cluster
echo "Creating a new Kind cluster..."
kind create cluster --wait 2m
# Configure kubectl to use the new Kind cluster
echo "Configuring kubectl to use the new Kind cluster..."
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
# Verify cluster status
echo "Verifying the cluster status..."
kubectl cluster-info
# Kind cluster setup with Kubernetes API machinery completed successfully
echo "Kind cluster setup with Kubernetes API machinery completed successfully!"#!/bin/bash
# Check if Go is installed
if ! command -v go &>/dev/null; then
    echo "Go is not installed. Installing Go..."
    sudo apt update
    sudo apt install -y golang
fi
# Check if kubectl is installed
if ! command -v kubectl &>/dev/null; then
    echo "kubectl is not installed. Installing kubectl..."
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    chmod +x kubectl
    sudo mv kubectl /usr/local/bin/
fi
# Check if kind is installed
if ! command -v kind &>/dev/null; then
    echo "kind is not installed. Installing kind..."
    curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.11.1/kind-linux-amd64
    chmod +x ./kind
    sudo mv ./kind /usr/local/bin/kind
fi
# Check if necessary Go environment variables are set
if [ -z "$GOPATH" ]; then
    echo "GOPATH environment variable is not set. Setting GOPATH..."
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    source ~/.bashrc
fi
# Check if kubernetes client-go library is installed
if ! go list -f '{{ .Name }}' -m k8s.io/client-go &>/dev/null; then
    echo "kubernetes client-go library is not installed. Installing..."
    go get -u k8s.io/client-go/...
fi
# Check if kubebuilder is installed
if ! command -v kubebuilder &>/dev/null; then
    echo "kubebuilder is not installed. Installing kubebuilder..."
    OS=$(go env GOOS)
    ARCH=$(go env GOARCH)
    curl -L "https://go.kubebuilder.io/dl/2.3.2/${OS}/${ARCH}" | tar -xz -C /tmp/
    sudo mv /tmp/kubebuilder_2.3.2_${OS}_${ARCH} /usr/local/kubebuilder
    export PATH=$PATH:/usr/local/kubebuilder/bin
    echo 'export PATH=$PATH:/usr/local/kubebuilder/bin' >> ~/.bashrc
fi
# Create a new Kind cluster
echo "Creating a new Kind cluster..."
kind create cluster --wait 2m
# Configure kubectl to use the new Kind cluster
echo "Configuring kubectl to use the new Kind cluster..."
export KUBECONFIG="$(kind get kubeconfig-path --name="kind")"
# Verify cluster status
echo "Verifying the cluster status..."
kubectl cluster-info

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

echo "Kind cluster setup with Kubernetes API machinery develop env completed successfully!"
