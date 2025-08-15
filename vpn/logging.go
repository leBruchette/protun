package vpn

import (
	"fmt"
	"github.com/mysteriumnetwork/go-openvpn/openvpn3"
	"os"
	"os/exec"
	"strings"
)

// StdoutLogger represents the stdout logger callback
type StdoutLogger func(text string)

// Log logs the given string to stdout logger
func (lc StdoutLogger) Log(text string) {
	lc(text)
}

type LoggingCallbacks struct {
	ProfileName string
}

var lastWasDNSLabel bool
var dnsServer string
var disableLogs = false

func (lc *LoggingCallbacks) Log(text string) {
	if !disableLogs {
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			fmt.Printf("[%s] Openvpn log >> %s\n", lc.ProfileName, line)
			if strings.Contains(line, "DNS Servers:") && len(dnsServer) == 0 {
				lastWasDNSLabel = true
			} else if lastWasDNSLabel {
				dnsServer = strings.TrimSpace(line)
				lastWasDNSLabel = false
				fmt.Printf("[%s] Openvpn log >> Updating dns server in : %s\n", lc.ProfileName, dnsServer)
				err := updateNameserverConfig(dnsServer) // found in openvpn log output
				if err != nil {
					fmt.Printf("[%s] Openvpn log >> !! Error Updating nameserver config: %s\n", lc.ProfileName, err.Error())
					os.Exit(1)
				}
			} else if strings.Contains(line, "Name:CONNECTED") {
				disableLogs = true
			}
		}
	}
}

func (lc *LoggingCallbacks) OnEvent(event openvpn3.Event) {
	fmt.Printf("[%s] Openvpn event >> %+v\n", lc.ProfileName, event)
}

var statCount = 0

func (lc *LoggingCallbacks) OnStats(stats openvpn3.Statistics) {
	statCount++
	if statCount%50 == 0 {
		fmt.Printf("[%s] Openvpn stats >> %+v\n", lc.ProfileName, stats)
	}
}

func updateNameserverConfig(ip string) error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf(`echo "nameserver %s" | tee /etc/resolv.conf`, ip))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to update /etc/resolv.conf: %w", err)
	}

	return nil
}
