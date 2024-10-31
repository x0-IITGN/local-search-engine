import os
from datetime import datetime
from pymongo import MongoClient

# Connect to MongoDB and create a collection
def init_db(db_name='file_index_db', collection_name='files'):
    client = MongoClient('mongodb://localhost:27017/')
    db = client[db_name]
    collection = db[collection_name]
    return collection

# Insert file metadata into MongoDB
def insert_file_data(collection, name, path, size, modified_time):
    # Check if the file already exists
    existing_file = collection.find_one({'path': path})
    if existing_file is None:
        try:
            collection.insert_one({
                'name': name,
                'path': path,
                'size': size,
                'modified_time': modified_time
            })
            print(f"Indexed: {path}")
        except Exception as e:
            print(f"Failed to index {path}: {e}")
    # else:
        # print(f"Skipped: {path} (already indexed)")

# Crawl a directory and index files
def crawl_directory(root_dir, collection):
    # Define file patterns to exclude (add more as needed)
    exclude_patterns = ['.cph', '.tmp', '.log', '.prob']  # Example extensions to exclude

    for root, _, files in os.walk(root_dir):
        for file in files:
            file_path = os.path.join(root, file)
            
            # Check if the file path contains any exclude patterns
            if any(pattern in file_path for pattern in exclude_patterns) or file.startswith('.'):
                # print(f"Skipped: {file_path} (excluded by pattern)")
                continue

            try:
                size = os.path.getsize(file_path)
                modified_time = datetime.fromtimestamp(os.path.getmtime(file_path)).strftime('%Y-%m-%d %H:%M:%S')
                insert_file_data(collection, file, file_path, size, modified_time)
            except Exception as e:
                print(f"Failed to index {file_path}: {e}")

# Search for files in MongoDB
def search_files(collection, keyword):
    query = {
        '$or': [
            {'name': {'$regex': keyword, '$options': 'i'}},
            {'path': {'$regex': keyword, '$options': 'i'}}
        ]
    }
    results = collection.find(query)
    return list(results)

# Display search results
def display_results(results):
    if results:
        print(f"{'Name':<30} {'Path':<50} {'Size (Bytes)':<15} {'Modified Time'}")
        print('-' * 110)
        for result in results:
            print(f"{result['name']:<30} {result['path']:<50} {result['size']:<15} {result['modified_time']}")
    else:
        print("No matching files found.")

if __name__ == "__main__":
    # Initialize MongoDB collection
    collection = init_db()

    # Crawl directory (provide the path to the directory you want to index)
    root_directory = input("Enter the directory to crawl: ")
    crawl_directory(root_directory, collection)

    # Search for files by name or path
    keyword = input("\nEnter the keyword to search for files: ")
    results = search_files(collection, keyword)
    display_results(results)
