# Contributing
Thank you for your interest in contributing to Horahora! Think you can do better than what we've done? You probably can!

## Architectural Overview
![](https://raw.githubusercontent.com/horahoradev/horahora/master/Architectural_Drawing.png)

Horahora's architecture is microservice-based. The main microservices are:
- front_api: which is the RESTful API to the rest of the services
- userservice: which does all authentication and handles user storage/permissions
- videoservice: which does all video storage, uploads to the origin (e.g. s3/backblaze), queries, transcoding, etc
- scheduler: which handles content archival requests and downloads
Communication between userservice, videosercvice, and scheduler is GRPC-based. For more details on individual microservice architecture, see the README for whichever microservice.

Postgresql is used as the database for each service, but Redis is also used for very specific purposes (e.g. distributed locking). Schema migrations can be found within the "migrations" directory within each service. As an example, [here's the migrations directory for Videoservice](https://github.com/horahoradev/horahora/tree/master/video_service/migrations). Migrations are applied using Flyway as a means of providing schema versioning.

Since microservices communicate via GRPC, the API for each service is defined by the domain specific proto3 language. [Here's an example.](https://github.com/horahoradev/horahora/blob/master/video_service/protocol/videoservice.proto) This file in particular defines the API for videoservice, which other services will use to invoke videoservice's functionality. We use this file to generate interface and struct definitions in Golang (our target language), and then implement every method. GRPC implementations for videoservice can be found here: https://github.com/horahoradev/horahora/blob/master/video_service/internal/grpcserver/grpc.go . [Here's a minimal method implementation example](https://github.com/horahoradev/horahora/blob/master/video_service/internal/grpcserver/grpc.go#L442), which rates videos: the implementation simply calls the AddRatingToVideoID method from the videos model (which is really more like an example of the repository pattern) with the supplied GRPC arguments, and returns a response.

## Communication and Coordination
Coordination happens via Discord, here: https://discord.gg/vfwfpctJRZ

Just ping me (Otoman) with a message to the effect of what you're interested in doing, and I'll help you out as best I can.

I've currently listed a few issues for the repo as "help-wanted" and "good first issue"; these represent good places to start, but if you don't see anything you like, or there's nothing that's unassigned, ping me, and I'll think of something that suits your goals.
