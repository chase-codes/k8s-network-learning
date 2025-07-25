#!/bin/bash

# OSI Model Kubernetes Lab Setup Script
# This script creates a kind cluster, deploys nginx, captures packets, and prepares data for analysis

set -e

# Configuration
CLUSTER_NAME="netlab-osi"
NAMESPACE="default"
CAPTURE_FILE="modules/01-osi-model/assets/https-nginx.pcap"
ASSETS_DIR="modules/01-osi-model/assets"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if Docker is running
is_docker_running() {
    docker info >/dev/null 2>&1
}

# Function to start Docker
start_docker() {
    print_status "Docker is not running. Attempting to start Docker..."
    
    # Detect OS and start Docker accordingly
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS - try different Docker installation methods
        
        # Check for Homebrew Docker service first
        if command_exists brew && brew services list 2>/dev/null | grep -q "docker"; then
            print_status "Starting Docker via Homebrew services..."
            brew services start docker
        elif command_exists dockerd; then
            print_status "Starting Docker daemon..."
            # Start Docker daemon in background
            sudo dockerd > /dev/null 2>&1 &
            sleep 3
        # Fallback to Docker Desktop if available
        elif [ -d "/Applications/Docker.app" ]; then
            print_status "Starting Docker Desktop..."
            open -a Docker
        elif [ -d "/Applications/Docker Desktop.app" ]; then
            print_status "Starting Docker Desktop..."
            open -a "Docker Desktop"
        else
            print_error "Docker daemon not available"
            print_error ""
            print_error "You have Docker CLI installed, but no daemon is running."
            print_error "On macOS, you need either:"
            print_error "• Docker Desktop: https://docker.com"
            print_error "• Or install Docker daemon: 'brew install docker docker-machine'"
            print_error ""
            print_error "If you have Docker Desktop, please start it manually and try again."
            return 1
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux - start Docker daemon
        if command_exists systemctl; then
            print_status "Starting Docker service with systemctl..."
            sudo systemctl start docker
        elif command_exists service; then
            print_status "Starting Docker service..."
            sudo service docker start
        else
            print_error "Unable to start Docker automatically on this system"
            print_error "Please start Docker manually and try again"
            return 1
        fi
    else
        print_error "Unsupported operating system: $OSTYPE"
        print_error "Please start Docker manually and try again"
        return 1
    fi
    
    return 0
}

# Function to wait for Docker to be ready
wait_for_docker() {
    print_status "Waiting for Docker to be ready..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if is_docker_running; then
            print_success "Docker is ready!"
            return 0
        fi
        
        print_status "Waiting for Docker... (attempt $attempt/$max_attempts)"
        sleep 2
        ((attempt++))
    done
    
    print_error "Docker failed to start within 60 seconds"
    print_error "Please start Docker manually and try again"
    return 1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    local missing_deps=()
    
    if ! command_exists docker; then
        missing_deps+=("docker")
    fi
    
    if ! command_exists kind; then
        missing_deps+=("kind")
    fi
    
    if ! command_exists kubectl; then
        missing_deps+=("kubectl")
    fi
    
    if ! command_exists tcpdump; then
        missing_deps+=("tcpdump")
    fi
    
    if ! command_exists tshark; then
        missing_deps+=("tshark")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        print_error "Missing required dependencies: ${missing_deps[*]}"
        echo ""
        echo "Please install the missing dependencies:"
        echo "  macOS: brew install docker kind kubernetes-cli tcpdump wireshark"
        echo "  Linux: Use your package manager to install the above tools"
        echo ""
        exit 1
    fi
    
    # Check if Docker is running, start if needed
    if ! is_docker_running; then
        if ! start_docker; then
            exit 1
        fi
        
        if ! wait_for_docker; then
            exit 1
        fi
    else
        print_success "Docker is already running"
    fi
    
    print_success "All prerequisites are satisfied"
}

# Function to create assets directory
create_assets_dir() {
    print_status "Creating assets directory..."
    mkdir -p "$ASSETS_DIR"
    print_success "Assets directory created at $ASSETS_DIR"
}

# Function to create kind cluster
create_cluster() {
    print_status "Creating kind cluster '$CLUSTER_NAME'..."
    
    # Check if cluster already exists
    if kind get clusters | grep -q "^$CLUSTER_NAME$"; then
        print_warning "Cluster '$CLUSTER_NAME' already exists. Deleting it first..."
        kind delete cluster --name "$CLUSTER_NAME"
    fi
    
    # Create cluster with custom config
    cat <<EOF | kind create cluster --name "$CLUSTER_NAME" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30080
    hostPort: 30080
    protocol: TCP
EOF
    
    # Set kubectl context
    kubectl cluster-info --context kind-"$CLUSTER_NAME"
    
    print_success "Kind cluster '$CLUSTER_NAME' created successfully"
}

# Function to deploy nginx
deploy_nginx() {
    print_status "Deploying nginx to the cluster..."
    
    # Create nginx deployment
    cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:alpine
        ports:
        - containerPort: 80
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "100m"
---
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  selector:
    app: nginx
  ports:
  - protocol: TCP
    port: 80
    targetPort: 80
  type: ClusterIP
EOF

    # Wait for nginx to be ready
    print_status "Waiting for nginx deployment to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/nginx
    
    print_success "Nginx deployed successfully"
}

# Function to deploy busybox for testing
deploy_busybox() {
    print_status "Deploying busybox test pod..."
    
    cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: busybox
  labels:
    app: busybox
spec:
  containers:
  - name: busybox
    image: busybox:latest
    command: ['sh', '-c', 'sleep 3600']
    resources:
      requests:
        memory: "32Mi"
        cpu: "25m"
      limits:
        memory: "64Mi"
        cpu: "50m"
EOF

    # Wait for busybox to be ready
    print_status "Waiting for busybox pod to be ready..."
    kubectl wait --for=condition=ready --timeout=300s pod/busybox
    
    print_success "Busybox pod deployed successfully"
}

# Function to capture packets using host tcpdump (fallback)
capture_with_host_tcpdump() {
    local nginx_ip="$1"
    print_status "Using host tcpdump to capture traffic..."
    
    # Get docker network for kind
    local docker_network=$(docker network inspect kind | grep '"Subnet"' | head -1 | sed 's/.*"\([^"]*\)".*/\1/')
    
    if [ -z "$docker_network" ]; then
        print_error "Could not determine docker network for kind cluster"
        return 1
    fi
    
    print_status "Capturing traffic on network: $docker_network"
    
    # Start host tcpdump in background
    sudo tcpdump -i any -w "$CAPTURE_FILE" net "$docker_network" and host "$nginx_ip" &
    local tcpdump_pid=$!
    
    # Wait for tcpdump to start
    sleep 3
    
    # Make HTTP requests
    print_status "Making HTTP request from busybox to nginx..."
    kubectl exec busybox -- wget -qO- http://nginx/ > /dev/null || true
    kubectl exec busybox -- wget -qO- http://nginx/ > /dev/null || true
    kubectl exec busybox -- wget -qO- http://nginx/ > /dev/null || true
    
    # Wait and stop capture
    sleep 3
    print_status "Stopping host packet capture..."
    sudo kill $tcpdump_pid 2>/dev/null || true
    sleep 2
    
    if [ -f "$CAPTURE_FILE" ]; then
        print_success "Host-based packet capture completed"
        return 0
    else
        print_error "Host-based packet capture failed"
        return 1
    fi
}

# Function to capture packets
capture_packets() {
    print_status "Setting up packet capture..."
    
    # Get the kind container name
    local container_name="${CLUSTER_NAME}-control-plane"
    
    # Get nginx pod IP
    local nginx_ip=$(kubectl get pod -l app=nginx -o jsonpath='{.items[0].status.podIP}')
    print_status "Nginx pod IP: $nginx_ip"
    
    # Install tcpdump in the kind container if not present
    print_status "Installing tcpdump in kind container..."
    if ! docker exec "$container_name" sh -c "apt-get update > /dev/null 2>&1 && apt-get install -y tcpdump > /dev/null 2>&1"; then
        print_warning "Failed to install tcpdump in container, using host-based capture..."
        capture_with_host_tcpdump "$nginx_ip"
        return $?
    fi
    
    # Start packet capture in background (use /var/tmp to avoid tmpfs issues)
    print_status "Starting packet capture on kind container..."
    docker exec -d "$container_name" sh -c "tcpdump -i any -w /var/tmp/capture.pcap host $nginx_ip"
    
    # Wait a moment for tcpdump to start
    sleep 3
    
    # Make HTTP request from busybox to nginx
    print_status "Making HTTP request from busybox to nginx..."
    kubectl exec busybox -- wget -qO- http://nginx/ > /dev/null || true
    kubectl exec busybox -- wget -qO- http://nginx/ > /dev/null || true
    kubectl exec busybox -- wget -qO- http://nginx/ > /dev/null || true
    
    # Wait a bit more to capture all packets
    sleep 3
    
    # Stop packet capture properly
    print_status "Stopping packet capture..."
    docker exec "$container_name" pkill -f tcpdump || true
    sleep 2  # Give tcpdump time to flush buffers and write file
    
    # Verify capture file exists before copying
    print_status "Verifying packet capture file..."
    if docker exec "$container_name" test -f /var/tmp/capture.pcap; then
        local file_size=$(docker exec "$container_name" stat -c%s /var/tmp/capture.pcap)
        print_status "Found capture file (${file_size} bytes)"
        
        # Copy capture file from container
        print_status "Copying packet capture file..."
        docker cp "$container_name:/var/tmp/capture.pcap" "$CAPTURE_FILE"
    else
        print_error "Capture file not found in container. Checking for running tcpdump processes..."
        docker exec "$container_name" ps aux | grep tcpdump || true
        print_warning "Container-based capture failed, trying host-based approach..."
        capture_with_host_tcpdump "$nginx_ip"
        return $?
    fi
    
    if [ -f "$CAPTURE_FILE" ]; then
        print_success "Packet capture saved to $CAPTURE_FILE"
        
        # Show basic packet info
        local packet_count=$(tshark -r "$CAPTURE_FILE" -q -z io,stat,0 | grep -E "Interval|packets" | tail -1 | awk '{print $3}')
        print_status "Captured $packet_count packets"
        print_success "Lab setup completed successfully!"
    else
        print_error "Failed to save packet capture file"
        return 1
    fi
}

# Function to create OSI diagram
create_osi_diagram() {
    print_status "Creating OSI model diagram..."
    
    cat > "$ASSETS_DIR/osi-diagram.txt" <<'EOF'
┌─────────────────────────────────────────────────────────────────┐
│                        OSI MODEL LAYERS                        │
├─────────────────────────────────────────────────────────────────┤
│  7  │ APPLICATION  │ HTTP, FTP, SMTP, DNS, SSH                │
│     │    LAYER     │ User Interface & Network Services        │
├─────┼──────────────┼───────────────────────────────────────────┤
│  6  │ PRESENTATION │ SSL/TLS, JPEG, MPEG, Encryption         │
│     │    LAYER     │ Data Translation & Encryption            │
├─────┼──────────────┼───────────────────────────────────────────┤
│  5  │   SESSION    │ NetBIOS, RPC, SQL Sessions               │
│     │    LAYER     │ Session Management & Control             │
├─────┼──────────────┼───────────────────────────────────────────┤
│  4  │  TRANSPORT   │ TCP, UDP, Port Numbers                   │
│     │    LAYER     │ End-to-End Delivery & Flow Control       │
├─────┼──────────────┼───────────────────────────────────────────┤
│  3  │   NETWORK    │ IP, ICMP, OSPF, BGP, Routing            │
│     │    LAYER     │ Logical Addressing & Path Selection      │
├─────┼──────────────┼───────────────────────────────────────────┤
│  2  │  DATA LINK   │ Ethernet, Wi-Fi, MAC Addresses          │
│     │    LAYER     │ Node-to-Node Delivery & Error Detection  │
├─────┼──────────────┼───────────────────────────────────────────┤
│  1  │   PHYSICAL   │ Cables, Radio, Fiber, Electrical Signals│
│     │    LAYER     │ Physical Transmission Medium             │
└─────┴──────────────┴───────────────────────────────────────────┘

Mnemonic: "Please Do Not Throw Sausage Pizza Away"
          Physical, Data Link, Network, Transport, Session, Presentation, Application

In our Kubernetes lab:
- Layer 7: HTTP GET request from busybox to nginx
- Layer 4: TCP connection on port 80
- Layer 3: IP routing between pod IPs
- Layer 2: Ethernet frames within the kind container
- Layer 1: Virtual network interfaces (veth pairs)
EOF
    
    print_success "OSI diagram created at $ASSETS_DIR/osi-diagram.txt"
}

# Function to show lab info
show_lab_info() {
    print_success "Lab setup completed successfully!"
    echo ""
    echo "📋 Lab Components:"
    echo "  • Kind cluster: $CLUSTER_NAME"
    echo "  • Nginx service running on cluster IP"
    echo "  • Busybox pod for testing"
    echo "  • Packet capture: $CAPTURE_FILE"
    echo "  • OSI diagram: $ASSETS_DIR/osi-diagram.txt"
    echo ""
    echo "🔍 Next Steps:"
    echo "  1. Run 'netlab module 01-osi-model' to start the interactive module"
    echo "  2. Press 'v' in the module to access the packet analysis lab"
    echo "  3. Use 'kubectl' commands to explore the cluster:"
    echo "     kubectl get pods"
    echo "     kubectl get services"
    echo "     kubectl exec busybox -- wget -qO- http://nginx/"
    echo ""
    echo "🧹 Cleanup:"
    echo "  • Run 'kind delete cluster --name $CLUSTER_NAME' when finished"
    echo ""
}

# Function to stop Docker
stop_docker() {
    print_status "Stopping Docker..."
    
    # Detect OS and stop Docker accordingly
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS - try different Docker stopping methods
        
        # Check for Homebrew Docker service first
        if command_exists brew && brew services list | grep -q "docker.*started"; then
            print_status "Stopping Docker via Homebrew services..."
            brew services stop docker
            print_success "Docker service stopped"
        # Check for Docker Desktop processes
        elif pgrep -f "Docker Desktop" > /dev/null; then
            print_status "Quitting Docker Desktop..."
            osascript -e 'quit app "Docker Desktop"' 2>/dev/null || true
            # Alternative method if osascript fails
            killall "Docker Desktop" 2>/dev/null || true
            print_success "Docker Desktop stopped"
        # Try to stop dockerd if running
        elif pgrep -f "dockerd" > /dev/null; then
            print_status "Stopping Docker daemon..."
            sudo pkill dockerd 2>/dev/null || true
            print_success "Docker daemon stopped"
        else
            print_status "Docker is not running"
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux - stop Docker daemon
        if command_exists systemctl; then
            print_status "Stopping Docker service with systemctl..."
            sudo systemctl stop docker
        elif command_exists service; then
            print_status "Stopping Docker service..."
            sudo service docker stop
        else
            print_warning "Unable to stop Docker automatically on this system"
            print_status "You may need to stop Docker manually"
            return 1
        fi
        print_success "Docker service stopped"
    else
        print_warning "Unsupported operating system: $OSTYPE"
        print_status "You may need to stop Docker manually"
        return 1
    fi
    
    return 0
}

# Function to ask user confirmation
ask_user_confirmation() {
    local question="$1"
    local default="${2:-n}"
    
    if [ "$default" = "y" ]; then
        prompt="$question [Y/n]: "
    else
        prompt="$question [y/N]: "
    fi
    
    read -p "$prompt" response
    
    case "$response" in
        [yY]|[yY][eE][sS])
            return 0
            ;;
        [nN]|[nN][oO])
            return 1
            ;;
        "")
            if [ "$default" = "y" ]; then
                return 0
            else
                return 1
            fi
            ;;
        *)
            echo "Please answer yes or no."
            ask_user_confirmation "$question" "$default"
            ;;
    esac
}

# Function to cleanup
cleanup() {
    print_status "Cleaning up NetLab OSI environment..."
    
    # Delete kind cluster
    if kind get clusters | grep -q "^$CLUSTER_NAME$"; then
        kind delete cluster --name "$CLUSTER_NAME"
        print_success "Cluster '$CLUSTER_NAME' deleted"
    else
        print_status "Cluster '$CLUSTER_NAME' was not found"
    fi
    
    # Ask if user wants to stop Docker
    echo ""
    print_status "Lab cleanup completed."
    echo ""
    echo "Do you want to stop Docker as well?"
    echo "Note: This will stop Docker completely, which may affect other containers you have running."
    echo ""
    
    if ask_user_confirmation "Stop Docker?"; then
        stop_docker
    else
        print_status "Docker left running"
        echo ""
        echo "💡 If you want to stop Docker later:"
        if [[ "$OSTYPE" == "darwin"* ]]; then
            echo "   macOS: Quit Docker Desktop from the menu bar"
        else
            echo "   Linux: sudo systemctl stop docker"
        fi
    fi
}

# Main function
main() {
    echo "🌐 NetLab OSI Model - Kubernetes Lab Setup"
    echo "=========================================="
    echo ""
    
    # Parse arguments
    case "${1:-setup}" in
        "setup")
            check_prerequisites
            create_assets_dir
            create_cluster
            deploy_nginx
            deploy_busybox
            capture_packets
            create_osi_diagram
            show_lab_info
            ;;
        "cleanup")
            cleanup
            ;;
        "capture")
            capture_packets
            ;;
        *)
            echo "Usage: $0 [setup|cleanup|capture]"
            echo "  setup   - Create the full lab environment (default)"
            echo "  cleanup - Delete the kind cluster"
            echo "  capture - Re-run packet capture only"
            exit 1
            ;;
    esac
}

# Run main function
main "$@" 