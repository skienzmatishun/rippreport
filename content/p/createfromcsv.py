import os
import csv
import re

folder_path = "."  # set the folder path to the current directory

# define a function to normalize a string for fuzzy matching
def normalize_string(s):
    return re.sub(r'[^\w\s]', '', s).lower()

# define a function to create a markdown file
def create_markdown_file(folder, filename, frontmatter, content):
    filepath = os.path.join(folder, filename + ".md")
    with open(filepath, "w") as f:
        # write the front matter
        f.write("---\n")
        for key, value in frontmatter.items():
            if isinstance(value, list):
                f.write(f"{key}: {', '.join(value)}\n")
            else:
                f.write(f"{key}: {value}\n")
        f.write("---\n\n")
        # write the content
        f.write(content)

# loop through all csv files in the folder and its subfolders
for root, dirs, files in os.walk(folder_path):
    for file in files:
        if file.endswith(".csv"):
            filepath = os.path.join(root, file)
            # read the csv file
            with open(filepath) as f:
                reader = csv.DictReader(f)
                for row in reader:
                    # check if the comment_post_title is a close match to any folder in the current directory
                    folder_name = normalize_string(row["comment_post_title"])
                    folder_match = None
                    for d in dirs:
                        if normalize_string(d) == folder_name:
                            folder_match = d
                            break
                        elif normalize_string(d).startswith(folder_name[:3]):
                            folder_match = d
                    if folder_match:
                        # create a markdown file inside the matching folder for each row
                        frontmatter = {
                            "title": row["comment_author"],
                            "date": row["comment_date"],
                            "authors": ["Ripp"],
                            "categories": ["Holiday"]
                        }
                        content = row["comment_content"]
                        filename = row["comment_ID"]
                        create_markdown_file(os.path.join(root, folder_match), filename, frontmatter, content)
                    else:
                        # if no close match, print the comment_post_title to the console
                        print(row["comment_post_title"])
