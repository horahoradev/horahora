# Horahora
## Locally archive, browse, and share videos from nearly any site

Horahora is a collaborative archival management tool.

It allows you to:
- download and continuously sync videos from any link supported by yt-dlp
- browse through downloaded videos by channel, tag, views, rating, upload date, etc
- manage archival with a group of friends or untrusted users, with downloads being prioritized by the number of users subscribed to the video's category
- manage site user permissions, ban users, delete videos, and view audit logs for admin/moderator actions

Join our Discord: https://discord.gg/7te9UrbydS

![](https://github.com/horahoradev/horahora-designs/blob/master/Screenshot%20from%202022-10-09%2011-56-34.png?raw=true)

![](https://github.com/horahoradev/horahora-designs/blob/master/Screenshot%20from%202022-10-09%2011-54-48.png?raw=true)

![](https://github.com/horahoradev/horahora-designs/blob/master/Screenshot%20from%202022-10-09%2011-57-35.png?raw=true)

![](https://github.com/horahoradev/horahora-designs/blob/master/Screenshot%20from%202022-10-09%2011-57-52.png?raw=true)

Archival capabilities are provided by yt-dlp (a fork of youtube-dl).

A word of warning: this application is pretty heavy, and setup can be complicated. If you're looking for something simpler, check out: https://github.com/tubearchivist/tubearchivist

## Usage Instructions (START HERE)

1. Install docker, python3, and docker-compose
2. ./up.sh
   1. respond to requests for input
2. Wait a minute, then visit localhost:80
3. Login as admin/admin
    - note that with the current video approval workflow, non-admin users won't be able to view unapproved videos
    - it's recommended to visit /password-reset immediately to change the admin user's default password if using in an untrusted environment
4. navigate to the archival requests page from the hamburger menu, add a link, and wait a few minutes

That's it for basic usage, and should work. If that doesn't work, bug me on Discord.

## Contributing
Contributions are always welcome. Please see [CONTRIBUTING.md](https://github.com/horahoradev/horahora/blob/master/CONTRIBUTING.md) for details, including an architectural rundown.

## Designs
Designs are listed here:
https://github.com/horahoradev/horahora-designs

## More Detailed Feature List
- performant at 170k videos, even for all varieties of search queries
- support for videos which have been deleted from the origin (e.g. if the original site deletes the video, there's no impact on your instance)
- support for comments, view count, user ratings
- video approval workflow which prevents normal users from seeing videos before they've been approved
- support for TOS/privacy policy
- content archival modeled as one-to-many user subscriptions, so users "subscribe" to a category (link), and links are prioritized according to the number of subscribers
- support for any website supported by yt-dlp which has the required metadata (but I only use YT/nicovideo atm)
- artificial user creation: archived videos will be grouped under a Horahora user created for the archived website's user (e.g. if I archive from Russia Today, then a Russia TOday user will be created on Horahora)
- support for yt-dlp tunneling via Gluetun (see below for setup)
- dark mode toggle

## Advanced Use Cases
### Other Storage Backends (s3, backblaze, anything s3-compatible)
By default, Horahora will store videos locally using Minio.

If you don't want videos to be stored locally, modify .env, adding the relevant values for your use case.

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
