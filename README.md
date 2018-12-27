# MiningHQ Miner Manager

The MiningHQ Miner manager GUI

## About

The MiningHQ Miner Manager is a graphical application for interacting
with the miner installed on your local machine.

It interacts with the installed
[Miner Controller](https://github.com/donovansolms/miner-controller) for
showing various statistics and information of your local setup.

On first launch, the manager acts as an installer for your rigs.

## Building

The following should work for Linux, Windows and MacOS.

1. You must have a working [Go](https://golang.org/) installation
2. Install the required libraries

  ```
  go get -u github.com/asticode/go-astilectron
  go get -u github.com/asticode/go-astilectron-bundler/...
  go get -u github.com/asticode/go-astilectron-bootstrap
  ```

## License

The software is licensed under the GNU GPL v3, you can find the
[full license](LICENSE) in the root of this repository.

---
*A MiningHQ project*
[https://www.mininghq.io](https://www.mininghq.io)
