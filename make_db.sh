#!/usr/bin/env sh
#
# Make an SQLite database for the server.

if [ -f "messages.sqlite" ]
then
	printf "Database already exists.\n"
	exit
fi

touch messages.sqlite
sqlite3 ./messages.sqlite 'CREATE TABLE messages ( time TIMESTAMP, author TEXT, data TEXT);'
