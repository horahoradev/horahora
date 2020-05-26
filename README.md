# horahora
Horahora is a microservice-based video hosting website with additional functionality for group content archival from Niconico, Bilibili, and Youtube.

This project is a WIP, and under active development. 

Join our discord here: https://discord.gg/psuj8QQ

## Architecture
TODO

## How to Use
Currently, only local use is supported.
To run horahora locally, follow these steps:
1. First install the following depdendencies:
  - Flyway
  - Docker
  - Kubernetes

2. Start minikube, and use `./run-local.sh` in the Kubernetes directory. If the database migrations fail to apply, keep running run-local.sh until they succeed.
3. `./run-tests.sh` will run local integration tests. Currently, this will send an archival request to scheduler for all YTPMVs on Niconico.

## Missing Essential Features
The following is a non-exhaustive list of features which should be added for the MVP:
1. A frontend
2. redis locking (or some other form of distributed locking) for video downloads to prevent concurrent downloads of the same videos from two categories of content
3. Extended archival request support. Currently, only tags archive requests from Niconico are supported.
4. Expanded unit tests
5. Expanded integration tests
6. A less awkward local development workflow
7. All necessary AWS infrastructure:
  - autoscaling EKS cluster
  - log aggregation
  
## Missing Non-essential Features
1. L7 load balancing between services with Envoy

