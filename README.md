# MiningHQ Miner Manager

The MiningHQ Miner manager GUI

## About

The MiningHQ Miner Manager is a graphical application for interacting
with the miner installed on your local machine.

It interacts with the installed
[Miner Controller](https://github.com/donovansolms/miner-controller) for
showing various statistics and information of your local setup.

On first launch, the manager acts as an installer for your rigs.

## Note

You might wonder why we don't have a single binary that accepts different
flags to install/uninstall/manage the rig. We initially had it working this
way, except when Windows testing started we found that the command line doesn't
work for GUI apps. This made it so that we could no longer route commands
correctly under Windows and decided to break it up into smaller pieces.

We might revisit this in the future, for now it gives a consistent experience.

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
