#!/bin/bash

# start the vm with custom arguments
/usr/bin/qemu-system-x86_64 \
    -monitor stdio \
	-vga qxl \
    -smp 4 \
    -soundhw ac97 \
    -k de-ch \
    -machine type=q35,accel=kvm \
    -m 8192 \
    -hda /home/lukas/qemu/windows_hd.qcow2 \
    -boot once=d,menu=off \
    -rtc base=localtime \
    -name "Windows"
