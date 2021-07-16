#!/bin/bash
set -e -x -o pipefail

  # http://wiki.webmproject.org/adaptive-streaming/instructions-to-playback-adaptive-webm-using-dash was used as a reference

MP4Box -dash 4000 -frag 4000 -rap \
-segment-name '$RepresentationID$_' -fps 24 \
${1}_320x240_600k.mp4#video:id=${1}_240p \
${1}_640x360_1000k.mp4#video:id=${1}_360p \
${1}_audio_128k.m4a#audio:id=${1}_128k:role=main \
-out ${1}.mpd
