#!/bin/bash
set -e -x -o pipefail

# http://wiki.webmproject.org/adaptive-streaming/instructions-to-playback-adaptive-webm-using-dash was used as a reference

ffmpeg \
 -f webm_dash_manifest -i ${1}_160x90_250k.webm \
 -f webm_dash_manifest -i ${1}_320x180_500k.webm \
 -f webm_dash_manifest -i ${1}_640x360_750k.webm \
 -f webm_dash_manifest -i ${1}_640x360_1000k.webm \
 -f webm_dash_manifest -i ${1}_1280x720_500k.webm \
 -f webm_dash_manifest -i ${1}_audio_128k.webm \
 -c copy -map 0 -map 1 -map 2 -map 3 -map 4 -map 5 \
 -f webm_dash_manifest \
 -adaptation_sets "id=0,streams=0,1,2,3,4 id=1,streams=5" \
 ${1}.mpd