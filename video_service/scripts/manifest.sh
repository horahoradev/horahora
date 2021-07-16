#!/bin/bash
set -e -x -o pipefail

  # http://wiki.webmproject.org/adaptive-streaming/instructions-to-playback-adaptive-webm-using-dash was used as a reference

ffmpeg \
 -f webm_dash_manifest -i ${1}_320x240_600k.mp4 \
 -f webm_dash_manifest -i ${1}_640x360_1000k.mp4 \
 -f webm_dash_manifest -i ${1}_audio_128k.webm \
 -c copy -map 0 -map 1 -map 2 \
 -f webm_dash_manifest \
 -adaptation_sets "id=0,streams=0,1 id=1,streams=2" \
 ${1}.mpd

#  -f webm_dash_manifest -i ${1}_640x360_1000k.webm \
# -f webm_dash_manifest -i ${1}_1280x720_500k.webm \