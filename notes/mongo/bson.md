# Binary JSON

Reference: [Data Formats: BSON](https://www.mongodb.com/docs/drivers/go/current/data-formats/bson/)

MongoDB stores documents in a binary representation called BSON that allows for easy and flexible data processing.

The Go driver provides four main types for working with BSON data:

- **A**: An ordered representation of a BSON document (array)
- **M**: An unordered representation of a BSON document (map)
- **D**: An ordered representation of a BSON document (slice)
- **E**: A single element inside a D type

## Struct Tags
