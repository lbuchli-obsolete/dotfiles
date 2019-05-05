#!/bin/bash

date=$(date '+%d-%m-%Y');
path=~/Pictures/screenshot-$date

# test, build, and run.
if [ -e $path\.png ]
then
	# if the name already existed, try adding a number to the end of the
	# file name and increasing it if there still is a file with that name.
	count=1
	while [ -e $path\_$count\.png ]
	do
		count=$(($count + 1))
	done

	grim -g "$(slurp)" $path\_$count\.png
else
	grim -g "$(slurp)" $path\.png
fi
