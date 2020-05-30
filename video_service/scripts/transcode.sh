#!/bin/bash
set -e -x -o pipefail

# http://wiki.webmproject.org/adaptive-streaming/instructions-to-playback-adaptive-webm-using-dash was used as a reference

# This script transcodes the input video to multiple different qualities
# In the future, there should be smarter, more customizable behavior here.
# E.g. if the input video is 360p, there's no need to output a version at 1080p.

# $2 sho

VP9_DASH_PARAMS="-tile-columns 4 -frame-parallel 1"
COMMON_VIDEO_ARGS="-keyint_min 150 -g 150 -threads 8"

# TODO: two-pass encoding/various other parameter tweaks
# See: https://trac.ffmpeg.org/wiki/Encode/VP9


if [ $2 -eq 0 ]
  then
    echo "No value given for argument 2, which controls encoding speed"
    exit 1
fi

ffmpeg -i ${1} -c:v libvpx-vp9 -s 160x90 -b:v 250k ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 ${1}_160x90_250k.webm


ffmpeg -i ${1} -c:v libvpx-vp9 -s 320x180 -b:v 500k ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 ${1}_320x180_500k.webm


ffmpeg -i ${1} -c:v libvpx-vp9 -s 640x360 -b:v 750k ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 ${1}_640x360_750k.webm


ffmpeg -i ${1} -c:v libvpx-vp9 -s 640x360 -b:v 1000k ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 ${1}_640x360_1000k.webm


ffmpeg -i ${1} -c:v libvpx-vp9 -s 1280x720 -b:v 1500k ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 ${1}_1280x720_500k.webm


ffmpeg -i ${1} -c:a libopus -b:a 128k -vn -f webm -dash 1 ${1}_audio_128k.webm

