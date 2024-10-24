import os
from datetime import datetime

# Function to crawl a directory and index files
def crawl_directory(root_dir):
    indexed_files = []  # List to hold indexed file metadata
    exclude_patterns = ['.cph', '.tmp', '.log']  # Define file patterns to exclude

    for root, _, files in os.walk(root_dir):
        for file in files:
            file_path = os.path.join(root, file)
            
            # Check if the file path contains any exclude patterns
            if any(pattern in file_path for pattern in exclude_patterns) or file.startswith('.'):
                print(f"Skipped: {file_path} (excluded by pattern)")
                continue

            try:
                size = os.path.getsize(file_path)
                modified_time = datetime.fromtimestamp(os.path.getmtime(file_path)).strftime('%Y-%m-%d %H:%M:%S')
                
                # Store file metadata in the list
                indexed_files.append({
                    'name': file,
                    'path': file_path,
                    'size': size,
                    'modified_time': modified_time
                })
                print(f"Indexed: {file_path}")
            except Exception as e:
                print(f"Failed to index {file_path}: {e}")
    
    return indexed_files

# Function to search for files in the indexed list
def search_files(indexed_files, keyword):
    results = [file for file in indexed_files if keyword.lower() in file['name'].lower() or keyword.lower() in file['path'].lower()]
    return results

# Function to display search results
def display_results(results):
    if results:
        print(f"{'Name':<30} {'Path':<50} {'Size (Bytes)':<15} {'Modified Time'}")
        print('-' * 110)
        for result in results:
            print(f"{result['name']:<30} {result['path']:<50} {result['size']:<15} {result['modified_time']}")
    else:
        print("No matching files found.")

if __name__ == "__main__":
    # Crawl directory (provide the path to the directory you want to index)
    root_directory = input("Enter the directory to crawl: ")
    indexed_files = crawl_directory(root_directory)

    # Search for files by name or path
    keyword = input("\nEnter the keyword to search for files: ")
    results = search_files(indexed_files, keyword)
    display_results(results)
