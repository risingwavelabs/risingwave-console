#!/bin/bash

# This script manages Apache License v2 headers in Go source files
# It can add, update, or remove license headers as needed

# Default excluded directories
EXCLUDED_DIRS=("./ee/*")

YEAR=$(date +%Y)

LICENSE_HEADER="// Copyright $YEAR RisingWave Labs
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License."

# Function to detect if a file has a license header
has_license_header() {
    local file="$1"
    # Look for any copyright line with RisingWave Labs
    grep -q "Copyright.*RisingWave Labs" "$file"
}

# Function to remove existing license header
remove_license_header() {
    local file="$1"
    local temp_file=$(mktemp)
    
    # Use awk to remove the license block
    awk '
    BEGIN { in_license = 0 }
    /^\/\/ Copyright.*RisingWave Labs/ { in_license = 1; next }
    /^\/\/ See the License for the specific language governing permissions and/ { in_license = 0; next }
    in_license == 0 { print }
    ' "$file" > "$temp_file"
    
    mv "$temp_file" "$file"
}

# Function to update year in existing license header
update_license_year() {
    local file="$1"
    local temp_file=$(mktemp)
    
    # Replace the year in the copyright line
    sed "s/Copyright [0-9]\{4\} RisingWave Labs/Copyright $YEAR RisingWave Labs/" "$file" > "$temp_file"
    mv "$temp_file" "$file"
}

# Build the find command with exclusions
FIND_CMD="find . -type f -name \"*.go\""
for dir in "${EXCLUDED_DIRS[@]}"; do
    FIND_CMD="$FIND_CMD -not -path \"$dir\""
done

# Process all Go source files
eval "$FIND_CMD" | while read -r file; do
    if has_license_header "$file"; then
        # Update the year in existing license
        echo "Updating license year in $file"
        update_license_year "$file"
    else
        # Add new license header
        echo "Adding license header to $file"
        temp_file=$(mktemp)
        echo "$LICENSE_HEADER" > "$temp_file"
        echo "" >> "$temp_file"
        cat "$file" >> "$temp_file"
        mv "$temp_file" "$file"
    fi
done

echo "License headers have been processed in all Go source files (excluding specified directories)."
