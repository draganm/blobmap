# Blobmap

**Blobmap** is a data structure designed for efficient storage and retrieval of Binary Large Objects (Blobs), using a continuous keyspace of 64-bit unsigned integers (`uint64`). The keyspace starts at a specified value `n` and covers `m` consecutive keys, ranging from `n` to `n + m - 1`.

The data is stored in a read-only memory-mapped file for constant-time (`O(1)`) access, allowing scalability for large datasets.

## File Format

### Structure

#### 1. Header
The file begins with a header containing essential metadata:
- Number of blobs
- Key offset (starting key `n`)

#### 2. Offset Records
Each record stores the **end offset** of the corresponding blob (64-bit big-endian integer). The start of a blob is implicitly the end of the previous blob.

#### 3. Blob Data
The blob data is stored sequentially, accessed by the offset records.

## Example Layout

```
+----------------+----------------------+-------------------+
|     Header     |  Offset Records      |    Blob Data      |
+----------------+----------------------+-------------------+
| Num of Blobs   | End of Blob 1        | Blob 1 Data Bytes |
| Key Offset (n) | End of Blob 2        | Blob 2 Data Bytes |
|                | End of Blob 3        | Blob 3 Data Bytes |
+----------------+----------------------+-------------------+
```

In this format, the data for each blob is accessed by calculating the byte range between consecutive offsets. The first blob starts immediately after the offset records, and each subsequent blob starts where the previous one ends.