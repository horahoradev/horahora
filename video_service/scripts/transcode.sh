#!/bin/bash
set -e -x -o pipefail

# http://wiki.webmproject.org/adaptive-streaming/instructions-to-playback-adaptive-webm-using-dash was used as a reference
# https://developers.google.com/media/vp9/settings/vod was used as another resource
# https://streaminglearningcenter.com/blogs/encoding-vp9-in-ffmpeg-an-update.html helpful! led me to the google recommended settings
# https://developers.google.com/media/vp9/bitrate-modes is helpful too
# -g is keyframe interval length
# so I believe this would determine the distance between iframes?

# This script transcodes the input video to multiple different qualities
# In the future, there should be smarter, more customizable behavior here.
# E.g. if the input video is 360p, there's no need to output a version at 1080p.

VP9_DASH_PARAMS="-frame-parallel 1 -row-mt 1 -crf 25"
COMMON_VIDEO_ARGS="-keyint_min 240 -g 240 -threads 8"

# TODO: two-pass encoding/various other parameter tweaks
# See: https://trac.ffmpeg.org/wiki/Encode/VP9


if [ $2 -eq 0 ]
  then
    echo "No value given for argument 2, which controls encoding speed"
    exit 1
fi

# I'm a little unsure about the CRF values, even though this is what google recommends
# TODO: improve understanding of CRF in the context of targeted avg video bitrate

QUAL_FACT=2

# 240p
ffmpeg -i ${1} -c:v libvpx-vp9 -vf scale=320x240 -b:v $(echo $((150 * $QUAL_FACT)))k -minrate 75k -maxrate $(echo $((218 * $QUAL_FACT)))k -tile-columns 0 ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 ${1}_320x240_150k.webm
#ffmpeg -i ${1} -c:v libvpx-vp9 -vf scale=320x240 -b:v 150k -minrate 75k -maxrate 218k -crf 37 -pass 2 -speed 1 -tile-columns 0 ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 -y ${1}_320x240_150k.webm

# 360p
ffmpeg -i ${1} -c:v libvpx-vp9 -vf scale=640x360 -b:v $(echo $((276 * $QUAL_FACT)))k -minrate 138k -maxrate $(echo $((400 * $QUAL_FACT)))k -tile-columns 1 ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 ${1}_640x360_276k.webm
#ffmpeg -i ${1} -c:v libvpx-vp9 -vf scale=640x360 -b:v 276k -minrate 138k -maxrate 400k -crf 36 -pass 2 -speed 1 -tile-columns 1 ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 -y ${1}_640x360_276k.webm

ffmpeg -i ${1} -c:a libopus -b:a 128k -vn -f webm -dash 1 ${1}_audio_128k.webm
