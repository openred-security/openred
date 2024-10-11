#!/bin/bash

# Source folder
SOURCE_FOLDER="openred"

# Destination file
DEST_FILE="file_contents.txt"

# Clear the destination file if it already exists
> "$DEST_FILE"

# Check if the source folder exists
if [[ ! -d "$SOURCE_FOLDER" ]]; then
  echo "The folder $SOURCE_FOLDER does not exist."
  exit 1
fi

echo "Processing files in folder: $SOURCE_FOLDER"

# Recursive function to process files
process_files() {
  for FILE in "$1"/*; do
    if [[ -f "$FILE" ]]; then
      # Skip LICENSE files and .sum files
      if [[ "$(basename "$FILE")" == "LICENSE" ]] || [[ "$FILE" == *.sum ]]; then
        echo "Skipping file: $(basename "$FILE")"
        continue
      fi

      # Write the filename as a header
      echo "Processing file: $(basename "$FILE")"
      echo "File: $(basename "$FILE")" >> "$DEST_FILE"

      # Add the file content
      cat "$FILE" >> "$DEST_FILE"

      # Add a blank line for separation
      echo "" >> "$DEST_FILE"
    elif [[ -d "$FILE" ]]; then
      # If it's a directory, call the function recursively
      process_files "$FILE"
    fi
  done
}

# Call the recursive function on the source folder
process_files "$SOURCE_FOLDER"

echo "Process completed. Contents saved in $DEST_FILE"