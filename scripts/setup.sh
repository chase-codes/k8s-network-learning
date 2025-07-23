#!/bin/bash

# NetLab Setup Script
# Validates toolchain presence and provides installation guidance

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    NetLab Setup Validation                  ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_section() {
    echo -e "\n${BLUE}▶ $1${NC}"
}

check_command() {
    local cmd=$1
    local name=$2
    local required=$3
    local install_hint=$4

    if command -v "$cmd" &> /dev/null; then
        local version=$($cmd --version 2>&1 | head -n1 || echo "version unknown")
        echo -e "  ${GREEN}✓${NC} $name: $version"
        return 0
    else
        if [ "$required" = "true" ]; then
            echo -e "  ${RED}✗${NC} $name: ${RED}REQUIRED - Not found${NC}"
            if [ -n "$install_hint" ]; then
                echo -e "    💡 Install: $install_hint"
            fi
            return 1
        else
            echo -e "  ${YELLOW}⚠${NC} $name: ${YELLOW}Optional - Not found${NC}"
            if [ -n "$install_hint" ]; then
                echo -e "    💡 Install: $install_hint"
            fi
            return 0
        fi
    fi
}

check_go_version() {
    if command -v go &> /dev/null; then
        local version=$(go version | cut -d' ' -f3 | sed 's/go//')
        local major=$(echo $version | cut -d'.' -f1)
        local minor=$(echo $version | cut -d'.' -f2)
        
        if [ "$major" -gt 1 ] || ([ "$major" -eq 1 ] && [ "$minor" -ge 21 ]); then
            echo -e "  ${GREEN}✓${NC} Go: go$version (meets requirement ≥1.21)"
            return 0
        else
            echo -e "  ${RED}✗${NC} Go: go$version ${RED}(requires ≥1.21)${NC}"
            echo -e "    💡 Update: https://golang.org/dl/"
            return 1
        fi
    else
        echo -e "  ${RED}✗${NC} Go: ${RED}REQUIRED - Not found${NC}"
        echo -e "    💡 Install: https://golang.org/dl/"
        return 1
    fi
}

# Main validation
main() {
    print_header
    
    local all_good=true
    local missing_optional=()

    print_section "Core Requirements"
    
    # Check Go version specifically
    if ! check_go_version; then
        all_good=false
    fi

    print_section "Optional Tools for Advanced Modules"
    
    # Docker
    if ! check_command "docker" "Docker" "false" "https://docs.docker.com/get-docker/"; then
        missing_optional+=("Docker")
    fi

    # kubectl
    if ! check_command "kubectl" "kubectl" "false" "https://kubernetes.io/docs/tasks/tools/"; then
        missing_optional+=("kubectl")
    fi

    # kind
    if ! check_command "kind" "kind" "false" "go install sigs.k8s.io/kind@latest"; then
        missing_optional+=("kind")
    fi

    # Network tools (platform-specific)
    case "$(uname)" in
        "Darwin")
            if ! check_command "tcpdump" "tcpdump" "false" "brew install tcpdump"; then
                missing_optional+=("tcpdump")
            fi
            ;;
        "Linux")
            if ! check_command "tcpdump" "tcpdump" "false" "sudo apt-get install tcpdump"; then
                missing_optional+=("tcpdump")
            fi
            if ! check_command "ip" "iproute2" "false" "sudo apt-get install iproute2"; then
                missing_optional+=("iproute2")
            fi
            if ! check_command "iptables" "iptables" "false" "sudo apt-get install iptables"; then
                missing_optional+=("iptables")
            fi
            ;;
    esac

    print_section "Summary"
    
    if $all_good && [ ${#missing_optional[@]} -eq 0 ]; then
        echo -e "${GREEN}🎉 Perfect! All tools are installed and ready.${NC}"
        echo -e "${GREEN}✅ NetLab is ready to run with full functionality.${NC}"
    elif $all_good; then
        echo -e "${GREEN}✅ Core requirements met! NetLab is ready to run.${NC}"
        if [ ${#missing_optional[@]} -gt 0 ]; then
            echo -e "${YELLOW}⚠️  Optional tools missing: $(IFS=', '; echo "${missing_optional[*]}")${NC}"
            echo -e "   ${YELLOW}Some advanced modules may have limited functionality.${NC}"
        fi
    else
        echo -e "${RED}❌ Missing required dependencies!${NC}"
        echo -e "   ${RED}Please install the required tools and run this script again.${NC}"
        exit 1
    fi

    print_section "Next Steps"
    echo "1. Build NetLab:        make build"
    echo "2. Start learning:      make start"
    echo "3. Run diagnostics:     make doctor"
    echo "4. See all commands:    make help"
    
    echo -e "\n${BLUE}Happy learning! 🚀${NC}"
}

# Check if script is being run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi 