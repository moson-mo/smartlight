package cli

import (
	"fmt"
	netrpc "net/rpc"
	"os"
	"strings"

	"github.com/moson-mo/smartlight/internal/rpc"
)

// Run starts this whole thing :) and runs commands against the RPC server of the service
func Run() {
	if len(os.Args) < 2 || wantsHelp(os.Args[1]) {
		fmt.Println(infoTXT)
		return
	}

	arg := os.Args[1]

	client, err := netrpc.Dial("tcp", "127.0.0.1:31987")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	r := &rpc.Response{}
	err = client.Call("smartlight.Execute", arg, r)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	if arg == "status" {
		fmt.Println("Current status is: " + r.Message)
		return
	}

	fmt.Println("Command executed.\nService returned: " + r.Message)
}

// holy crap, yes, i do want/need help
func wantsHelp(arg string) bool {
	if strings.Contains(arg, "-h") || strings.Contains(arg, "help") || strings.Contains(arg, "?") {
		return true
	}
	return false
}
