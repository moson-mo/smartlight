package tray

import (
	"fmt"
	netrpc "net/rpc"
	"strings"

	"github.com/gen2brain/beeep"

	"github.com/moson-mo/smartlight/internal/rpc"

	"github.com/getlantern/systray"
)

var (
	icons           map[string][]byte
	iconErrPath     string
	iconSuccessPath string
)

// Run starts this whole mess :)
func Run() {
	path, err := installCacheFiles()
	if err != nil {
		fmt.Println("Failed to install cache files: " + err.Error())
	}
	iconErrPath = path + "error.svg"
	iconSuccessPath = path + "ok.svg"
	systray.Run(makeTrayIcon, killTrayIcon)
}

// creates tray icon & calls function for building menu
func makeTrayIcon() {
	var err error
	icons, err = getIcons()
	if err != nil {
		notifyError(err.Error())
		return
	}

	systray.SetIcon(icons["ion"])
	systray.SetTitle("smartlight")
	systray.SetTooltip("Control your keyboard backlight")

	buildMenu()
}

// builds the menu and registers click events (what happens when an item is clicked)
func buildMenu() {

	mir := systray.AddMenuItem("Re-connect", "Tries to reconnect to the service")
	mie := systray.AddMenuItem("Enable", "Enables the smartlight service")
	mid := systray.AddMenuItem("Disable", "Disables the smartlight service")
	systray.AddSeparator()
	miquit := systray.AddMenuItem("Quit", "Quit systray app")

	setItemStatus(mie, mid, mir)

	// disable entry (stop)
	go func() {
		for {
			<-mid.ClickedCh
			msg, err := callRPCFunc("stop")
			if err != nil {
				notifyError(err.Error())
				continue
			}
			notifySuccess(msg)
			mid.Disable()
			mie.Enable()
			systray.SetIcon(icons["ioff"])
		}
	}()

	// enable entry (start)
	go func() {
		for {
			<-mie.ClickedCh
			msg, err := callRPCFunc("start")
			if err != nil {
				notifyError(err.Error())
				continue
			}
			notifySuccess(msg)
			mie.Disable()
			mid.Enable()
			systray.SetIcon(icons["ion"])
		}
	}()

	// re-connect enty
	go func() {
		for {
			<-mir.ClickedCh
			setItemStatus(mie, mid, mir)
		}
	}()

	// quit entry
	go func() {
		<-miquit.ClickedCh
		systray.Quit()
	}()
}

// error notification
func notifyError(err string) {
	beeep.Notify("smartlight - Error", err, iconErrPath)
}

// success notification (hopefully more than the one before :) )
func notifySuccess(msg string) {
	beeep.Notify("smartlight - Success", "smartlight service "+msg, iconSuccessPath)
}

// set menu item status according to service status / state, status, what, who cares :)
func setItemStatus(mie, mid, mir *systray.MenuItem) {
	status, err := callRPCFunc("status")
	if err != nil {
		if strings.Contains(err.Error(), "connection refused") {
			notifyError(err.Error())
			mie.Disable()
			mid.Disable()
			systray.SetIcon(icons["ierr"])
			return
		}
		systray.Quit()
	}
	if status == "started" {
		mie.Disable()
		mid.Enable()
		mir.Hide()
		systray.SetIcon(icons["ion"])
		return
	}
	mid.Disable()
	mie.Enable()
	mir.Hide()
	systray.SetIcon(icons["ioff"])
}

// communicate with the service
func callRPCFunc(com string) (string, error) {
	client, err := netrpc.Dial("tcp", "127.0.0.1:31987")
	if err != nil {
		return "", err
	}
	defer client.Close()

	r := &rpc.Response{}
	err = client.Call("smartlight.Execute", com, r)
	if err != nil {
		return "", err
	}
	return r.Message, nil
}

// doesn't so shit atm
func killTrayIcon() {

}
