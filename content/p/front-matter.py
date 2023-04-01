import os
import re

def process_file(filename):
    # Check if the file is a Markdown file
    if not filename.endswith('.md'):
        return

    # Open the file and read the contents
    with open(filename, 'r') as f:
        contents = f.read()

    # Check if the file has front matter
    match = re.match(r'^---\n(.*?)\n---\n(.*)', contents, re.DOTALL)
    if not match:
        return

    front_matter = match.group(1)
    content = match.group(2)

    # Replace the Ripp in the list labeled authors with the original comment_author
    front_matter = re.sub(r'^authors:\n- Ripp', r'comment_author: \g<1>', front_matter, flags=re.MULTILINE)

    # Remove the categories front matter
    front_matter = re.sub(r'^categories:\n- comment\n', '', front_matter, flags=re.MULTILINE)

    # Replace date with comment_date
    front_matter = re.sub(r'^date: (.+)$', r'comment_date: \1', front_matter, flags=re.MULTILINE)

    # Remove the title and comment_author from the front matter
    front_matter = re.sub(r'^title: .+$\n?', '', front_matter, flags=re.MULTILINE)
    front_matter = re.sub(r'^comment_author: .+$\n?', '', front_matter, flags=re.MULTILINE)

    # Write the modified contents to the file
    with open(filename, 'w') as f:
        f.write('---\n')
        f.write(front_matter)
        f.write('---\n')
        f.write(content)

def process_directory(directory):
    for filename in os.listdir(directory):
        path = os.path.join(directory, filename)
        if os.path.isdir(path):
            process_directory(path)
        else:
            process_file(path)

if __name__ == '__main__':
    process_directory('.')
