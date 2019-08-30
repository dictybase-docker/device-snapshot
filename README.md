# device-snapshot
A command line runner to generate webpage snapshots using headless remote chrome browser.

# Usage
It's recommend to use the headless [docker
image](https://hub.docker.com/r/chromedp/headless-shell) as a remote chrome
instance.

## Command

```
NAME:
   snapshot - generate webpage snapshot using remote chrome browser

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-format value                      format of the logging out, either of json or text. (default: "json")
   --log-level value                       log level for the application (default: "error")
   --host value, -H value                  remote host address
   --path value, -p value                  webpage paths for which the snapshots will be taken
   --remote-chrome-host value, --rh value  remote chrome host
   --remote-chrome-port value, --rp value  remote chrome port (default: 9222)
   --output value, -o value                output path for saving all the files
   --help, -h                              show help
   --version, -v                           print the version
```
