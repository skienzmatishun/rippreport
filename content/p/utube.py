import os
import re

# Define the regular expression to match YouTube URLs
youtube_regex = r'(https?://)?(www\.)?youtube\.com/watch\?v=([\w-]+)(\S+)?'

# Define the replacement string for the embedded YouTube video
embed_string = r'<iframe width="560" height="315" src="https://www.youtube.com/embed/\3" frameborder="0" allowfullscreen></iframe>'

# Define a function to replace the YouTube URL with the embedded version in a single file
def replace_youtube_urls(file_path):
    with open(file_path, 'r+') as f:
        contents = f.read()
        contents = re.sub(youtube_regex, embed_string, contents)
        f.seek(0)
        f.write(contents)
        f.truncate()

# Recursively search through the current directory for markdown files
for root, dirs, files in os.walk('.'):
    for file in files:
        if file.endswith('.md'):
            file_path = os.path.join(root, file)
            replace_youtube_urls(file_path)
