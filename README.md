# Horahora
## Self-hosted Video-hosting Website and Video Archival Manager for Niconico, Bilibili, and Youtube
![](https://github.com/horahoradev/horahora-designs/blob/master/archive.png?raw=true)

![](https://raw.githubusercontent.com/horahoradev/horahora-designs/master/Archives_1.png)

![](https://github.com/horahoradev/horahora-designs/blob/master/Video.png?raw=true)

Note: the above images are designs, rather than screenshots of the current frontend. The current frontend is essentially a worse version of the above. PRs are welcome 😉 (😭)

Horahora is a microservice-based video hosting website with additional functionality for content archival from Niconico, Bilibili, and Youtube. Users can upload their own content, or schedule categories of content from other websites to be archived (e.g. a given channel on Niconico, a tag on Youtube, or a playlist from Bilibili). Content archived from other websites will be accessible in the same manner as user-uploaded videos, and will be organized under the same metadata (author, tags) associated with the original video.

https://discord.gg/vfwfpctJRZ

## Local Use Instructions

1. Install docker, docker-compose, flyway, and either make a Backblaze account (for b2) or AWS account (for s3).
2. Modify secrets.env.template, adding the relevant values for your use case.
    - ORIGIN_FQDN: this will be the public URL of your Backblaze bucket WITH NO TRAILING SLASH. E.g. for me it's: https://f002.backblazeb2.com/file/otomads for backblaze, or https://horahora-dev-otomads.s3-us-west-1.amazonaws.com for s3.
    - STORAGE_BACKEND: 'b2' or 's3' (depending on which you want to use)
    - STORAGE_API_ID: the API ID for your Backblaze account if using backblaze, otherwise blank
    - STORAGE_API_KEY: The API key for your Backblaze account, otherwise blank
    - BUCKET_NAME: the storage bucket name for b2 or s3

  Once done, rename secrets.env.template to secrets.env.
  If you want to use S3, you need to include your aws credentials and config in $HOME/.aws. The config and credentials will be mounted into the relevant services at runtime. See https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html for more information.

3. Run generate.sh , which will generate docker-compose.yml
4. docker-compose build
5. docker-compose up
6. run ./sql/create_and_apply_migrations.sh
7. Visit localhost:8082 (or if it doesn't work initially, try after a few minutes)
    - if it never works, check the container logs, and/or bug me on discord
    - you'll need to login as admin/horahora to view videos that have been encoded. There's an approval workflow which prevents unapproved videos from being viewed by regular users.
    - there's a delay between videos being downloaded/uploaded and being visible, as they need to be transcoded for DASH
8. If everything comes up correctly, once you're logged in, visit the archival requests tab, and add a new category of content to be archived. If everything works, videos will start to be downloaded, and will be made available after a delay.

## Contributing
Contributions are always welcome (and quite needed atm). If you'd like to contribute, and either aren't sure where to start, or lack familiarity with the relevant components of the project, please send me a message on Discord, and I'll help you out as best I can.

# Horahora

## Designs
Designs are listed here:
https://github.com/horahoradev/horahora-designs

## Task Roadmap
Missing features are tracked using Trello.

Our Trello board is:
https://trello.com/b/Rm5TPR4Q/horahora

Note: this is currently outdated

## Backup Restoration
Backup_service writes psql dumps of the three databases (userservice, videoservice, scheduler) to backblaze. To restore, place the three latest dumps in the sql dir, `docker-compose up`, run migrations, then run restore.sh from within the container.
