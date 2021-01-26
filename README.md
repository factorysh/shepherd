ShepherD
========

ShepherD removes exited containers, with a grace time.

ShepherD talks to Sentry when a container is closed with an exit code =! 0.

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
  upperdir    Upperdir
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

```
$ ./bin/shepherd upperdir --help

Explore content of upperdir layer of your containers.

You can pipe result : ./bin/shepherd upperdir| xargs tree -s
Or using JSON output: ./bin/shepherd upperdir -j | jq .

Usage:
  shepherd upperdir [flags]

Aliases:
  upperdir, upper

Flags:
  -a, --all    all containers
  -h, --help   help for upperdir
  -j, --json   json output
```

Build it
--------

    make build

In a CI, or cross compiling

    make docker-build

License
-------

3 terms BSD License
