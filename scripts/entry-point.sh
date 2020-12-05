#!/bin/bash

/bin/bash ./init-mongodbs.sh &
/bin/bash ./init-replica.sh &
# Wait for MongoDB replica set initialisation and run the server.
sleep 100 
./server
