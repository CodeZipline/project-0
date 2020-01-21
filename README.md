# project-0
## Author: Joshua Nguyen
This is my first GO project. Implementing a concurrent, fast, key-value pair database, with the support of badger database methods.

## Instructions
Use go get https://github.com/dgraph-io/badger, and go get https://github.com/DataDog/zstd.
    The first package will be the main database methods that this program will use, and second is used for compression of files into different levels of storage within a Log Structured Merge Tree.  

## Configurations
The configuration files will contain the path to the directory in which the keys log and value log will be stored to, in addition to the changes logs made to the database and the key registry.

## Garbage Collection Notes
1. Due to values being stored seperately in the value logs from keys, values are not effected when keys are removed
2. Concurrent read/writes will accumlate different values for a single key and different version of them
    This will be handled with a call to GC(Garbage Collection function) to answer the two issues above.
        -The GC function will pick value logs to clean up from the statistics off of compaction of the LSM tree, IF no files are found then this will randomly selected files until a found is file that does not need to be cleaned up. 
        -Once a file is selected it is checked for possiblility of rewritting it into a smaller file that satifies the discard ratio space value.
        -Note that only one GC function can be run, other wise errRejected will occur, and This will have a spike in LSM tree when running.