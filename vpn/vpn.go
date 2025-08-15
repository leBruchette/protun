package vpn

import (
	"fmt"
	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
	"os"
	"strings"
	"time"
)

const (
	VpnConfig     = "VPN_CONFIG"
	VpnUser       = "VPN_USER"
	VpnPass       = "VPN_PASS"
	VpnConfigsDir = "ovpn/configs"
)

var vpnProfile, vpnUser, vpnPass, ovpnConfigFile string

func StartSession(startedChan chan struct{}) {
	config := configureVpn()
	session := createNewVpnSession(config)
	session.Start()

	//delay to allow fetching of the ip address
	time.Sleep(3 * time.Second)
	startedChan <- struct{}{}
	err := session.Wait()
	if err != nil {
		fmt.Println("Openvpn3 error: ", err)
	} else {
		fmt.Println("Graceful exit")
	}
}

func configureVpn() openvpn3.Config {
	loadVpnCredentials()
	openvpn3.SelfCheck(performLibraryCheck())
	bytes, err := os.ReadFile(ovpnConfigFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return openvpn3.NewConfig(string(bytes))
}

func loadVpnCredentials() {
	getEnvOrExit := func(envVar string) string {
		val := os.Getenv(envVar)
		if val == "" {
			fmt.Printf("%s environment variable not set, exiting...", envVar)
			os.Exit(1)
		}
		return val
	}

	vpnProfile = getEnvOrExit(VpnConfig)
	vpnUser = getEnvOrExit(VpnUser)
	vpnPass = getEnvOrExit(VpnPass)
	ovpnConfigFile = fmt.Sprintf("%s/%s.ovpn", VpnConfigsDir, vpnProfile)
}

func createNewVpnSession(config openvpn3.Config) *openvpn3.Session {
	return openvpn3.NewSession(config, openvpn3.UserCredentials{
		Username: vpnUser,
		Password: vpnPass,
	}, &LoggingCallbacks{ProfileName: vpnProfile})
}

func performLibraryCheck() StdoutLogger {
	var logger StdoutLogger = func(text string) {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			fmt.Println("Library check >>", line)
		}
	}
	return logger
}
