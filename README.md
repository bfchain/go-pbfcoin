## pbfcoin Go

Official golang implementation of the pbfcoin protocol

          | Linux   | OSX | ARM | Windows | Tests
----------|---------|-----|-----|---------|------
## pbfcoin Go

Official golang implementation of the pbfcoin protocol

          | Linux   | OSX | ARM | Windows | Tests
----------|---------|-----|-----|---------|------
develop   | (http://www.paybf.com/)
master    | (http://www.paybf.com/)
## Automated development builds

The following builds are build automatically by our build servers after each push to the [develop](http://www.paybf.com/) branch.

## Building the source

For prerequisites and detailed build instructions please read the
[Installation Instructions](https://github.com/pbfcoin/go-pbfcoin/wiki/Building-pbfcoin)
on the wiki.

Building gpbf requires two external dependencies, Go and GMP.
You can install them using your favourite package manager.
Once the dependencies are installed, run

    make gpbf

## Executables

Go pbfcoin comes with several wrappers/executables found in 
[the `cmd` directory](https://github.com/pbfcoin/go-pbfcoin/tree/develop/cmd):

 Command  |         |
----------|---------|
`gpbf` | pbfcoin CLI (pbfcoin command line interface client) |
`bootnode` | runs a bootstrap node for the Discovery Protocol |
`pbftest` | test tool which runs with the [tests](https://github.com/pbfcoin/tests) suite: `/path/to/test.json > pbftest --test BlockTests --stdin`.
`evm` | is a generic pbfcoin Virtual Machine: `evm -code 60ff60ff -gas 10000 -price 0 -dump`. See `-h` for a detailed description. |
`disasm` | disassembles EVM code: `echo "6001" | disasm` |
`rlpdump` | prints RLP structures |

## Command line options

`gpbf` can be configured via command line options, environment variables and config files.

To get the options available:

    gpbf help

For further details on options, see the [wiki](https://github.com/pbfcoin/go-pbfcoin/wiki/Command-Line-Options)

## Contribution

If you'd like to contribute to go-pbfcoin please fork, fix, commit and
send a pull request. Commits who do not comply with the coding standards
are ignored (use gofmt!). If you send pull requests make absolute sure that you
commit on the `develop` branch and that you do not merge to master.
Commits that are directly based on master are simply ignored.

See [Developers' Guide](https://github.com/pbfcoin/go-pbfcoin/wiki/Developers'-Guide)
for more details on configuring your environment, testing, and
dependency management.
