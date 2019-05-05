#!/bin/bash

# start the ssh agent
if ps -p $SSH_AGENT_PID > /dev/null
then
   echo "SSH-Agent is already running"
else
	eval `ssh-agent -s`
fi

# if gitssh was added
if ssh-add -L | grep gitssh
then
	echo "Using gitssh key."
	exit 0
else
	if ssh-add ~/.ssh/gitssh
	then
		echo "Authentication succeeded."
		exit 0
	else
		echo "Error: Could not add key ~/.ssh/gitssh"
		exit 1
	fi
fi

exit 1
