# Video Service

Video Service is a microservice written in Golang which handles management of video uploads, video metadata storage, video storage, transcoding and chunking as required for DASH.

## Package Overview
- config:
- dashutils: utilities relating to transcoding and chunking as required for DASH.
- grpcserver: implements Video Service's GRPC API
- model: abstractions over database operations for videos

## Overview of Workflow
### Video Uploads
1. Scheduler uploads the video to Video Service via GRPC
2. Video Service reads the incoming bytes for the video, and writes them to a temporary file
3. The temporary video file is transcoded and chunked for DASH, and the DASH manifest is generated (see scripts.sh for the scripts used for transcoding and generating the manifest).
4. The transcoded video files and DASH manifest are uploaded to AWS S3
5. If the video is foreign (it was downloaded from another website via Scheduler), Video Service will use User Service's GRPC API to check whether a domestic user for that author already exists. If one doesn't exist, it will be created.
6. The video is written to the videos table along with the author's domestic user ID.
At this point, the video will be returned to the frontend via the getVideoList API.

### TODO
- creation of domestic users for foreign authors will fail if a user already exists with their username