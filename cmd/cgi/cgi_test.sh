#!/bin/bash
#
printf "Content-type: text/plain\n\n"

printf "example error message\n" > /dev/stderr

if [ "POST" = "$REQUEST_METHOD" -a -n "$CONTENT_LENGTH" ]; then
	read -n "$CONTENT_LENGTH" POST_DATA
	fi

	printf "AUTH_TYPE         [%s]\n" $AUTH_TYPE
	printf "CONTENT_LENGTH    [%s]\n" $CONTENT_LENGTH
	printf "CONTENT_TYPE      [%s]\n" $CONTENT_TYPE
	printf "GATEWAY_INTERFACE [%s]\n" $GATEWAY_INTERFACE
	printf "PATH_INFO         [%s]\n" $PATH_INFO
	printf "PATH_TRANSLATED   [%s]\n" $PATH_TRANSLATED
	printf "POST_DATA         [%s]\n" $POST_DATA
	printf "QUERY_STRING      [%s]\n" $QUERY_STRING
	printf "REMOTE_ADDR       [%s]\n" $REMOTE_ADDR
	printf "REMOTE_HOST       [%s]\n" $REMOTE_HOST
	printf "REMOTE_IDENT      [%s]\n" $REMOTE_IDENT
	printf "REMOTE_USER       [%s]\n" $REMOTE_USER
	printf "REQUEST_METHOD    [%s]\n" $REQUEST_METHOD
	printf "SCRIPT_EXEC       [%s]\n" $SCRIPT_EXEC
	printf "SCRIPT_NAME       [%s]\n" $SCRIPT_NAME
	printf "SERVER_NAME       [%s]\n" $SERVER_NAME
	printf "SERVER_PORT       [%s]\n" $SERVER_PORT
	printf "SERVER_PROTOCOL   [%s]\n" $SERVER_PROTOCOL
	printf "SERVER_SOFTWARE   [%s]\n" $SERVER_SOFTWARE

	exit 0
