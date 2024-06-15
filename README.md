# Nakama RPC File Handler

## Overview
This project implements an RPC function for the Nakama game server, which is responsible for reading a file from the disk, calculating its hash, and storing the data in Nakama's database. The function ensures data integrity by comparing provided hashes and handles various edge cases such as missing files or mismatched hashes.

## Assumptions
- The Nakama database is used to store the file content and its metadata.
- Files don't change after being stored on disk, so if there is a corresponding item in the database, we can use it as a cache.
- If the hash provided in the payload is null, return the content of the file in the response.
- The RPC function is made only for server-to-server communication (no user ID in context is allowed).
- Errors are properly logged and returned to the caller in case of issues such as missing files or validation failures.

## Solution Explanation

The implementation includes:

- An RPC function (RpcTransferFile) that processes the payload, validates it, reads the file from disk, calculates its hash, and stores the data in the Nakama database.
- Validation for the type and version fields to ensure they follow expected patterns.
- Unit tests to cover various scenarios, including valid and invalid inputs, file not found errors, and hash mismatches.
- The hashing algorithm used in this implementation is SHA-256, which provides a good balance between security and performance. SHA-256 generates a 256-bit hash value, ensuring data integrity and detection of any tampering.

## Examples
Request
```json
{
   "type": "core",
   "version": "1.0.0",
   "hash": "cbfab3df1f0156ba9eb8e292b754b8cd4f802582ce44b0a0551e918cf3d09092"
}
```
Response (Hash Matches)
```json
{
   "type": "core",
   "version": "1.0.0",
   "hash": "cbfab3df1f0156ba9eb8e292b754b8cd4f802582ce44b0a0551e918cf3d09092",
   "content": {
      "and": "even more data",
      "even": "more data",
      "more": "data",
      "some": "data"
   }
}
```
Response (Hash Does Not Match)
```json
{
   "type": "core",
   "version": "1.0.0",
   "hash": "cbfab3df1f0156ba9eb8e292b754b8cd4f802582ce44b0a0551e918cf3d09092",
   "content": null
}
```

Database Item Example
```json
{
   "type": "core",
   "version": "1.0.0",
   "hash": "292b754b8cd4f80",
   "content": {
      "key": "value"
   }
}
```

## Running the Code
### Prerequisites
- Docker
- Docker Compose
- Make

### run nakama server
```bash
make run
```

### run tests
```bash
make test
```

## Improvements
- Have a conversation with stakeholders to better understand requirements and expectations for this function and, based on this conversation, introduce corresponding improvements.
- Add linters.
- Better API documentation using OpenAPI.
