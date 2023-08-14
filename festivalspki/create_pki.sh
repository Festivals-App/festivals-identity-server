#!/bin/bash
#
# create_pki.sh 1.0.0
#
# Creates a CA and than uses it to create a certain amount of server certificates
#
# (c)2020-2023 Simon Gaus
#

# Check if all passwords are supplied
#
if [ $# -ne 3 ]; then
    echo "$0: usage: sudo ./create_pki.sh <ca_root_pw> <mysql_backup_pw> <read_write_pw>"
    exit 1
fi