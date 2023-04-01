import os
import re

# define the directory to search for .md files
dir_path = "./p"

# regular expression to match YouTube URLs
youtube_regex = r"https:\/\/(www\.)?youtu(be\.com|\.be)\/(watch\?v=)?([a-zA-Z0-9_-]{11})"

# function to replace YouTube URLs with embedded iframes in a file
def replace_youtube_links(file_path):
    with open(file_path, "r") as f:
        content = f.read()
    # replace YouTube URLs with embedded iframes using regex
    content = re.sub(youtube_regex, r'<iframe width="560" height="315" src="https://www.youtube.com/embed/\4" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>', content)
    with open(file_path, "w") as f:
        f.write(content)
    # print the path of the modified file to the console
    print("Modified file:", file_path)

# loop through all .md files in the directory and its subdirectories
for subdir, dirs, files in os.walk(dir_path):
    for file in files:
        if file.endswith(".md"):
            file_path = os.path.join(subdir, file)
            # call replace_youtube_links function to replace YouTube URLs in the file
            replace_youtube_links(file_path)
