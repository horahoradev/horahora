(ivan's fork of)

# Horahora

Horahora is a microservice-based video hosting website with additional functionality for content archival from Niconico, Bilibili, and Youtube. Users can upload their own content, or schedule categories of content from other websites to be archived (e.g. a given channel on Niconico, a tag on Youtube, or a playlist from Bilibili). Content archived from other websites will be accessible in the same manner as user-uploaded videos, and will be organized under the same metadata (author, tags) associated with the original video.

This fork is a WIP of a WIP, and not meant for production use. I ([SEAPUNK](https://github.com/SEAPUNK)) am maintaining and making modifications of this project when I feel like it. Don't expect anything to come out of it. If you want to contact me, open up an issue in this repo.

Horahora's repo: https://github.com/horahoradev/horahora

Horahora's discord: https://discord.gg/vfwfpctJRZ

## Run Horahora locally (on Linux/MacOS)

1. Install Docker and docker-compose
2. Clone this repository
3. Run `make up` inside the cloned repo
4. Wait a while, and then navigate to localhost:8082

   Admin username is `admin` and the password is `admin`

There will be a delay between videos being downloaded/uploaded and showing up on the site, as they need to be transcoded first.
