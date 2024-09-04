# Blobmap

**Blobmap** is a specialized data structure for efficient storage and retrieval of Binary Large Objects (Blobs) within a continuous keyspace of 64-bit unsigned integers (`uint64`). The keyspace begins at a specified value `n` and spans `m` consecutive keys, covering the range from `n` to `n + m - 1`.

The data is stored in a read-only memory-mapped file to enable constant-time (`O(1)`) access, making it scalable for handling large datasets.

## File Format

### Structure Overview

#### 1. Header
The file begins with a header containing metadata critical for blob management:
- **Number of blobs**: The total count of blobs in the file.
- **Key offset**: The starting key `n` for the keyspace.

#### 2. Offset Table
An array of **offset records** follows the header. Each record stores the **end offset** of a blob, encoded as a 64-bit big-endian integer. The start of each blob is implicitly defined by the end offset of the preceding blob.

#### 3. Blob Data
The blobs themselves are stored sequentially in the file. The data for each blob can be accessed by determining its byte range from the offset records.

#### 4. Integrity Check
The file concludes with an [xxHash](https://xxhash.com/) checksum, covering all preceding data. This can be used to verify the integrity of the **blobmap** during reads.

## Example Layout

```
+----------------+----------------------+-------------------+-------------+
|     Header     |    Offset Table      |    Blob Data      |   xxHash    |
+----------------+----------------------+-------------------+-------------+
| Num of Blobs   | End Offset of Blob 1 | Blob 1 Data Bytes | Hash Value  |
| Key Offset (n) | End Offset of Blob 2 | Blob 2 Data Bytes |             |
|                | End Offset of Blob 3 | Blob 3 Data Bytes |             |
+----------------+----------------------+-------------------+-------------+
```

- **Header**: Stores the number of blobs and the starting key.
- **Offset Table**: Defines the end offsets of each blob.
- **Blob Data**: Contains the actual binary data of each blob, laid out sequentially.
- **xxHash**: Provides a checksum to ensure data integrity.

### Blob Access
To access a specific blob, compute its byte range using the corresponding offsets in the table:
- The start of blob `i` is the end offset of blob `i-1` (or immediately after the offset table for the first blob).
- The end of blob `i` is the offset at position `i` in the table.

This layout enables fast, direct access to any blob, minimizing overhead and maximizing scalability for large datasets.