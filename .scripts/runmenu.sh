#!/bin/bash

# create a place to save temporary files
rm -rf /tmp/runmenu
mkdir /tmp/runmenu

# make screenshots
swaygrab /tmp/runmenu/current.png
# TODO find a way to screenshot other workspaces

# blur the screenshot of the current workspace
convert /tmp/runmenu/current.png -blur 0x8 /tmp/runmenu/current_blurred.png 

# draw the time on top
time=`date -u +"%H:%M"`
convert /tmp/runmenu/current_blurred.png -pointsize 64 -fill white \
	-gravity North -annotate +0+16 $time /tmp/runmenu/background.png

feh /tmp/runmenu/background.png

# run rofi
#rofi -show drun
