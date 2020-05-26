#!/bin/bash
set -e -x -o pipefail

# http://wiki.webmproject.org/adaptive-streaming/instructions-to-playback-adaptive-webm-using-dash was used as a reference

# This script transcodes the input video to multiple different qualities
# In the future, there should be smarter, more customizable behavior here.
# E.g. if the input video is 360p, there's no need to output a version at 1080p.

VP9_DASH_PARAMS="-tile-columns 4 -frame-parallel 1"


ffmpeg -i ${1}.mp4 -c:v libvpx-vp9 -s 160x90 -b:v 250k -keyint_min 150 -g 150 ${VP9_DASH_PARAMS} -an -f webm -dash 1 ${1}_160x90_250k.webm


ffmpeg -i ${1}.mp4 -c:v libvpx-vp9 -s 320x180 -b:v 500k -keyint_min 150 -g 150 ${VP9_DASH_PARAMS} -an -f webm -dash 1 ${1}_320x180_500k.webm


ffmpeg -i ${1}.mp4 -c:v libvpx-vp9 -s 640x360 -b:v 750k -keyint_min 150 -g 150 ${VP9_DASH_PARAMS} -an -f webm -dash 1 ${1}_640x360_750k.webm


ffmpeg -i ${1}.mp4 -c:v libvpx-vp9 -s 640x360 -b:v 1000k -keyint_min 150 -g 150 ${VP9_DASH_PARAMS} -an -f webm -dash 1 ${1}_640x360_1000k.webm


ffmpeg -i ${1}.mp4 -c:v libvpx-vp9 -s 1280x720 -b:v 1500k -keyint_min 150 -g 150 ${VP9_DASH_PARAMS} -an -f webm -dash 1 ${1}_1280x720_500k.webm


ffmpeg -i ${1}.mp4 -c:a libopus -b:a 128k -vn -f webm -dash 1 ${1}_audio_128k.webm

