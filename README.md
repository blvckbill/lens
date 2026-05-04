# Lens

A distributed search engine built from scratch in Go. 

Lens is being developed from first principles, strictly avoiding AI code-generation for core systems, to deeply explore information retrieval, memory management, and distributed systems architecture.

## 🚀 Current Status
**Phase 1: Ingestion & Brute-Force Retrieval**
Currently capable of reading local text files and performing raw byte-level string matching using UTF-8 encoding.

## 🗺️ Roadmap
The architecture is being built iteratively from the inside out:
- [x] Disk I/O and File Ingestion
- [ ] Text Analyzer (Tokenization, Lowercasing, Punctuation Stripping)
- [ ] In-Memory Inverted Index
- [ ] Relevance Scoring (TF-IDF / BM25)
- [ ] Disk Persistence (Segment Merging)
- [ ] Distributed Scatter-Gather Routing

## 🛠️ How to Run


```bash
# Clone the repository
git clone [https://github.com/blvckbill/lens.git](https://github.com/blvckbill/lens.git)
cd lens

# Build the executable
go build ./...

# Run the engine
./lens.exe