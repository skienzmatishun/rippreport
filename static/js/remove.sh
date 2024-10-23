#!/bin/bash

# Check if a file is passed as an argument
if [ -z "$1" ]; then
  echo "Error: No file specified."
  echo "Usage: $0 <filename>"
  exit 1
fi

# Check if the file exists
if [ ! -f "$1" ]; then
  echo "Error: File '$1' not found."
  exit 1
fi

# Define the URL to remove
URL="https://web.archive.org/web/20240919181358/"

# Count occurrences of the URL before making changes
occurrences=$(grep -o "$URL" "$1" | wc -l)

if [ "$occurrences" -eq 0 ]; then
  echo "No instances of the URL found in '$1'."
  exit 0
else
  echo "$occurrences instances of the URL found in '$1'."
fi

# Use sed to remove all occurrences of the URL from the file
sed -i '' "s|$URL||g" "$1"

# Verify if the URL was successfully removed
new_occurrences=$(grep -o "$URL" "$1" | wc -l)

if [ "$new_occurrences" -eq 0 ]; then
  echo "Successfully removed all instances of the URL from '$1'."
else
  echo "Error: Some instances of the URL remain in '$1'."
fi
