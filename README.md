# revv

To ease my application deployment in cloud, I need a small utility to auto-restart
my application when the binary is updated. It should be more customizable than a
shell script, easier to configure than supervisord, simpler than
[codegangster's gin](https://github.com/codegangsta/gin). So I wrote revv, a tool
to reload program on change.

## Usage

    revv <program> [args...]

## License

MIT