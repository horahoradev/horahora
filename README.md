# Horahora
Horahora is a microservice-based video hosting website with additional functionality for group content archival from Niconico, Bilibili, and Youtube. Users can upload their own content, or schedule categories of content from other websites to be archived (e.g. a given channel on Niconico, a tag on Youtube, or a playlist from Bilibili). Content archived from other websites will be accessible in the same manner as user-uploaded videos, and will be organized under the same metadata (author, tags) associated with the original video.

This project is a WIP, and under active development. 

Join our discord here: https://discord.gg/6TAEmAA

## Architecture
![](Architectural_Drawing.png)

Currently, there are three microservices:
1. User Service, which handles registration, logins, and JWT validation
2. Video Service, which handles video uploads (both from Scheduler and from users), transcoding/chunking as required for DASH, uploads to the origin, and storage of metadata.
3. Scheduler, which handles content archival requests from users. For example, if a user specifies that they'd like all videos on Niconico with the tag "YTPMV" to be downloaded, Scheduler will download those videos, register them (and their associated creator) with video service and user service, and check that category of content regularly for new videos.

All microservices are horizontally scalable, containerized, and communicate via gRPC.

The MVP will also consist of a frontend service to handle HTML templating, and a Censorship Service (name is a WIP) to manage the workflow for video approvals, and censorship of obscene content from foreign websites.  

## How to Use
Currently, only local use is supported.
To run horahora locally, follow these steps:
1. First install the following depdendencies:
  - Flyway
  - Docker
  - Kubernetes

2. Start minikube, and use `./run-local.sh` in the Kubernetes directory. If the database migrations fail to apply, keep running run-local.sh until they succeed.
3. Run `./build.all.sh` to build all Docker images and send to the Docker daemon within minikube
4. `./run-tests.sh` will run local integration tests. Currently, this will send an archival request to scheduler for all YTPMVs on Niconico.

## Designs
Designs are listed here:
https://github.com/horahoradev/horahora-designs

![](https://github.com/horahoradev/horahora-designs/blob/master/Login.png?raw=true)

![](https://github.com/horahoradev/horahora-designs/blob/master/Video.png?raw=true)

![](https://github.com/horahoradev/horahora-designs/blob/master/Profile.png?raw=true)

## Task Roadmap
Missing features are tracked using Trello.

Our Trello board is:
https://trello.com/b/Rm5TPR4Q/horahora



