# gottem!

Trick a backup program that doesn't support ignore files into backing up whatever you like

## What/Why?

Google Backup and Sync, desktop app for Google One, does not support ignore files.
That means I can't tell it not to touch `node_modules/` etc. when backing up my
projects directory.

gottem takes all of the files in a directory, minus the ones you don't want,
and symlinks them elsewhere so that you can point your backup software there.

In my case, I want to back up `/Users/adam/Mac Projects`, so I use the following
`.gottemconfig` and point GB&S to `/Users/adam/Backup`:

```
# From
/Users/adam/Documents/Mac Projects
# To
/Users/adam/Backup

# Ignore rules
node_modules
.app/Contents
build
Build
```

## How to use it

Download a release or build with `go build gottem`, and place a `.gottemconfig`
in your home directory. The first (non-comment) line is taken as the 'from' directory,
the second as the 'to' directory, and the rest are rules for substrings that can't
exist in a synced file/folder's path.

Now just run gottem!

I have mine hooked up on a cronjob every hour.

## Speed/Technical stuff

For my over half a million project files (518,474), running gottem with the above
config takes 67 seconds on a 3.6GHz CPU with 16 threads.

```
$ time go run gottem
Gottem is linking /Users/adam/Documents/Mac Projects to /Users/adam/Backup...
Done, enjoy your backup!
go run gottem  2.92s user 67.09s system 363% cpu 19.253 total
```

It uses goroutines to take advantage of all your CPU's threads.

**NOTE:** You can pass a command line argument to limit the number of goroutines
(it defaults to 16) eg. `./gottem 8`

Don't worry, gottem isn't copying any files or taking up extra disk space, it's
just tricking programs into thinking it is. For this reason Finder etc. might show
the synced folder as taking space, but if you check your actual disk usage, it will
not go up.
