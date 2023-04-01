import pandas as pd
import os

# Read in the CSV file
df = pd.read_csv('ripp.csv')

# Iterate through each row
for index, row in df.iterrows():

    # Create a folder for this row's post title
    folder_name = row['comment_post_title']
    os.makedirs(folder_name, exist_ok=True)

    # Create the front matter for the markdown file
    front_matter = f"---\ncomment_author: {row['comment_author']}\ncomment_author_email: {row['comment_author_email']}\ncomment_author_url: {row['comment_author_url']}\ncomment_author_IP: {row['comment_author_IP']}\ncomment_date: {row['comment_date']}\ncomment_approved: {row['comment_approved']}\ncomment_parent: {row['comment_parent']}\n---\n"

    # Create the markdown file content
    content = row['comment_content']

    # Write the markdown file to the folder
    file_name = f"{row['comment_ID']}.md"
    file_path = os.path.join(folder_name, file_name)
    with open(file_path, 'w') as f:
        f.write(front_matter + content)
