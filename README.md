# Tyche
[Tyche](https://en.wikipedia.org/wiki/Tyche) like the goddess of destiny.

## TODO Basic Environment
- Initialize repository command
- Git Manager
  - automatic pull 🚧
  - automatic commit on change 🚧
  - prevent running eg. plugins while committing or pulling 🚧
    - how to handle "cron"jobs? 🚧
- File System Helper
  - File Tree ✅ 
  - metadata for .md files 
    - unique id
    - tags
- Config Infrastructure ✅ 
- Push Notification System via Matrix
- Dockerfile
- Plugin System
  - Plugin folder ✅
    - run on change 🚧
    - run on schedule ✅
    - Typescript Types ... how?
  - Docker Compose
  - NPM
- Plugins
  - Journal ✅
    - statistics for checklist in template 
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
    - statistics
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

##  Getting Started
**I tested Tyche only on linux or mac!**

### Dependencies
Be sure you have following dependencies installed.
- git
- node 
- npm

> ℹ️
> Tyche dose not handle the authentication process for git.  
> You have to set up the ssh key to use Tyche with an repository hoster like GitLab or GitHub.  
> Have a look at the [tutorial from github](https://docs.github.com/en/authentication/connecting-to-github-with-ssh/adding-a-new-ssh-key-to-your-github-account).

### Initialize

**⚠️ This is no yet developed - see ToDo list "Initialize repository command"**
You have to initialize Tyche so it knows which repo to use.  
`node run initialize --repo=git@github.com:i5heu/Tyche-Test.git`
(`git@github.com:i5heu/Tyche-Test.git` is the test repository i use, you can have a look at it if you need a working example)

### Run

To start Tyche run:  
`npm run start`

If you want to run this in production, use something like `supervisord` or use the docker container (⚠️ not yet build)

### Use

TODO
