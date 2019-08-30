# smartlight

A small Linux (! not tested on any other platform) app that dynamically switches keyboard backlight based on mouse and keyboard activity.</br>
When keyboard / mouse are idle for some time, the backlight is turned off.</br>
This duration can be specified as well as the off and on brightness levels.</br>

Note: I've written this app since my thinkpad (T-480s) does not support automatic keyboard brightness control.</br>
So far it has only been tested with that model on Manjaro Linux.

smartlight consists of 3 binaries:

* **smservice**
  * The service that set's the backlight level via dbus (org.freedesktop.UPower.KbdBacklight.SetBrightness)
  * Run it as a systemd --user service or set up it up in your startup applications in your favorite DE.
  * Use command line arguments to configure keyboard / mouse timeout and brightness levels

* **smcli**
  * command line interface to control smservice
  
* **smtray**
  * tray icon with menu to control smservice
  * set up as a startup application in your DE

## Building

See buildscript `build.sh`.

## How to get

go get it :)
```
go get github.com/moson-mo/smartlight
```

or cone git repo
```
git clone https://github.com/moson-mo/smartlight.git
```

or download binaries (you trust me ?) from the [releases](https://github.com/moson-mo/smartlight/releases) page...

## Installation

You can use the install script `install.sh` which copies the binaries to /usr/local/bin/.</br>
Note that you need to build first (or download binries from release on github).

In order to run things on startup set up smservice (and smtray if you like) as startup application on your favorite desktop environment.

## Configuration

On the first start of smservice, a config file is created here `~/.smartlight/config.json`, adjust the settings as needed.
You can also override those settings by running the application with parameters. Start the service with "-h" to see a list of available options.

## Based on

Following libraries are used:

* [godbus](https://github.com/godbus/dbus) - dbus communication
* [gohook](https://github.com/robotn/gohook) - catching keyboard / mouse events
* [systray](https://github.com/getlantern/systray) - tray icon (for smtray)
* [beeep](https://github.com/gen2brain/beeep) - notifications (from smtray)

## Disclamer

Although this software uses "safe methods" in order to set keyboard backlight levels I wouldn't say it can kills things. So handle with care.

In no event will we be liable for any loss or damage including without limitation, indirect or consequential loss or damage, or any loss or damage whatsoever arising from loss of data or profits arising out of, or in connection with, the use of this software.
(ok, that section I've copied from somewhere else ;-) , but ok it's still valid)

## License

What? FOSS brother. MIT

A few libraries are used here.</br>
If this project does not comply, please let me know.</br>
I don't know shit about those things. :) , bear with me.

*WORK IN PROGRESS*