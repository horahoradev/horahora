# Scheduler Microservice

This microservice exposes an API via GRPC to schedule categories of content for download. Scheduler polls the database for categories which it hasn't downloaded from recently, download those videos, and then upload them to Video Service along with the associated metadata.

This program is an example of the multiple publishers multiple subscribers problem. Each program has one or more publishers (see: schedule package) and one or more downloaders (see: downloader package).

## Package Overview
- config: extracts configuration information from the environment, and initializes database connections
- downloader: logic pertaining to downloading categories of content, and uploading videos to Video Service.
- schedule: logic pertaining to selecting categories of content to download, and sending them down to the downloader package.
- grpc: implementation of scheduler's GRPC API
- models: various structs providing abstracted APIs over data store operations

## Workflow
1. Client uses the GRPC API to schedule a category of content for download (e.g. the "YTPMV" tag from Niconico). This causes the download request for the YTPMV tag to be written to scheduler's Postgres database. See the "migrations" directory for information on the schema.
2. One of the database pollers from the schedule package selects the request from the database, and acquires the lock on it to prevent other pollers from publishing it. This transaction is executed in serializable mode, so it should be atomic and concurrency-safe.
3. The download request is sent into the download channel, where it will be picked up by one of the downloader workers.
4. All of the videos from that category of content will be extracted, and downloaded in order. Youtube-dl is used to download videos and extract metadata. Prior to downloading a given video, the downloader will attempt to acquire a lock for it; if the lock can't be acquired, the video will be skipped. This prevents duplicate concurrent downloads of the same video from different categories of content. Redlock is used for this purpose.
5. After a video has been downloaded, it will be uploaded to Video Service.
6. If the upload to Video Service succeeds, a record of the download will be inserted into the previous_downloads table, preventing it from being downloaded again for that category of content. Note: the use of this cache isn't enabled for all categories of content. Video Service will prevent duplicate uploads anyway.
