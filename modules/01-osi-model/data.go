package osimodel

// OSILayer represents a single layer of the OSI model
type OSILayer struct {
	Number      int
	Name        string
	ShortName   string
	Description string
	Function    string
	Protocols   []string
	Analogy     string
	HeaderType  string
	CLITools    []string
	ExternalDoc string
	Examples    []string
	KeyConcepts []string
}

// GetOSILayers returns all seven OSI layers with detailed information
func GetOSILayers() []OSILayer {
	return []OSILayer{
		{
			Number:      7,
			Name:        "Application Layer",
			ShortName:   "Application",
			Description: "Provides network services directly to end-user applications. This is where human-computer interaction happens through network-aware applications.",
			Function:    "Interface between applications and the network",
			Protocols:   []string{"HTTP", "HTTPS", "FTP", "SMTP", "POP3", "IMAP", "DNS", "DHCP", "SSH", "Telnet"},
			Analogy:     "Like the post office window where you interact with postal services",
			HeaderType:  "Application-specific headers (HTTP headers, email headers)",
			CLITools:    []string{"curl", "wget", "dig", "nslookup", "ssh", "telnet", "ftp"},
			ExternalDoc: "https://developer.mozilla.org/en-US/docs/Web/HTTP",
			Examples: []string{
				"Web browsers sending HTTP requests",
				"Email clients using SMTP/IMAP",
				"File transfer with FTP/SFTP",
				"DNS resolution queries",
			},
			KeyConcepts: []string{
				"User interface to network services",
				"Protocol-specific formatting",
				"Application data handling",
				"Service identification",
			},
		},
		{
			Number:      6,
			Name:        "Presentation Layer",
			ShortName:   "Presentation",
			Description: "Handles data translation, encryption, compression, and formatting. Ensures data sent by one system can be read by another.",
			Function:    "Data translation, encryption, and compression",
			Protocols:   []string{"SSL/TLS", "JPEG", "MPEG", "GIF", "PNG", "ASCII", "EBCDIC", "MIME"},
			Analogy:     "Like a translator who converts between languages and encrypts messages",
			HeaderType:  "Encryption headers, compression metadata",
			CLITools:    []string{"openssl", "gpg", "base64", "gzip", "tar"},
			ExternalDoc: "https://tools.ietf.org/html/rfc5246",
			Examples: []string{
				"SSL/TLS encryption for HTTPS",
				"Image compression (JPEG, PNG)",
				"Video encoding (MPEG, H.264)",
				"Character encoding (UTF-8, ASCII)",
			},
			KeyConcepts: []string{
				"Data encryption and decryption",
				"Compression and decompression",
				"Character set conversion",
				"Data format translation",
			},
		},
		{
			Number:      5,
			Name:        "Session Layer",
			ShortName:   "Session",
			Description: "Manages sessions between applications. Establishes, maintains, synchronizes, and terminates communication sessions.",
			Function:    "Session establishment, management, and termination",
			Protocols:   []string{"NetBIOS", "RPC", "PPTP", "L2TP", "SQL sessions", "NFS"},
			Analogy:     "Like a meeting coordinator who schedules, manages, and ends meetings",
			HeaderType:  "Session management headers, checkpoint markers",
			CLITools:    []string{"netstat", "ss", "rpcinfo", "showmount"},
			ExternalDoc: "https://tools.ietf.org/html/rfc1001",
			Examples: []string{
				"Database connection sessions",
				"Web application login sessions",
				"Remote procedure calls (RPC)",
				"Network file system sessions",
			},
			KeyConcepts: []string{
				"Session establishment",
				"Synchronization and checkpointing",
				"Session recovery",
				"Connection management",
			},
		},
		{
			Number:      4,
			Name:        "Transport Layer",
			ShortName:   "Transport",
			Description: "Provides reliable data transfer services to upper layers. Handles error detection, flow control, and segmentation.",
			Function:    "End-to-end data delivery and error recovery",
			Protocols:   []string{"TCP", "UDP", "SCTP", "SPX"},
			Analogy:     "Like a delivery service that ensures packages arrive intact and in order",
			HeaderType:  "TCP/UDP headers with ports, sequence numbers, checksums",
			CLITools:    []string{"netstat", "ss", "lsof", "tcpdump", "wireshark", "nmap"},
			ExternalDoc: "https://tools.ietf.org/html/rfc793",
			Examples: []string{
				"TCP reliable web traffic (port 80, 443)",
				"UDP streaming media (DNS port 53)",
				"TCP file transfers (FTP port 21)",
				"UDP gaming traffic",
			},
			KeyConcepts: []string{
				"Port numbers (0-65535)",
				"Reliable vs unreliable delivery",
				"Flow control and congestion control",
				"Segmentation and reassembly",
			},
		},
		{
			Number:      3,
			Name:        "Network Layer",
			ShortName:   "Network",
			Description: "Handles routing of data packets between different networks. Determines the best path for data across multiple networks.",
			Function:    "Routing and logical addressing",
			Protocols:   []string{"IP", "IPv6", "ICMP", "OSPF", "BGP", "RIP", "EIGRP"},
			Analogy:     "Like a GPS system that finds the best route between addresses",
			HeaderType:  "IP headers with source/destination addresses, TTL",
			CLITools:    []string{"ping", "traceroute", "route", "ip", "iptables", "mtr"},
			ExternalDoc: "https://tools.ietf.org/html/rfc791",
			Examples: []string{
				"IP routing between networks",
				"ICMP ping and traceroute",
				"Router forwarding decisions",
				"Subnet communication",
			},
			KeyConcepts: []string{
				"IP addresses (IPv4/IPv6)",
				"Routing tables and algorithms",
				"Subnetting and VLANs",
				"Packet forwarding",
			},
		},
		{
			Number:      2,
			Name:        "Data Link Layer",
			ShortName:   "Data Link",
			Description: "Provides node-to-node data transfer and error detection/correction for the physical layer. Handles MAC addressing.",
			Function:    "Node-to-node delivery and error detection",
			Protocols:   []string{"Ethernet", "Wi-Fi (802.11)", "PPP", "Frame Relay", "ATM"},
			Analogy:     "Like addressing an envelope with the recipient's street address",
			HeaderType:  "Ethernet frames with MAC addresses, frame check sequence",
			CLITools:    []string{"arp", "bridge", "brctl", "iwconfig", "ethtool"},
			ExternalDoc: "https://standards.ieee.org/standard/802_3-2018.html",
			Examples: []string{
				"Ethernet frame transmission",
				"Wi-Fi wireless communication",
				"Switch forwarding decisions",
				"ARP address resolution",
			},
			KeyConcepts: []string{
				"MAC addresses (48-bit hardware)",
				"Frame formatting and CRC",
				"Collision detection (CSMA/CD)",
				"Switch operation",
			},
		},
		{
			Number:      1,
			Name:        "Physical Layer",
			ShortName:   "Physical",
			Description: "Defines the electrical, mechanical, and procedural interface to the physical transmission medium. Raw bit transmission.",
			Function:    "Physical transmission of raw bits",
			Protocols:   []string{"Ethernet cables", "Fiber optic", "Wi-Fi radio", "Bluetooth", "USB"},
			Analogy:     "Like the actual roads and vehicles that carry the mail",
			HeaderType:  "No headers - raw electrical/optical signals",
			CLITools:    []string{"ethtool", "iwlist", "lshw", "dmesg", "lsusb"},
			ExternalDoc: "https://standards.ieee.org/standard/802_3-2018.html",
			Examples: []string{
				"Copper wire electrical signals",
				"Fiber optic light pulses",
				"Radio frequency transmission",
				"Cable specifications (Cat5e, Cat6)",
			},
			KeyConcepts: []string{
				"Electrical signal specifications",
				"Cable types and connectors",
				"Signal encoding and modulation",
				"Physical topology",
			},
		},
	}
}

// GetMnemonics returns popular mnemonic devices for remembering OSI layers
func GetMnemonics() []string {
	return []string{
		"Please Do Not Throw Sausage Pizza Away",
		"All People Seem To Need Data Processing",
		"Please Do Not Tell Secret Passwords Anywhere",
		"Please Do Not Touch Steve's Pet Alligator",
		"Physical Data Networking Transport Session Presentation Application",
	}
}

// GetRealWorldExample returns a detailed example of OSI layers in action
func GetRealWorldExample() map[string]string {
	return map[string]string{
		"title":       "HTTPS Web Request to example.com",
		"description": "When you type https://example.com in your browser and press Enter",
		"layer7":      "Browser formats HTTP request: 'GET / HTTP/1.1\\nHost: example.com'",
		"layer6":      "TLS encrypts the HTTP data and compresses it",
		"layer5":      "Session established between browser and web server",
		"layer4":      "TCP wraps data with port 443, sequence numbers, checksums",
		"layer3":      "IP adds source (your IP) and destination (example.com IP) addresses",
		"layer2":      "Ethernet frame adds your MAC and router's MAC address",
		"layer1":      "Electrical signals sent over network cable or Wi-Fi radio waves",
	}
}

// GetKubernetesContext returns how OSI layers relate to Kubernetes networking
func GetKubernetesContext() map[int]string {
	return map[int]string{
		7: "Ingress controllers handle HTTP/HTTPS traffic routing and load balancing",
		6: "TLS termination at Ingress for HTTPS, cert-manager for certificate management",
		5: "Service sessions, connection pooling in service meshes like Istio",
		4: "Service ports, load balancing, kube-proxy manages port translation",
		3: "Pod IPs, Service IPs, cluster CIDR, CNI manages IP address allocation",
		2: "CNI plugins handle container network interfaces, bridge networks",
		1: "Node network interfaces, physical/virtual network infrastructure",
	}
}
