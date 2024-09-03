# Blobmap

A Blobmap is a data structure designed to efficiently store and access a collection of Binary Large Objects (Blobs). Each Blobmap utilizes a continuous keyspace, where the keys are 64-bit unsigned integers (`uint64`). This keyspace starts at a specified value `n` and covers `m` consecutive keys, ranging from `n` to `n+m-1`.

The Blobmap is implemented as a read-only memory-mapped file, providing constant-time (`O(1)`) access to the values stored within it. This design ensures both high efficiency and scalability for large datasets.
