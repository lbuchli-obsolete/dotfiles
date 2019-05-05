#!/bin/bash

sudo guestmount -a /home/lukas/qemu/windows_hd.qcow2 -m /dev/sda2 /mnt
sudo ranger -r /home/lukas/.config/ranger /mnt/
sudo guestunmount /mnt/
