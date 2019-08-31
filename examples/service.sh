#!/bin/bash

# runs the service and in case it crashes, restarts it.

until smservice &>> /tmp/smartlight.log; do
    echo "smservice crashed with exit code $?. Respawning..." &>> /tmp/smartlight.log
    sleep 1
done