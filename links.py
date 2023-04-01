import os
import re

# Define regex pattern for detecting URLs
url_pattern = re.compile(r'http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+')

# Define function to convert a URL to a markdown link
def url_to_markdown_link(url):
    return f"[{url}]({url})"

# Define function to scan and update markdown files
def scan_and_update_files(directory):
    # Loop through all files and subdirectories in the given directory
    for root, dirs, files in os.walk(directory):
        for file in files:
            # Check if file is a markdown file
            if file.endswith('.md'):
                # Open file and read content
                with open(os.path.join(root, file), 'r') as f:
                    content = f.read()
                
                # Search for URLs in the content
                urls = url_pattern.findall(content)
                
                # Loop through URLs and check if they are wrapped properly
                for url in urls:
                    if not re.search(r'\[.*?\]\(.*?\)', url):
                        # If URL is not wrapped properly, replace with a markdown link
                        markdown_link = url_to_markdown_link(url)
                        content = content.replace(url, markdown_link)
                
                # Write updated content back to file
                with open(os.path.join(root, file), 'w') as f:
                    f.write(content)

# Example usage
scan_and_update_files('./')
