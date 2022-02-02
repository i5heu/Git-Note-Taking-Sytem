# Tyche
[Tyche](https://en.wikipedia.org/wiki/Tyche) like the goddess of destiny.

## TODO Basic Environment
- Git Manager
  - automatic pull
  - automatic commit on change
  - prevent running eg. plugins while committing or pulling
    - how to handle "cron"jobs?
- File System Helper
  - File Tree ✅ 
  - metadata for .md files 
    - unique id
    - tags
- Config Infrastructure ✅ 
- Plugin System
  - Plugin folder ✅
    - run on change ✅
    - run on schedule ✅
    - Typescript Types ... how?
  - Docker Compose
  - NPM
- Plugins
  - Journal ✅
  - SaveLinkAsPdfArchive
    - docker headless chromium for PDFs
    - compress size
  - PDF and Image OCR
    - place file with text into the OCR folder
    - link back to source file
    - delete file if source file was deleted
  - Indexer
    - index tags and hashes
  - ToDo System
    - unique id
    - priority
    - estimates of workload
    - dependencies
    - daily automated ToDo lists
    - recurring ToDos
    - generate file for ToDo, in the ToDo folder, with tag and link 
  - Markdown Table Calculator
  - Automatic Encrypted Backups 
    - S3 (Backblaze)
    - Deadman Switch
    - Error Push Notification
  - Compress Git Vectors
    - delete vectors for files in configurable paths and got deleted
    - delete vectors for binary files that got deleted
- Debugger

## TODO Plugins
- Journal 
- Link to PDF Archive
- Indexing
- Markdown Table Calculations

✅ Done  
🚧 Under Development  
🐛 Buggy
