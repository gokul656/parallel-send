# parallel-send

parallel-send is a Go code developed to share files blazingly fast. It utilizes the power of Go routines to split files into chunks and re-aggregate them efficiently.

Sharing large files quickly over networks can be challenging due to limitations in bandwidth and latency. parallel-send aims to address this challenge by leveraging Go routines to split files into smaller, manageable chunks. These chunks can then be transmitted in parallel, maximizing throughput and reducing transfer times.

## Features

- **Fast File Sharing**: Utilizes parallel processing to split and transmit files quickly.
- **Efficient Reassembly**: Re-aggregates file chunks seamlessly to reconstruct the original file.
- **Simple Interface**: Provides a straightforward API for splitting and reassembling files.
