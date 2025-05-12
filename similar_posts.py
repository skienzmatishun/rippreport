import json
import shutil
from pathlib import Path
import numpy as np
import frontmatter # for Markdown frontmatter handling
from tqdm import tqdm # For progress bar

# --- Configuration ---
# Path to your generated embeddings JSON file
EMBEDDINGS_FILE_PATH = Path("hugo_embeddings.json") # Or hugo_embeddings_sdk.json
# Root directory of your original Hugo content (where 'p' folder is located)
# If your embeddings 'path' is "p/post-name/index.md", ORIGINAL_CONTENT_ROOT should point to "content"
ORIGINAL_CONTENT_ROOT = Path("content")
# Directory where backup posts with added 'similar' frontmatter will be stored
BACKUP_ROOT_DIR = Path("../backup_posts")
# Number of similar posts to include
N_SIMILAR = 5

def load_embeddings(file_path: Path) -> dict:
    """Loads embeddings from a JSON file."""
    if not file_path.exists():
        print(f"Error: Embeddings file not found at {file_path}")
        return {}
    try:
        with open(file_path, "r", encoding="utf-8") as f:
            return json.load(f)
    except json.JSONDecodeError:
        print(f"Error: Could not decode JSON from {file_path}")
        return {}
    except Exception as e:
        print(f"An error occurred while reading {file_path}: {e}")
        return {}

def cosine_similarity(vec1: list[float], vec2: list[float]) -> float:
    """Calculates cosine similarity between two vectors."""
    if not vec1 or not vec2:
        return 0.0
    vec1_np = np.array(vec1)
    vec2_np = np.array(vec2)
    dot_product = np.dot(vec1_np, vec2_np)
    norm_vec1 = np.linalg.norm(vec1_np)
    norm_vec2 = np.linalg.norm(vec2_np)

    if norm_vec1 == 0 or norm_vec2 == 0:
        return 0.0  # Avoid division by zero
    return dot_product / (norm_vec1 * norm_vec2)

def find_top_n_similar(target_slug: str, all_embeddings_data: dict, n: int = 5) -> list[dict]:
    """Finds the top N similar posts for a target post."""
    if target_slug not in all_embeddings_data:
        print(f"Warning: Target slug '{target_slug}' not found in embeddings data.")
        return []

    target_embedding = all_embeddings_data[target_slug].get("embedding")
    if not target_embedding:
        print(f"Warning: Embedding not found for target slug '{target_slug}'.")
        return []

    similarities = []
    for slug, data in all_embeddings_data.items():
        if slug == target_slug:
            continue # Don't compare a post to itself

        embedding = data.get("embedding")
        if not embedding:
            # print(f"Warning: Skipping post '{slug}' due to missing embedding.")
            continue

        similarity = cosine_similarity(target_embedding, embedding)
        similarities.append({
            "slug": slug,
            "similarity": similarity,
            "url": data.get("url", f"/{data.get('path', slug).replace('index.md', '').strip('/')}/"), # Construct URL if missing
            "title": slug # Using slug as title, can be enhanced if title is in embeddings
        })

    # Sort by similarity in descending order
    similarities.sort(key=lambda x: x["similarity"], reverse=True)

    top_n = []
    for item in similarities[:n]:
        top_n.append({
            "title": item["title"], # This is currently the slug
            "url": item["url"]
        })
    return top_n

def add_similar_to_frontmatter(md_filepath: Path, similar_posts_data: list[dict]):
    """Adds or updates the 'similar' key in the frontmatter of a Markdown file."""
    try:
        if not md_filepath.exists():
            print(f"Error: Markdown file not found for updating frontmatter: {md_filepath}")
            return

        # Read the post using "r" for text mode
        with open(md_filepath, "r", encoding="utf-8") as f:
            post = frontmatter.load(f)

        # Update metadata
        post.metadata["similar"] = similar_posts_data

        # Serialize the entire post (frontmatter + content) to a string
        # The dumps() function should return a string.
        updated_content_str = frontmatter.dumps(post)

        # Write the string content back to the file, ensuring UTF-8 encoding
        with open(md_filepath, "w", encoding="utf-8") as f:
            f.write(updated_content_str) # Explicitly writing the string

        # print(f"Added similar posts to frontmatter of {md_filepath}") # Optional: uncomment for verbose logging

    except Exception as e:
        print(f"Error updating frontmatter for {md_filepath}: {e}")
        # For more detailed debugging, you can uncomment the following lines:
        # import traceback
        # traceback.print_exc()


def main():
    """Main script execution."""
    print("Starting script to add similar posts to backups...")

    all_embeddings = load_embeddings(EMBEDDINGS_FILE_PATH)
    if not all_embeddings:
        print("Exiting due to issues loading embeddings.")
        return

    if not ORIGINAL_CONTENT_ROOT.exists() or not ORIGINAL_CONTENT_ROOT.is_dir():
        print(f"Error: Original content root directory '{ORIGINAL_CONTENT_ROOT}' not found.")
        print("Please ensure ORIGINAL_CONTENT_ROOT is set correctly.")
        return

    # Create the main backup directory if it doesn't exist
    BACKUP_ROOT_DIR.mkdir(parents=True, exist_ok=True)
    print(f"Backup directory set to: {BACKUP_ROOT_DIR.resolve()}")

    processed_count = 0
    for post_slug, data in tqdm(all_embeddings.items(), desc="Processing posts"):
        original_relative_path_str = data.get("path")
        if not original_relative_path_str:
            print(f"Warning: 'path' missing for slug '{post_slug}'. Skipping.")
            continue

        original_file_path = ORIGINAL_CONTENT_ROOT / original_relative_path_str
        if not original_file_path.exists():
            print(f"Warning: Original file '{original_file_path}' not found for slug '{post_slug}'. Skipping.")
            continue

        # Determine backup path
        # original_relative_path_str is like "p/post-name/index.md"
        backup_file_path = BACKUP_ROOT_DIR / original_relative_path_str

        # Create parent directories for the backup file
        backup_file_path.parent.mkdir(parents=True, exist_ok=True)

        # 1. Copy the original file to the backup location
        try:
            shutil.copy2(original_file_path, backup_file_path)
        except Exception as e:
            print(f"Error copying '{original_file_path}' to '{backup_file_path}': {e}")
            continue # Skip this post if copy fails

        # 2. Find top N similar posts
        similar_posts = find_top_n_similar(post_slug, all_embeddings, N_SIMILAR)

        if not similar_posts:
            print(f"No similar posts found for '{post_slug}'. Frontmatter will not be updated for similarity.")
            # still the file is copied, so it's a backup.
        else:
            # 3. Add similar posts to the frontmatter of the *backed-up* file
            add_similar_to_frontmatter(backup_file_path, similar_posts)
        processed_count += 1

    print(f"\nProcessing complete. Processed {processed_count} posts.")
    print(f"Backed-up posts with 'similar' frontmatter (if applicable) are in: {BACKUP_ROOT_DIR.resolve()}")

if __name__ == "__main__":
    main()