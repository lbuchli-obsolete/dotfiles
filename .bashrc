#
# ~/.bashrc
#

# If not running interactively, don't do anything
[[ $- != *i* ]] && return

# Use .Xresources
xrdb ~/.Xresources

alias ls='ls --color=auto'
#PS1='[\u@\h \W]\$ '

# ll as a short ls -al
alias ll='ls -al'

# quit with q
alias q='exit'

# add go binaries to path
PATH=$PATH:/home/lukas/go/bin

# add cargo (rust) binaries to path
PATH=$PATH:/home/lukas/cargo/bin

# custom ps1
source ~/.scripts/prompt.sh

# start the ssh-agent
eval "$(ssh-agent -s)" > /dev/null

# fancy dotfiles management
alias config='/usr/bin/git --git-dir=/home/lukas/.cfg/ --work-tree=/home/lukas'
