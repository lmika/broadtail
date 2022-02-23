# broadtail

Broadtail is a web-app for subscribing to YouTube RSS feeds and downloading them using youtube-dl.
Broadtail runs as a daemon, similar to a Plex server.

## Installation

Broadtail is distributed as either a DEB or RPM package and requires a Linux x86 machine:

1. Install [youtube-dl](https://github.com/ytdl-org/youtube-dl).  Broadtail does not install this
   automatically, and makes no assumptions as to how it's installed (this is because the version
   of youtube-dl distributed via package managers may be out of date).
2. Download the DEB or RPM package from the latest release.
3. Install the package using either `apt` or `yum`, depending on the Linux OS you're using.

Broadtail is configured as a systemd daemon.  To start Broadtail, run the following as root:

```
systemctl start broadtail
```

If Broadtail has not started after installation, you may need to make some configuration changes.  See below.

## Development

If you're insterested in doing work on Broadtail, or just testing it out before installing it, you
can check it out and run it locally using the following procedure:

1. Install [Go](https://go.dev) and Node/NPM
2. Install [RWT](https://github.com/lmika/rwt): `go install -x github.com/lmika/rwt`
2. Check-out the latest version of `main`
3. Run `make init`
4. Run `make run`

The web-app will be running on port 3690.

## Configuration

The configuration for Broadtail is located at `/usr/local/etc/broadtail/config.yaml`.  To make changes, open
this up in a text editor and restart Broadtail.  The configuration should hopefully be self explainatory, but
certain key parameters are explained in detail below:

- `youtubedl_command`: This is command Broadtail executes when it needs to use youtube-dl.  This may need to
  be changed depending on how youtube-dl is installed.  By default Broadtail assumes that youtube-dl is installed
  via pip3.
- `data_dir`: This is the directory used by Broadtail to store any internal data or index files.  This does not include
  video files.
- `library_dir`: This is the base directory that Broadtail will use to store downloaded videos.  This should generally
  be on a partition with ample space.  This can be a directory configured as a Plex library.
- `library_owner`: This is the OS user that owns the library.  If set, Broadtail will attempt to change the owner of
  any downloaded videos or subdirectories to this user.  This should generally be the standard non-root user that owns
  the library directory.

## Known Limitations

- Newly created nested subdirectories will not be changed to the library owner.
- Errors during downloaded videos will not be retried.
