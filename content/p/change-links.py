import os
import re

def process_markdown_file(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        content = file.read()

    # Regular expression to find the specific markdown format
    pattern = re.compile(
        r"(## .+?\n\n### .+?\n\n!\[.*?\]\(https?://.+?\)\n\n.*?\n\nSource: \[.+?\]\(.+?\)\n?)",
        re.DOTALL
    )

    # Replacement format with the <div> wrapper
    def replacement(match):
        return f'<div class="link-preview">\n\n{match.group(0)}\n</div>'

    # Substitute the pattern with the replacement in the content
    new_content = re.sub(pattern, replacement, content)

    # Write the modified content back to the file
    with open(file_path, 'w', encoding='utf-8') as file:
        file.write(new_content)

def process_all_markdown_files(root_dir):
    for subdir, _, files in os.walk(root_dir):
        for file in files:
            if file.endswith('.md'):
                file_path = os.path.join(subdir, file)
                print(f'Processing {file_path}...')
                process_markdown_file(file_path)
                print(f'Processed {file_path}')

if __name__ == "__main__":
    current_directory = os.getcwd()
    process_all_markdown_files(current_directory)
