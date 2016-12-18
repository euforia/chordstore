# chordstore

Chordstore is a key-value store using the chord protocol.  It distributes and replicates
keys around a ring and moves them as nodes join and leave.

## Usage

### Build
```
make deps build
```

### Start demo cluster
```
./start-cluster.sh
```

# To Do

- [ ] Healing mechanism
- [ ] Persistent store
- [ ] Per key transaction log
