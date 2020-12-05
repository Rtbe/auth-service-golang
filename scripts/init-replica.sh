#!/bin/bash

sleep 30 

mongo <<EOF
   var cfg = {
        "_id": "rs0",
        "version": 1,
        "members": [
            {
                "_id": 0,
                "host": "127.0.0.1:27017",
                "priority": 2
            },
            {
                "_id": 1,
                "host": "127.0.0.1:27018",
                "priority": 0
            },
            {
                "_id": 2,
                "host": "127.0.0.1:27019",
                "priority": 0
            }
        ]
    };
    rs.initiate(cfg, { force: true });
EOF

sleep 60 

mongo <<EOF
   use admin;
   admin = db.getSiblingDB("admin");
   admin.createUser(
     {
	user: "admin",
        pwd: "password",
        roles: [ { role: "root", db: "admin" } ]
     });
     db.getSiblingDB("admin").auth("admin", "password");
     rs.status();
     use testTask;
     db.createCollection("tokens");
EOF