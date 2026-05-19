#!/bin/bash

# Database configuration
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-fintech_db}
MIGRATIONS_PATH="migrations"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to run migrations
run_migrations() {
    echo -e "${GREEN}Running migrations...${NC}"
    
    for file in $(ls ${MIGRATIONS_PATH}/*.up.sql | sort); do
        echo -e "${YELLOW}Executing: $(basename $file)${NC}"
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $file
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ $(basename $file) completed${NC}"
        else
            echo -e "${RED}✗ Error in $(basename $file)${NC}"
            exit 1
        fi
    done
    
    echo -e "${GREEN}All migrations completed successfully!${NC}"
}

# Function to rollback migrations
rollback_migrations() {
    echo -e "${YELLOW}Rolling back migrations...${NC}"
    
    for file in $(ls ${MIGRATIONS_PATH}/*.down.sql | sort -r); do
        echo -e "${YELLOW}Reverting: $(basename $file)${NC}"
        PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $file
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✓ $(basename $file) reverted${NC}"
        else
            echo -e "${RED}✗ Error reverting $(basename $file)${NC}"
            exit 1
        fi
    done
    
    echo -e "${GREEN}Rollback completed!${NC}"
}

# Main script logic
case "$1" in
    up)
        run_migrations
        ;;
    down)
        rollback_migrations
        ;;
    reset)
        rollback_migrations
        run_migrations
        ;;
    *)
        echo "Usage: $0 {up|down|reset}"
        echo "  up    - Run all migrations"
        echo "  down  - Rollback all migrations"
        echo "  reset - Rollback and re-run all migrations"
        exit 1
        ;;
esac