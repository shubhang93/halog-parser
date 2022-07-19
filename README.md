# How to use

## Run the build script

- `./build.sh`
- Copy the emitted binary file `halogp_linux` onto an HAProxy VM using `scp`

```shell
scp halogp_linux <USER-NAME>@<VM-IP>:/home/<USER-NAME>/
```

- ssh into the vm and start a tmux
- Run the following command in the tmux session

```shell
tail -f /var/log/haproxy/haproxy-info.log |  halog -H | ./halogp_linux parse
```

- Detach the tmux session
- Exit the VM

**The parser runs as long as there is input being produced through `STDIN` once there is a stop
of input it stops parsing logs, to manually halt the parser, attach the tmux session and hit
`CTRL+C` to terminate the parsing. The result is streamed to an output file called `out.csv` and not written
all at once to conserve memory. The `out.csv` can be used to further analyze the logs** 