# OSI Model Module

## Overview

This comprehensive module teaches the OSI (Open Systems Interconnection) model through two interactive experiences:

1. **Interactive Layer Explorer**: Navigate through each OSI layer with detailed explanations, protocols, and real-world context
2. **Hands-On Packet Lab**: Analyze real network traffic from a Kubernetes cluster to see OSI layers in action

## Learning Objectives

By the end of this module, you will understand:

- The seven layers of the OSI model and their specific functions
- How data flows through each layer during network communication
- Real-world protocols and technologies that operate at each layer
- How to use CLI tools to inspect network traffic at different layers
- The relationship between OSI layers and Kubernetes networking
- How to analyze packet captures to understand network behavior

## Prerequisites

### Basic Knowledge
- Basic understanding of computer networks
- Familiarity with common internet protocols (HTTP, TCP/IP)
- Basic command line usage

### For Packet Lab (Optional)
- Docker installed and running
- `kind` (Kubernetes in Docker)
- `kubectl` (Kubernetes CLI)
- `tcpdump` and `tshark` (packet analysis tools)

## Module Structure

### Phase 1: Interactive Layer Explorer

Navigate through the OSI layers with:
- **Detailed explanations** of each layer's function
- **Common protocols** used at each layer
- **Real-world analogies** to understand concepts
- **CLI tools** for inspecting each layer
- **Kubernetes context** showing how each layer applies to container networking
- **Memory aids** to help remember the layer order

**Navigation:**
- `↑/↓` or `j/k` - Navigate between layers
- `m` - Show mnemonic devices for remembering layers
- `v` - Advance to packet analysis lab
- `q` - Return to main menu

### Phase 2: Hands-On Packet Lab

Experience the OSI layers through real packet analysis:
- **Live Kubernetes cluster** using kind
- **Real HTTP traffic** between nginx and busybox pods
- **Packet capture** with tcpdump
- **Layer-by-layer walkthrough** of actual network data
- **Header analysis** showing each layer's contribution

## Key Concepts Covered

### Layer 7 - Application Layer
- **Protocols**: HTTP/HTTPS, FTP, SMTP, DNS, SSH
- **Function**: User interface and network services
- **CLI Tools**: `curl`, `wget`, `dig`, `nslookup`
- **Kubernetes**: Ingress controllers, application routing

### Layer 6 - Presentation Layer
- **Protocols**: SSL/TLS, JPEG, MPEG, encryption
- **Function**: Data translation, encryption, compression
- **CLI Tools**: `openssl`, `gpg`, `base64`
- **Kubernetes**: TLS termination, certificate management

### Layer 5 - Session Layer
- **Protocols**: NetBIOS, RPC, SQL sessions
- **Function**: Session establishment and management
- **CLI Tools**: `netstat`, `ss`, `rpcinfo`
- **Kubernetes**: Service sessions, connection pooling

### Layer 4 - Transport Layer
- **Protocols**: TCP, UDP
- **Function**: End-to-end delivery, port management
- **CLI Tools**: `netstat`, `ss`, `tcpdump`, `nmap`
- **Kubernetes**: Service ports, load balancing

### Layer 3 - Network Layer
- **Protocols**: IP, ICMP, routing protocols
- **Function**: Routing and logical addressing
- **CLI Tools**: `ping`, `traceroute`, `route`, `ip`
- **Kubernetes**: Pod IPs, Service IPs, CNI

### Layer 2 - Data Link Layer
- **Protocols**: Ethernet, Wi-Fi, MAC addressing
- **Function**: Node-to-node delivery, error detection
- **CLI Tools**: `arp`, `bridge`, `iwconfig`
- **Kubernetes**: CNI plugins, container interfaces

### Layer 1 - Physical Layer
- **Protocols**: Ethernet cables, fiber, radio
- **Function**: Physical transmission of bits
- **CLI Tools**: `ethtool`, `lshw`, `iwlist`
- **Kubernetes**: Node network interfaces, infrastructure

## Memory Aids

**Popular Mnemonics:**
1. **"Please Do Not Throw Sausage Pizza Away"** (bottom-up)
2. **"All People Seem To Need Data Processing"** (top-down)
3. **"Please Do Not Tell Secret Passwords Anywhere"**

## Lab Setup and Usage

### Quick Start
```bash
# Run the interactive module
netlab module 01-osi-model

# Within the module:
# - Navigate with ↑/↓ keys
# - Press 'm' for mnemonics
# - Press 'v' for packet lab
```

### Manual Lab Setup
```bash
# Set up the Kubernetes lab environment
./scripts/k8s_lab.sh setup

# Parse packet captures (after lab setup)
./scripts/parse_packets.sh parse

# Run the module with lab data
netlab module 01-osi-model
```

### Lab Management
```bash
# Clean up the lab environment
./scripts/k8s_lab.sh cleanup

# Re-run packet capture only
./scripts/k8s_lab.sh capture

# Validate parsed data
./scripts/parse_packets.sh validate
```

## Generated Files

The lab creates several files for analysis:

### Assets Directory (`modules/01-osi-model/assets/`)
- **`https-nginx.pcap`** - Raw packet capture from Kubernetes traffic
- **`parsed-packets.json`** - Structured analysis of OSI layers
- **`packet-summary.txt`** - Human-readable analysis summary
- **`raw-packets.txt`** - Hex dumps of raw packet data
- **`osi-diagram.txt`** - ASCII art OSI model diagram

### Packet Analysis Data
The lab analyzes real HTTP traffic showing:
- **Ethernet frames** with MAC addresses
- **IP headers** with Pod IPs and routing
- **TCP segments** with ports and connection management
- **HTTP requests** with application-layer data

## Troubleshooting

### Common Issues

**Docker not running:**
```bash
# Start Docker Desktop (macOS) or Docker daemon (Linux)
# Verify with: docker info
```

**Missing dependencies:**
```bash
# macOS
brew install docker kind kubernetes-cli wireshark

# Linux (Ubuntu/Debian)
sudo apt-get update
sudo apt-get install docker.io kubectl tcpdump tshark

# Install kind separately
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

**Cluster creation fails:**
```bash
# Check Docker resources and try again
kind delete cluster --name netlab-osi
./scripts/k8s_lab.sh setup
```

**No packets captured:**
```bash
# Re-run packet capture
./scripts/k8s_lab.sh capture
```

## Kubernetes Networking Context

This module provides the foundation for understanding Kubernetes networking:

- **Pod networking** uses layers 2-4 for container communication
- **Service discovery** operates across layers 3-7 for application routing
- **Ingress controllers** primarily work at layer 7 for HTTP/HTTPS
- **Network policies** function at layers 3-4 for security
- **CNI plugins** handle layers 2-3 for container network interfaces

## Next Steps

After completing this module:

1. **Deep dive into TCP/IP**: `netlab module 02-tcp-ip`
2. **Explore subnetting**: `netlab module 03-subnetting`
3. **Learn Kubernetes networking**: `netlab module 05-k8s-networking`
4. **Return to main menu**: `netlab start`

## Educational Value

This module bridges theoretical networking concepts with practical Kubernetes implementation, showing students:

- **How abstract OSI layers map to real network protocols**
- **Why understanding layers helps with troubleshooting**
- **How Kubernetes leverages these fundamental concepts**
- **What tools to use for network analysis at each layer**
- **How container networking builds on traditional networking** 