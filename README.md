Shepherd
=======

Shepherd removes exited containers, with a grace time.

Shepherd talks to Sentry when a container is closed with an exit code =! 0.

Help
----

```
$ ./bin/shepherd -h
shepherd cleans the mess

Usage:
  shepherd [command]

Available Commands:
  event       Displays the report for a container, as it should be sent to Sentry
  help        Help about any command
  version     Version of shepherd
  watch       Watch docker and clean its mess

Flags:
  -h, --help   help for shepherd

Use "shepherd [command] --help" for more information about a command.
```

```
$ ./bin/shepherd watch --help
Watch docker and clean its mess

Usage:
  shepherd watch [flags]

Flags:
  -a, --admin string    Listen admin http address (default "localhost:4012")
  -c, --config string   config file
  -h, --help            help for watch
```

Build it
--------

    make build

In a CI, or cross compiling

    make docker-build

License
-------

3 terms BSD License