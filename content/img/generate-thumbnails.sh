#!/bin/bash

# Loop through each image file in the current directory
for img in *.{jpg,jpeg,png,gif}; do
    # Check if the file exists
    if [ -f "$img" ]; then
        # Extract the filename without the extension
        filename="${img%.*}"
        # Extract the file extension
        extension="${img##*.}"
        # Create the output filename
        output="${filename}-thumb.${extension}"
        # Resize the image to 60% of its original size and save it as the output file
        convert "$img" -resize 60% "$output"
        echo "Processed $img -> $output"
    fi
done
