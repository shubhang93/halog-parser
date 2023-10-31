# How to use

## Prerequisites

- Requires golang `v1.18` to build

## Run the build script

- ssh into the vm and start a tmux
- Run the following command in the tmux session

```shell
tail -f /var/log/haproxy/haproxy-info.log |  halog -H | ./halogp_linux parse
```

**The parser runs as long as there is input being produced through `STDIN` once there is a stop
of input it stops parsing logs, to manually halt the parser, attach the tmux session and hit
`CTRL+C` to terminate the parsing. The result is streamed to an output file called `out.csv` and not written
all at once to conserve memory. The `out.csv` can be used to further analyze the logs** 
