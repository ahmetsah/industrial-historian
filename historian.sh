#!/bin/bash

# =============================================================================
# Historian Platform Management Script
# =============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OPS_DIR="$SCRIPT_DIR/ops"
COMPOSE_FILE="$OPS_DIR/docker-compose.yml"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Print banner
print_banner() {
    echo -e "${CYAN}"
    echo "╔═══════════════════════════════════════════════════════════════╗"
    echo "║              🏭 HISTORIAN PLATFORM MANAGER                    ║"
    echo "║                   Industrial Data Historian                   ║"
    echo "╚═══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Print usage
usage() {
    echo -e "${BOLD}Usage:${NC} $0 <command>"
    echo ""
    echo -e "${BOLD}Commands:${NC}"
    echo -e "  ${GREEN}start${NC}      Start all services"
    echo -e "  ${RED}stop${NC}       Stop all services"
    echo -e "  ${YELLOW}restart${NC}    Restart all services"
    echo -e "  ${BLUE}status${NC}     Show status of all services"
    echo -e "  ${CYAN}logs${NC}       Show logs (use: logs [service_name])"
    echo -e "  ${CYAN}errors${NC}     Show services with errors"
    echo -e "  ${CYAN}build${NC}      Build all services"
    echo -e "  ${CYAN}rebuild${NC}    Rebuild and restart all services"
    echo ""
    echo -e "${BOLD}Examples:${NC}"
    echo "  $0 start"
    echo "  $0 status"
    echo "  $0 logs engine"
    echo "  $0 errors"
    echo ""
}

# Check prerequisites
check_prereqs() {
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}Error: Docker is not installed${NC}"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${RED}Error: Docker Compose is not installed${NC}"
        exit 1
    fi
    
    if [ ! -f "$COMPOSE_FILE" ]; then
        echo -e "${RED}Error: docker-compose.yml not found at $COMPOSE_FILE${NC}"
        exit 1
    fi
}

# Run docker-compose command
dc() {
    cd "$OPS_DIR"
    if docker compose version &> /dev/null 2>&1; then
        docker compose "$@"
    else
        docker-compose "$@"
    fi
}

# Start all services
start_services() {
    echo -e "${GREEN}▶ Starting Historian Platform...${NC}"
    echo ""
    dc up -d
    echo ""
    echo -e "${GREEN}✅ All services started${NC}"
    echo ""
    show_status
    show_access_info
}

# Stop all services
stop_services() {
    echo -e "${RED}■ Stopping Historian Platform...${NC}"
    echo ""
    dc down
    echo ""
    echo -e "${RED}✅ All services stopped${NC}"
}

# Restart all services
restart_services() {
    echo -e "${YELLOW}↻ Restarting Historian Platform...${NC}"
    echo ""
    dc restart
    echo ""
    echo -e "${GREEN}✅ All services restarted${NC}"
    echo ""
    show_status
}

# Show status of all services
show_status() {
    echo -e "${BLUE}${BOLD}📊 Service Status${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    # Get all containers from our compose project
    containers=$(docker ps -a --filter "name=ops-" --format "{{.Names}}\t{{.Status}}\t{{.Ports}}" 2>/dev/null || true)
    
    if [ -z "$containers" ]; then
        echo -e "${YELLOW}No Historian containers found${NC}"
        return
    fi
    
    printf "${BOLD}%-30s %-25s %s${NC}\n" "SERVICE" "STATUS" "PORTS"
    echo "────────────────────────────────────────────────────────────────────"
    
    while IFS=$'\t' read -r name status ports; do
        # Determine color based on status
        if [[ "$status" == *"Up"* ]]; then
            if [[ "$status" == *"unhealthy"* ]]; then
                status_color="${YELLOW}⚠ ${status}${NC}"
            else
                status_color="${GREEN}✓ ${status}${NC}"
            fi
        else
            status_color="${RED}✗ ${status}${NC}"
        fi
        
        # Clean up service name
        service_name="${name#ops-}"
        
        printf "%-30s %-35b %s\n" "$service_name" "$status_color" "$ports"
    done <<< "$containers"
    
    echo ""
}

# Show errors and unhealthy services
show_errors() {
    echo -e "${RED}${BOLD}🔴 Services with Errors${NC}"
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    has_errors=false
    
    # Check for exited containers
    exited=$(docker ps -a --filter "name=ops-" --filter "status=exited" --format "{{.Names}}" 2>/dev/null || true)
    
    if [ -n "$exited" ]; then
        has_errors=true
        echo -e "${RED}Exited Containers:${NC}"
        for container in $exited; do
            echo -e "  ${RED}✗${NC} $container"
            echo -e "  ${YELLOW}Last 10 log lines:${NC}"
            docker logs --tail 10 "$container" 2>&1 | sed 's/^/    /'
            echo ""
        done
    fi
    
    # Check for unhealthy containers
    unhealthy=$(docker ps --filter "name=ops-" --filter "health=unhealthy" --format "{{.Names}}" 2>/dev/null || true)
    
    if [ -n "$unhealthy" ]; then
        has_errors=true
        echo -e "${YELLOW}Unhealthy Containers:${NC}"
        for container in $unhealthy; do
            echo -e "  ${YELLOW}⚠${NC} $container"
            echo -e "  ${YELLOW}Last 10 log lines:${NC}"
            docker logs --tail 10 "$container" 2>&1 | sed 's/^/    /'
            echo ""
        done
    fi
    
    # Check for containers with restart loops
    restarting=$(docker ps --filter "name=ops-" --filter "status=restarting" --format "{{.Names}}" 2>/dev/null || true)
    
    if [ -n "$restarting" ]; then
        has_errors=true
        echo -e "${RED}Restarting Containers (possible crash loop):${NC}"
        for container in $restarting; do
            echo -e "  ${RED}↻${NC} $container"
            echo -e "  ${YELLOW}Last 10 log lines:${NC}"
            docker logs --tail 10 "$container" 2>&1 | sed 's/^/    /'
            echo ""
        done
    fi
    
    if [ "$has_errors" = false ]; then
        echo -e "${GREEN}✅ No errors found! All services are healthy.${NC}"
    fi
    echo ""
}

# Show logs for a service
show_logs() {
    local service="$1"
    
    if [ -z "$service" ]; then
        echo -e "${CYAN}Showing logs for all services (Ctrl+C to exit)...${NC}"
        dc logs -f --tail 50
    else
        container="ops-$service"
        if docker ps -a --format "{{.Names}}" | grep -q "^$container$"; then
            echo -e "${CYAN}Showing logs for $service (Ctrl+C to exit)...${NC}"
            docker logs -f --tail 100 "$container"
        else
            echo -e "${RED}Error: Service '$service' not found${NC}"
            echo ""
            echo "Available services:"
            docker ps -a --filter "name=ops-" --format "  - {{.Names}}" | sed 's/ops-//'
        fi
    fi
}

# Build all services
build_services() {
    echo -e "${CYAN}🔨 Building all services...${NC}"
    echo ""
    dc build
    echo ""
    echo -e "${GREEN}✅ Build complete${NC}"
}

# Rebuild and restart
rebuild_services() {
    echo -e "${CYAN}🔨 Rebuilding and restarting all services...${NC}"
    echo ""
    dc build --no-cache
    dc up -d --force-recreate
    echo ""
    echo -e "${GREEN}✅ Rebuild complete${NC}"
    echo ""
    show_status
}

# Show access information
show_access_info() {
    echo -e "${CYAN}${BOLD}🌐 Access Information${NC}"
    echo -e "${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "  ${BOLD}Dashboard:${NC}        http://localhost:3000"
    echo -e "  ${BOLD}Config UI:${NC}        http://localhost:3001"
    echo -e "  ${BOLD}Engine API:${NC}       http://localhost:8081"
    echo -e "  ${BOLD}Config API:${NC}       http://localhost:8090"
    echo -e "  ${BOLD}Auth API:${NC}         http://localhost:8080"
    echo -e "  ${BOLD}Alarm API:${NC}        http://localhost:8083"
    echo -e "  ${BOLD}Audit API:${NC}        http://localhost:8082"
    echo -e "  ${BOLD}NATS:${NC}             http://localhost:8222 (monitoring)"
    echo -e "  ${BOLD}MinIO Console:${NC}    http://localhost:9001"
    echo -e "  ${BOLD}PgAdmin:${NC}          http://localhost:5050"
    echo ""
}

# Main
main() {
    print_banner
    check_prereqs
    
    case "${1:-}" in
        start)
            start_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services
            ;;
        status)
            show_status
            show_access_info
            ;;
        logs)
            show_logs "${2:-}"
            ;;
        errors)
            show_errors
            ;;
        build)
            build_services
            ;;
        rebuild)
            rebuild_services
            ;;
        *)
            usage
            exit 1
            ;;
    esac
}

main "$@"
