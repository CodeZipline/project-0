# project-0
## Author: Joshua Nguyen
This is my first GO project. Implementing a concurrent, fast, key-value pair database, with the support of badger database methods.

## Instructions
Use `go get https://github.com/dgraph-io/badger`, and `go get https://github.com/DataDog/zstd`.
    The first package will be the main database methods that this program will use, and second is used for compression of files into different levels of storage within a Log Structured Merge Tree.  

Running this cli program can be done in under the project folder using `go run cmd/project-0cli/project-0cli.go`
    Alternative could be to build the program by using `go build cmd/project-0cli/project-0cli.go`, and the resulting executable is in the project-0 folder and ran with `./project-0cli`

Testing is performed in its own subfolder so unit test are called in isolation from the main cli project. To run this users must be in the directory to run this properly so enter in `cd cmd/testingscripts` then run `go run maintesting.go`. Te unit test calls a specfic set of functions in order and output statements easy to follow to check that Badger functions are working properly. In addition to the unit test, in the same file, exist a set of go routines that will write to the database and read from the datbase. This accomplished with go channels to help regular writes and read order to prevent lookup of keys that do not exist and with RWMutex to have access to resources reserved for a set of routines while blocks others to prevent race conditions. 

## Functions
r - will perform a lookup in the key value database on a given database, which is passed in for the user in this program and a key string which is collected from user input after the read() shortcut r has been entered.

r_keys - will perform create an iterator to index through keys only, which is much faster than regular iteration as it now only checks keys stored in memory on the LSM tree, while values are stored seperately. And values are not prebatched when iterating through items in the database while in this mode. 

w - will perform a key value pairm update to this database, if a value already exist its overrides the value and if not then the key value pair is added to this database. The key is captured first then the value with user input from the command line.

w_TTL - will perform a write action with an additional function that sets a marker for how long this key value pair will exist in the database, the default value is 24 hours, (but can be adjusted with the ttl flag), which requires a interger value and will be used to determine how many hours the item will exist for. 

gc - will perform a unquie garbage collection method. First it will clear up any older values that have been updated or delete any values that are paired with deleted keys. Then it will perform a check of files determine from satistics from compaction of the LSM tree to be used for determining maximun space reclaimation. If no statistics exist random files will be selected until no GC is needed to be perform on that file. When a file is selected it is used to determine if it can be removed and rewriten to save space, so this function by default will run every 10 sec(can be set with gci flag), for a duration of 1 min (gcd flag), with a discard ratio space of 0.5(gcdrs flag, This ratio stands for a file that exist that can be rewritten as half of its current occupying space, this is a write amplification of 2).Calling this would result of at max one file to be removed, otherwise its states that no files have been removed and rewritten due to not reaching the dsr threshold.

d - will perform a delete operation of a key within the database, specifically its places a delete marker in memory to note to GC that this will be eligible for deletion at GC runtime. Note this is similar to updates as an update marker will be places and older values of the key will be also eligible for deletion at GC runtime.

q - exits the database safely.

## Configurations
The configuration files will contain the path to the directory in which the keys log and value log will be stored to, in addition to the changes logs made to the database and the key registry.

## Garbage Collection Notes
1. Due to values being stored seperately in the value logs from keys, values are not effected when keys are removed
2. Concurrent read/writes will accumlate different values for a single key and different version of them
    This will be handled with a call to GC(Garbage Collection function) to answer the two issues above.
        -The GC function will pick value logs to clean up from the statistics off of compaction of the LSM tree, IF no files are found then this will randomly selected files until a found is file that does not need to be cleaned up. 
        -Once a file is selected it is checked for possiblility of rewritting it into a smaller file that satifies the discard ratio space value.
        -Note that only one GC function can be run, other wise errRejected will occur, and This will have a spike in LSM tree when running.