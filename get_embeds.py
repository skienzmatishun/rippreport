import os
import json
from pathlib import Path
import frontmatter # To parse Markdown front matter
from tqdm import tqdm # For progress bar

# Attempt to import the user-specified lmstudio SDK
try:
    import lmstudio as lms
except ImportError:
    print("Error: The 'lmstudio' Python SDK is not installed or cannot be found.")
    print("Please ensure you have installed the correct SDK that provides 'import lmstudio as lms'.")
    print("If you meant to use the OpenAI-compatible API, please use the previous script version.")
    exit(1)

# --- Configuration ---
# Model identifier for the lmstudio SDK.
# This MUST match how your specific 'lmstudio' SDK identifies the model.
# It could be the Hugging Face name, a GGUF filename, or an alias from LM Studio.
MODEL_IDENTIFIER = "text-embedding-mxbai-embed-large-v1" # Verify this!
CONTENT_DIR = Path("content/p") # Path to your Hugo posts parent directory
OUTPUT_FILE = Path("hugo_embeddings_sdk.json")
# mxbai-embed-large-v1 often performs better for retrieval tasks with a specific prefix
# See: https://huggingface.co/mixedbread-ai/mxbai-embed-large-v1
EMBEDDING_PROMPT_PREFIX = "Represent this sentence for searching relevant passages: "

# --- Initialize LM Studio Embedding Model ---
embedding_model = None
try:
    print(f"Attempting to load embedding model via lmstudio SDK: {MODEL_IDENTIFIER}")
    # This assumes the SDK works as per your example:
    # model = lms.embedding_model("model-identifier")
    embedding_model = lms.embedding_model(MODEL_IDENTIFIER)
    print("LM Studio embedding model loaded successfully.")
except Exception as e:
    print(f"Error loading embedding model via lmstudio SDK: {e}")
    print("Please ensure LM Studio is running, the model is downloaded,")
    print(f"and '{MODEL_IDENTIFIER}' is the correct identifier for your SDK.")
    exit(1)

def get_embedding(text: str) -> list[float] | None:
    """
    Generates an embedding for the given text using the loaded
    lmstudio SDK embedding model.
    """
    global embedding_model
    if not text.strip():
        print("Warning: Attempted to embed empty text.")
        return None
    if embedding_model is None:
        print("Error: Embedding model is not initialized.")
        return None

    try:
        # Prepend the recommended prefix for mxbai-embed-large-v1
        text_to_embed = f"{EMBEDDING_PROMPT_PREFIX}{text}"

        # This assumes the SDK's embed method works as per your example:
        # embedding_vector = model.embed("text")
        # And that it returns a list of floats directly.
        # If it returns a more complex object or a list of embeddings for a list input,
        # this part might need adjustment (e.g., embedding_data[0] if it returns a list).
        embedding_vector = embedding_model.embed(text_to_embed)

        if isinstance(embedding_vector, list) and all(isinstance(item, float) for item in embedding_vector):
            return embedding_vector
        # Handling cases where embed might return a list containing the embedding (e.g. for batch input)
        elif isinstance(embedding_vector, list) and len(embedding_vector) == 1 and \
             isinstance(embedding_vector[0], list) and all(isinstance(item, float) for item in embedding_vector[0]):
            return embedding_vector[0]
        else:
            print(f"Unexpected embedding format received: {type(embedding_vector)}")
            return None

    except Exception as e:
        print(f"Error getting embedding using lmstudio SDK: {e}")
        return None

def extract_content_from_markdown(file_path: Path) -> str:
    """
    Extracts the main content from a Markdown file, skipping front matter.
    """
    try:
        with open(file_path, "r", encoding="utf-8") as f:
            post = frontmatter.load(f)
            return post.content
    except Exception as e:
        print(f"Error reading or parsing markdown file {file_path}: {e}")
        return ""

def main():
    """
    Main function to walk through Hugo content, generate embeddings,
    and save them using the lmstudio SDK.
    """
    all_embeddings_data = {}
    markdown_files_to_process = []

    print(f"Scanning for markdown files in {CONTENT_DIR}...")
    if not CONTENT_DIR.exists() or not CONTENT_DIR.is_dir():
        print(f"Error: Content directory '{CONTENT_DIR}' not found or is not a directory.")
        print("Please check the CONTENT_DIR variable in the script.")
        return

    for post_folder in CONTENT_DIR.iterdir():
        if post_folder.is_dir():
            index_md_path = post_folder / "index.md"
            if index_md_path.exists():
                markdown_files_to_process.append(index_md_path)

    if not markdown_files_to_process:
        print(f"No 'index.md' files found in the subdirectories of {CONTENT_DIR}.")
        print("Please ensure your content structure is 'content/p/post-name/index.md'")
        return

    print(f"Found {len(markdown_files_to_process)} markdown files to process.")

    for md_file_path in tqdm(markdown_files_to_process, desc="Processing files"):
        print(f"\nProcessing: {md_file_path}")
        content = extract_content_from_markdown(md_file_path)

        if content:
            embedding = get_embedding(content)
            if embedding:
                post_slug = md_file_path.parent.name
                relative_path = str(md_file_path.relative_to(CONTENT_DIR.parent))
                all_embeddings_data[post_slug] = {
                    "path": relative_path,
                    "url": f"/{CONTENT_DIR.name}/{post_slug}/",
                    "embedding": embedding,
                }
            else:
                print(f"Could not generate embedding for: {md_file_path}")
        else:
            print(f"No content found in: {md_file_path}")

    if all_embeddings_data:
        try:
            with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
                json.dump(all_embeddings_data, f, indent=4)
            print(f"\nSuccessfully saved embeddings for {len(all_embeddings_data)} posts to {OUTPUT_FILE}")
        except Exception as e:
            print(f"Error saving embeddings to JSON file: {e}")
    else:
        print("\nNo embeddings were generated.")

if __name__ == "__main__":
    main()