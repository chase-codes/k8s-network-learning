# OSI Model Module

## Learning Objectives

By the end of this module, you will understand:

- The seven layers of the OSI model and their functions
- How data flows through each layer during network communication
- Real-world protocols and technologies that operate at each layer
- The relationship between OSI layers and modern networking
- How the OSI model applies to Kubernetes networking concepts

## Prerequisites

- Basic understanding of computer networks
- Familiarity with common internet protocols (HTTP, TCP/IP)

## Module Structure

This module presents the OSI model through:

1. **Interactive TUI Interface**: Navigate through each layer with detailed explanations
2. **ASCII Diagrams**: Visual representation of the layer stack
3. **Real-World Examples**: How web browsing maps to the OSI layers
4. **Practice Questions**: Test your understanding
5. **Kubernetes Context**: Connect OSI concepts to container networking

## Key Concepts Covered

### Layer 7 - Application Layer
- HTTP/HTTPS, FTP, SMTP, DNS
- User interface and application services

### Layer 6 - Presentation Layer  
- Data encryption, compression, formatting
- SSL/TLS, image/video codecs

### Layer 5 - Session Layer
- Session management and control
- NetBIOS, RPC, SQL sessions

### Layer 4 - Transport Layer
- TCP (reliable) vs UDP (fast)
- Port numbers and end-to-end delivery

### Layer 3 - Network Layer
- IP addressing and routing
- Packet forwarding across networks

### Layer 2 - Data Link Layer
- Frame formatting and MAC addresses
- Switch operations and error detection

### Layer 1 - Physical Layer
- Physical transmission media
- Electrical signals and hardware

## Memory Aid

**"Please Do Not Throw Sausage Pizza Away"**
- Physical
- Data Link
- Network
- Transport
- Session
- Presentation
- Application

## Next Steps

After completing this module, proceed to:
- `netlab module 02-tcp-ip` - Deep dive into the TCP/IP protocol stack
- Or return to the main menu with `netlab start`

## Usage

Run this module directly with:
```bash
netlab module 01-osi-model
```

Navigate using:
- **↑/↓ arrows**: Scroll through content
- **q/Ctrl+C**: Exit back to main menu 