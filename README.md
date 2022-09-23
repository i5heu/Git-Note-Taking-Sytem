# Tyche
[Tyche](https://en.wikipedia.org/wiki/Tyche) like the goddess of destiny.  
  
Tyche is the central Server that manages a git repository to make it useable as a Personal Knowledge Base.  
It is written in Go and facilitates a microservice architecture in which the clients are connected via websockets.  
For debugging purposes, the server also provides a web interface.  

## Scope of this project
This Server dose following things:  

### Coordination of one git repository
- This single git repo is used as a "database for files with version control".  
- The Server creates a new branch for the current epoch in which every client will get a commit if they write something.
- If all clients give up their write access and tag the repo as "finished", the Server squashes the commits into main and pushes them to remote.
- Preventing Race Conditions
  - The Server grants write access to the repo only to one client at a time.
  - All other clients are granted read access to the repo.
  - A clients write access is revoked if no write was made for 10 seconds. The made changes of this client will be reverted.

### WebSocket API
- Authentication and Registration 🚧
- push and receive notifications
  - Subscribe to file changes
- Requesting write access ✅
  - timeout write authority after 10 seconds of inactivity
  - revoke write authority
- health check
  - check clients
  - check server
- Files
  - Read ✅
  - Read directory ✅
  - Read file history
  - Write ✅
  - Delete 
  - Hash file for raze conditions
  - Hash tree for comparison of file trees

## to consider:
- Full text search API
- Query API
  - Query Files like they are in a sql database
  - allow a database abstraction layer that actually uses plain text files for storage

### Web Interface
- Authentication
- Debugging for the WebSocket API


## Websocket API Documentation
