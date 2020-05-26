# Scheduler Microservice

This microservice will expose an interface via GRPC to schedule categories of videos to be downloaded.

This program is an example of the multiple publishers multiple subscribers problem.

On graceful shutdown, all published items are guaranteed to be dealt with.