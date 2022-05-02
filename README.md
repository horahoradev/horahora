# Horahora
## Self-hosted Video-hosting Website and yt-dlp Video Archival Manager for Niconico, Bilibili, and Youtube
![](https://raw.githubusercontent.com/horahoradev/horahora-designs/master/video_page.png)

![](https://raw.githubusercontent.com/horahoradev/horahora-designs/master/homepage.png)

![](https://raw.githubusercontent.com/horahoradev/horahora-designs/master/Archival_requests_new.png)

![](https://raw.githubusercontent.com/horahoradev/horahora-designs/master/Audit_logs.png)

Horahora is a collaborative archival management tool.

It allows you to:
- download and sync videos from any link supported by yt-dlp
- browse through downloaded videos by channel, tag, views, rating, upload date, etc
- manage archival with a group of friends or untrusted users, with downloads being prioritized by the number of users subscribed to the video's category
- manage site user permissions, ban users, delete videos, and view audit logs for admin/moderator actions

Archival capabilities are provided by yt-dlp (a fork of youtube-dl).

https://discord.gg/vfwfpctJRZ

## Local Use Instructions (START HERE)

1. Install docker and docker-compose
2. sudo make up
3. Wait a minute, then visit localhost:3000
4. Login as admin/admin
    - note that with the current video approval workflow, non-admin users won't be able to view unapproved videos
    - it's recommended to visit /password-reset immediately to change the admin user's default password if using in an untrusted environment
  
If that doesn't work, bug me on Discord.

## Contributing
Contributions are always welcome (and quite needed atm). If you'd like to contribute, and either aren't sure where to start, or lack familiarity with the relevant components of the project, please send me a message on Discord, and I'll help you out as best I can.

## Designs
Designs are listed here:
https://github.com/horahoradev/horahora-designs

## Advanced Use Cases
### Other Storage Backends (s3, backblaze, anything s3-compatible)
By default, Horahora will storage videos locally using Minio.

If you don't want videos to be stored locally, modify secrets.env.template, adding the relevant values for your use case.

    - ORIGIN_FQDN: this will be the public URL of your Backblaze bucket WITH NO TRAILING SLASH. E.g. for me it's: https://f002.backblazeb2.com/file/otomads for backblaze, or https://horahora-dev-otomads.s3-us-west-1.amazonaws.com for s3.
    - STORAGE_BACKEND: 'b2' or 's3' (depending on which you want to use)
    - STORAGE_API_ID: the API ID for your Backblaze account if using backblaze, otherwise blank
    - STORAGE_API_KEY: The API key for your Backblaze account, otherwise blank
    - BUCKET_NAME: the storage bucket name for b2 or s3
  If you want to use S3, you need to include your aws credentials and config in $HOME/.aws. The config and credentials will be mounted into the relevant services at runtime. See https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html for more information.

### Tunneling yt-dlp Traffic
Horahora comes with Gluetun support out of the box. To enable it, you'll need to set the proper values in the "vpn config" section of the secrets.env.template file. This will enable your yt-dlp traffic to be tunneled through your VPN provider via a local Gluetun HTTP proxy.

### Backup Restoration
(this currently isn't functioning, I'll fix it later)

Backup_service writes psql dumps of the three databases (userservice, videoservice, scheduler) to backblaze. To restore, place the three latest dumps in the sql dir, `docker-compose up`, run migrations, then run restore.sh from within the container.
