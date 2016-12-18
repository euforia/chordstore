#!/bin/bash
#
# This script spins up a three node chordstore cluster.  It will kill all 3
# nodes on an interrupt or term.  Trying to interrupt before all 3 nodes are spun
# up will leave behind straggling processes.
#

[ -x "./chordstore" ] || { echo "chordstore not found!"; exit 1; }

./chordstore &
p1=$!

sleep 2;
./chordstore -b 127.0.0.1:3244 -j 127.0.0.1:3243 -http 127.0.0.1:9091 &
p2=$!


sleep 2;
./chordstore -b 127.0.0.1:3245 -j 127.0.0.1:3244 -http 127.0.0.1:9092 &
p3=$!

trap "{ kill $p1; kill $p2; kill $p3; }" SIGINT SIGTERM

wait
