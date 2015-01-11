watch a number of mountpoints and log their
status to graphite.

I have some cheap USB drive bays that periodically goe offline. This
detects that and makes it easy for me to monitor them and reconnect
them when that happens.

## build

    $ go build .


## configure

copy the sample `config.json` and point it at your own mountpoints,
and specify metric names.

## run

    $ ./mountwatch -config=/path/to/your/config.json -interval=60 \
    -graphite=graphite.example.com:2003 -prefix=yourserver.drives

with the existing config, that would submit metrics something like the
following (if they are all OK):

    server.leo.drives.sata1 0 1420998941
    server.leo.drives.sata2 0 1420998941
    server.leo.drives.sata3 0 1420998941
    server.leo.drives.sata4 0 1420998941
    server.leo.drives.sata5 0 1420998941
    server.leo.drives.sata6 0 1420998941
    server.leo.drives.sata7 0 1420998941
    server.leo.drives.sata8 0 1420998941
    server.leo.drives.sata9 0 1420998941
    server.leo.drives.sata10 0 1420998941

every 60 seconds to your graphite server.

If a drive is unavailable, it will be a '1' instead of a zero for that
metric. If you're alerting off your metrics, that should be pretty
easy to work with.

