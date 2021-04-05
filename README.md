# Horahora
Horahora is a microservice-based video hosting website with additional functionality for content archival from Niconico, Bilibili, and Youtube. Users can upload their own content, or schedule categories of content from other websites to be archived (e.g. a given channel on Niconico, a tag on Youtube, or a playlist from Bilibili). Content archived from other websites will be accessible in the same manner as user-uploaded videos, and will be organized under the same metadata (author, tags) associated with the original video.

This project is a WIP, and under active development. The MVP for local use is not yet complete, but nearly done. Contributions and suggestions are always welcome! If you have any questions regarding this project, feel free to contact me via email.

https://discord.gg/vfwfpctJRZ

## Local Use Instructions

1. Install docker, docker-compose, flyway, and make a Backblaze account
2. modify docker-compose.yaml, adding values for the following environment variables: 
    - OriginFQDN: this will be the public URL of your Backblaze bucket WITH NO TRAILING SLASH. E.g. for me it's: https://f002.backblazeb2.com/file/otomads
    - StorageAPIID: the API ID for your Backblaze account
    - StorageAPIKey: The API key for your Backblaze account
    - SocksConn (optional): youtube-dl socks5 flag value string, allowing downloads to be proxied
3. docker-compose up
4. run ./sql/create_and_apply_migrations.sh
5. Visit localhost:8082 (or if it doesn't work initially, try after a few minutes)
    - if it never works, check the container logs, and/or bug me on discord
    - you'll need to login as admin/horahora to view videos that have been encoded. There's an approval workflow which prevents unapproved videos from being viewed by regular users.
    - there's a delay between videos being downloaded/uploaded and being visible, as they need to be transcoded for DASH

## Architecture
![](Architectural_Drawing.png)

Currently, there are three microservices:
1. User Service, which handles registration, logins, and JWT validation
2. Video Service, which handles video uploads (both from Scheduler and from users), transcoding/chunking as required for DASH, uploads to the origin, and storage of metadata.
3. Scheduler, which handles content archival requests from users. For example, if a user specifies that they'd like all videos on Niconico with the tag "YTPMV" to be downloaded, Scheduler will download those videos, register them (and their associated creator) with video service and user service, and check that category of content regularly for new videos.

For more in-depth information on a given microservice, consult its README.

All microservices are horizontally scalable, containerized, and communicate via gRPC.

The MVP will also consist of a frontend service to handle HTML templating, and a Censorship Service (name is a WIP) to manage the workflow for video approvals, and censorship of obscene content from foreign websites.  

## How to Use
Currently, only local use is supported.
To run Horahora locally, follow these steps:
1. First install the following dependencies:
  - Flyway
  - Docker
  - Kubernetes

2. Start minikube with `minikube start --memory=3072`
3. Run `./build.all.sh` (in the Kubernetes directory) to build all Docker images and send to the Docker daemon within minikube
4. `./run-local.sh` to create Kubernetes deployments/services and apply database migrations. Keep running run-local.sh until migrations succeed.
5. `./run-tests.sh` will run local integration tests.
6. Navigate to `localhost:8080` to view the current state of the website.

## Designs
Designs are listed here:
https://github.com/horahoradev/horahora-designs

![](https://github.com/horahoradev/horahora-designs/blob/master/archive.png?raw=true)

![](https://github.com/horahoradev/horahora-designs/blob/master/Video.png?raw=true)

![](https://github.com/horahoradev/horahora-designs/blob/master/Profile.png?raw=true)

## Task Roadmap
Missing features are tracked using Trello.

Our Trello board is:
https://trello.com/b/Rm5TPR4Q/horahora
