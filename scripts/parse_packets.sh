#!/bin/bash

# OSI Model Packet Parser Script
# This script parses packet capture files and extracts headers for each OSI layer

set -e

# Configuration
PCAP_FILE="${1:-modules/01-osi-model/assets/https-nginx.pcap}"
OUTPUT_DIR="modules/01-osi-model/assets"
JSON_OUTPUT="$OUTPUT_DIR/parsed-packets.json"

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

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check prerequisites
check_prerequisites() {
    print_status "Checking prerequisites..."
    
    if ! command_exists tshark; then
        print_error "tshark is required but not installed"
        echo "Install with: brew install wireshark (macOS) or apt-get install tshark (Linux)"
        exit 1
    fi
    
    if [ ! -f "$PCAP_FILE" ]; then
        print_error "Packet capture file not found: $PCAP_FILE"
        echo "Run './scripts/k8s_lab.sh setup' first to generate packet capture"
        exit 1
    fi
    
    print_success "Prerequisites checked"
}

# Function to extract packet information
extract_packet_info() {
    local packet_num=$1
    local output_file="$OUTPUT_DIR/packet-${packet_num}.txt"
    
    print_status "Extracting packet $packet_num information..."
    
    # Extract detailed packet information using tshark
    tshark -r "$PCAP_FILE" -Y "tcp and frame.number==$packet_num" -V > "$output_file" 2>/dev/null || {
        print_error "Failed to extract packet $packet_num"
        return 1
    }
    
    print_success "Packet $packet_num information saved to $output_file"
}

# Function to parse ethernet headers
parse_ethernet() {
    local packet_num=$1
    
    print_status "Parsing Ethernet headers for packet $packet_num..."
    
    # Extract Ethernet frame information
    tshark -r "$PCAP_FILE" -Y "frame.number==$packet_num" -T fields \
        -e eth.dst -e eth.src -e eth.type -e frame.len \
        2>/dev/null | while read dst_mac src_mac eth_type frame_len; do
        
        cat <<EOF
{
  "layer": 2,
  "name": "Data Link Layer (Ethernet)",
  "headers": {
    "Destination MAC": "$dst_mac",
    "Source MAC": "$src_mac",
    "EtherType": "$eth_type",
    "Frame Length": "$frame_len bytes"
  },
  "explanation": "Ethernet frame provides MAC-to-MAC delivery within the local network segment."
}
EOF
    done
}

# Function to parse IP headers
parse_ip() {
    local packet_num=$1
    
    print_status "Parsing IP headers for packet $packet_num..."
    
    # Extract IP header information
    tshark -r "$PCAP_FILE" -Y "frame.number==$packet_num and ip" -T fields \
        -e ip.version -e ip.hdr_len -e ip.src -e ip.dst -e ip.proto -e ip.ttl -e ip.len \
        2>/dev/null | while read version hdr_len src_ip dst_ip protocol ttl total_len; do
        
        # Convert protocol number to name
        case $protocol in
            6) proto_name="TCP" ;;
            17) proto_name="UDP" ;;
            1) proto_name="ICMP" ;;
            *) proto_name="Protocol $protocol" ;;
        esac
        
        cat <<EOF
{
  "layer": 3,
  "name": "Network Layer (IP)",
  "headers": {
    "Version": "$version (IPv4)",
    "Header Length": "$hdr_len bytes",
    "Source IP": "$src_ip",
    "Destination IP": "$dst_ip",
    "Protocol": "$protocol ($proto_name)",
    "TTL": "$ttl",
    "Total Length": "$total_len bytes"
  },
  "explanation": "IP header provides logical addressing for routing packets across networks."
}
EOF
    done
}

# Function to parse TCP headers
parse_tcp() {
    local packet_num=$1
    
    print_status "Parsing TCP headers for packet $packet_num..."
    
    # Extract TCP header information
    tshark -r "$PCAP_FILE" -Y "frame.number==$packet_num and tcp" -T fields \
        -e tcp.srcport -e tcp.dstport -e tcp.seq -e tcp.ack -e tcp.flags -e tcp.window_size_value -e tcp.checksum \
        2>/dev/null | while read src_port dst_port seq_num ack_num flags window checksum; do
        
        # Parse TCP flags
        flag_names=""
        if [[ $flags == *"0x002"* ]]; then flag_names="SYN "; fi
        if [[ $flags == *"0x010"* ]]; then flag_names+="ACK "; fi
        if [[ $flags == *"0x001"* ]]; then flag_names+="FIN "; fi
        if [[ $flags == *"0x004"* ]]; then flag_names+="RST "; fi
        
        cat <<EOF
{
  "layer": 4,
  "name": "Transport Layer (TCP)",
  "headers": {
    "Source Port": "$src_port",
    "Destination Port": "$dst_port",
    "Sequence Number": "$seq_num",
    "Acknowledgment": "$ack_num",
    "Flags": "${flag_names:-$flags}",
    "Window Size": "$window",
    "Checksum": "$checksum"
  },
  "explanation": "TCP header provides reliable, connection-oriented communication with flow control."
}
EOF
    done
}

# Function to parse HTTP headers
parse_http() {
    local packet_num=$1
    
    print_status "Parsing HTTP headers for packet $packet_num..."
    
    # Check if packet contains HTTP data
    if tshark -r "$PCAP_FILE" -Y "frame.number==$packet_num and http" -T fields -e http.request.method 2>/dev/null | grep -q "."; then
        
        # Extract HTTP request information
        tshark -r "$PCAP_FILE" -Y "frame.number==$packet_num and http" -T fields \
            -e http.request.method -e http.request.uri -e http.request.version -e http.host -e http.user_agent \
            2>/dev/null | while read method uri version host user_agent; do
            
            cat <<EOF
{
  "layer": 7,
  "name": "Application Layer (HTTP)",
  "headers": {
    "Method": "$method",
    "URI": "$uri",
    "HTTP Version": "$version",
    "Host": "$host",
    "User-Agent": "$user_agent"
  },
  "explanation": "HTTP application layer contains the actual web request that users and applications care about."
}
EOF
        done
    else
        # No HTTP data in this packet
        return 1
    fi
}

# Function to create comprehensive packet analysis
create_packet_analysis() {
    print_status "Creating comprehensive packet analysis..."
    
    # Get total packet count
    local total_packets=$(tshark -r "$PCAP_FILE" -q -z io,stat,0 | grep packets | tail -1 | awk '{print $2}')
    print_status "Total packets in capture: $total_packets"
    
    # Find the first HTTP request packet
    local http_packet=$(tshark -r "$PCAP_FILE" -Y "http.request" -T fields -e frame.number | head -1)
    
    if [ -n "$http_packet" ]; then
        print_status "Found HTTP request in packet $http_packet"
        
        # Create JSON analysis for the HTTP packet
        cat > "$JSON_OUTPUT" <<EOF
{
  "capture_info": {
    "file": "$PCAP_FILE",
    "total_packets": $total_packets,
    "analysis_packet": $http_packet,
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  },
  "osi_layers": [
EOF
        
        # Parse each layer and add to JSON
        local first=true
        
        # Layer 2 - Ethernet
        if ethernet_data=$(parse_ethernet "$http_packet" 2>/dev/null); then
            if [ "$first" = false ]; then echo "," >> "$JSON_OUTPUT"; fi
            echo "$ethernet_data" >> "$JSON_OUTPUT"
            first=false
        fi
        
        # Layer 3 - IP
        if ip_data=$(parse_ip "$http_packet" 2>/dev/null); then
            if [ "$first" = false ]; then echo "," >> "$JSON_OUTPUT"; fi
            echo "$ip_data" >> "$JSON_OUTPUT"
            first=false
        fi
        
        # Layer 4 - TCP
        if tcp_data=$(parse_tcp "$http_packet" 2>/dev/null); then
            if [ "$first" = false ]; then echo "," >> "$JSON_OUTPUT"; fi
            echo "$tcp_data" >> "$JSON_OUTPUT"
            first=false
        fi
        
        # Layer 7 - HTTP
        if http_data=$(parse_http "$http_packet" 2>/dev/null); then
            if [ "$first" = false ]; then echo "," >> "$JSON_OUTPUT"; fi
            echo "$http_data" >> "$JSON_OUTPUT"
            first=false
        fi
        
        cat >> "$JSON_OUTPUT" <<EOF
  ]
}
EOF
        
        print_success "Packet analysis saved to $JSON_OUTPUT"
    else
        print_error "No HTTP packets found in capture"
        return 1
    fi
}

# Function to create text summary
create_text_summary() {
    print_status "Creating text summary..."
    
    local summary_file="$OUTPUT_DIR/packet-summary.txt"
    
    cat > "$summary_file" <<EOF
OSI Model Packet Analysis Summary
=================================

Capture File: $PCAP_FILE
Analysis Date: $(date)

This analysis shows how a simple HTTP request from a busybox Pod to nginx
demonstrates all the OSI model layers in action within a Kubernetes cluster.

OSI Layer Breakdown:
-------------------

Layer 7 - Application (HTTP):
  The actual web request that users care about. Contains HTTP methods,
  headers, and the request for the nginx welcome page.

Layer 4 - Transport (TCP):
  Provides reliable delivery with port numbers, sequence numbers, and
  connection management. Uses port 80 for HTTP.

Layer 3 - Network (IP):
  Routes packets between Pod IPs within the Kubernetes cluster using
  the cluster's internal network (usually 10.244.x.x range).

Layer 2 - Data Link (Ethernet):
  Handles MAC-to-MAC delivery within the kind container's virtual
  network using container network interfaces.

Layer 1 - Physical:
  Virtual network interfaces (veth pairs) that connect containers
  to the kind cluster's virtual bridge network.

Kubernetes Networking Context:
-----------------------------
- Pod-to-Pod communication uses the cluster's overlay network
- Service discovery resolves 'nginx' to the actual Pod IP
- CNI (Container Network Interface) manages IP allocation and routing
- kube-proxy handles service load balancing and port translation

To explore further:
- Run 'kubectl get pods -o wide' to see Pod IPs
- Use 'kubectl exec' to run commands inside the containers
- Examine service endpoints with 'kubectl get endpoints'
EOF
    
    print_success "Text summary saved to $summary_file"
}

# Function to extract raw packet data
extract_raw_data() {
    print_status "Extracting raw packet data..."
    
    # Get the first few packets as hex dump
    local raw_file="$OUTPUT_DIR/raw-packets.txt"
    
    tshark -r "$PCAP_FILE" -c 5 -x > "$raw_file" 2>/dev/null
    
    print_success "Raw packet data saved to $raw_file"
}

# Function to validate analysis
validate_analysis() {
    print_status "Validating packet analysis..."
    
    if [ -f "$JSON_OUTPUT" ] && [ -s "$JSON_OUTPUT" ]; then
        # Check if JSON is valid
        if command_exists python3; then
            python3 -m json.tool "$JSON_OUTPUT" >/dev/null 2>&1 && {
                print_success "JSON analysis is valid"
            } || {
                print_error "JSON analysis is invalid"
                return 1
            }
        fi
    else
        print_error "JSON analysis file is missing or empty"
        return 1
    fi
}

# Main function
main() {
    echo "🔍 NetLab OSI Model - Packet Parser"
    echo "==================================="
    echo ""
    
    case "${1:-parse}" in
        "parse")
            check_prerequisites
            create_packet_analysis
            create_text_summary
            extract_raw_data
            validate_analysis
            
            print_success "Packet parsing completed!"
            echo ""
            echo "📄 Generated Files:"
            echo "  • $JSON_OUTPUT - Structured packet analysis"
            echo "  • $OUTPUT_DIR/packet-summary.txt - Human-readable summary"
            echo "  • $OUTPUT_DIR/raw-packets.txt - Raw packet dumps"
            echo ""
            ;;
        "validate")
            validate_analysis
            ;;
        *)
            echo "Usage: $0 [parse|validate]"
            echo "  parse    - Parse packet capture and create analysis files (default)"
            echo "  validate - Validate existing analysis files"
            exit 1
            ;;
    esac
}

# Run main function
main "$@" 