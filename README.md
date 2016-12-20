# chordstore

Chordstore is a key-value and object store which uses the chord protocol.  It distributes
and replicates keys around a ring and moves them around as needed when nodes join and leave.

## Usage
```
To be written
```

### Build
```
make deps build
```

### Start demo cluster
```
./start-cluster.sh
```

# To Do

- [ ] Handle un-notified node exits
- [ ] Healing mechanism
- [ ] Persistent store
- [ ] Per key transaction log
