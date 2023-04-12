import os
import re

# Define a regular expression to match URLs with dates
date_url_regex = re.compile(r'(https?://[^\s]*?/)(\d{4})/(\d{2})/(\d{2})/(.*)')

# Define a function to replace matched URLs with dates with their updated version
def replace_date_urls(match):
    url_base = match.group(1)
    url_suffix = match.group(5)
    return url_base + url_suffix

# Define a function to recursively search for markdown files and replace date URLs in them
def search_and_replace(root_dir):
    for root, dirs, files in os.walk(root_dir):
        for file in files:
            if file.endswith('.md'):
                file_path = os.path.join(root, file)
                with open(file_path, 'r') as f:
                    file_contents = f.read()
                updated_contents = date_url_regex.sub(replace_date_urls, file_contents)
                with open(file_path, 'w') as f:
                    f.write(updated_contents)

# Call the search_and_replace function with the current directory as the root directory
search_and_replace('.')
