#!/bin/bash

# run the install script and set this as startup/login script in your desktop environment (f.e. XFCE: "Session and Startup")
# output will be logged to /tmp/smartlight.log, so you can check what is going on

smservice &> /tmp/smartlight.log &
smtray &