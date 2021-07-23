#!/bin/bash
set -e -x -o pipefail -u

# h264/AAC:
# https://blogs.gnome.org/rbultje/2015/09/28/vp9-encodingdecoding-performance-vs-hevch-264/
# https://accurate.video/docs/guides/encoding-multi-bitrate-content-optimal-dash-delivery/
# https://github.com/gpac/gpac/wiki/GPAC-Build-Guide-for-Linux

# vp9:
# http://wiki.webmproject.org/adaptive-streaming/instructions-to-playback-adaptive-webm-using-dash was used as a reference
# https://developers.google.com/media/vp9/settings/vod was used as another resource
# https://streaminglearningcenter.com/blogs/encoding-vp9-in-ffmpeg-an-update.html helpful! led me to the google recommended settings
# https://developers.google.com/media/vp9/bitrate-modes is helpful too
# -g is keyframe interval length
# so I believe this would determine the distance between iframes?

# This script transcodes the input video to multiple different qualities
# In the future, there should be smarter, more customizable behavior here.
# E.g. if the input video is 360p, there's no need to output a version at 1080p.

# VP9_DASH_PARAMS="-frame-parallel 1 -row-mt 1 -crf 33"
H264_DASH_PARAMS="-r 24 -x264opts keyint=48:min-keyint=48:no-scenecut -movflags faststart -preset medium -profile:v main -threads 8"
# COMMON_VIDEO_ARGS="-keyint_min 240 -g 240 -threads 8"

# TODO: two-pass encoding/various other parameter tweaks
# See: https://trac.ffmpeg.org/wiki/Encode/VP9

# I'm a little unsure about the CRF values, even though this is what google recommends
# TODO: improve understanding of CRF in the context of targeted avg video bitrate
# CRF is maximum quantization, so I think it'd set an upper bound on compression within a set of frames?
# If CRF was lower, I think it'd be the case that the bitrate would tend to be higher
# But how does the targeted bitrate come into play here? :thinking:
# https://developers.google.com/media/vp9/bitrate-modes#compression helpful
#

# TODO: improve understanding of two-pass encoding

QUAL_FACT=4

# 240p
# x264 --output ${1}_320x240_600k.mp4 --profile main --fps 24 --bitrate 600 --vbv-maxrate 1200 --vbv-bufsize 9600 --min-keyint 48 --keyint 48 --scenecut 0 --no-scenecut --pass 1 --video-filter "resize:width=320,height=240" ${1}
ffmpeg -i ${1} -c:v libx264 ${H264_DASH_PARAMS} -vf scale=320x240 -b:v $(echo $((150 * $QUAL_FACT)))k -minrate 75k -maxrate $(echo $((218 * 2 * $QUAL_FACT)))k -an -f mp4  ${1}_320x240_600k.mp4
#ffmpeg -i ${1} -c:v libvpx-vp9 -vf scale=320x240 -b:v 150k -minrate 75k -maxrate 218k -crf 37 -pass 2 -speed 1 -tile-columns 0 ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 -y ${1}_320x240_150k.webm

# 360p
# x264 --output ${1}_640x360_1000k.mp4 --profile main --fps 24 --bitrate 1100 --vbv-maxrate 2200 --vbv-bufsize 9600 --min-keyint 48 --keyint 48 --scenecut 0 --no-scenecut --pass 1 --video-filter "resize:width=320,height=240" ${1}

ffmpeg -i ${1} -c:v libx264 ${H264_DASH_PARAMS} -vf scale=640x360 -b:v $(echo $((276 * $QUAL_FACT)))k -minrate 138k -maxrate $(echo $((400 * 2 * $QUAL_FACT)))k -an -f mp4 ${1}_640x360_1000k.mp4
#ffmpeg -i ${1} -c:v libvpx-vp9 -vf scale=640x360 -b:v 276k -minrate 138k -maxrate 400k -crf 36 -pass 2 -speed 1 -tile-columns 1 ${COMMON_VIDEO_ARGS} ${VP9_DASH_PARAMS} ${2} -an -f webm -dash 1 -y ${1}_640x360_276k.webm

# TODO: --strict -2
ffmpeg -i ${1} -ac 2 -c:a aac -ar 48000 -b:a 128k -vn -strict -2 ${1}_audio_128k.m4a
