package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/joho/godotenv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	internetTestIP      string
	intranetPublicIP    string
	wireguardPort       string
	sambaServerIP       string
	sambaUser           string
	sambaPassword       string
	connectionTimeout   = 3 * time.Second
	checkExecutionDelay = 500 * time.Millisecond
)

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	internetTestIP = os.Getenv("INTERNET_TEST_IP")
	intranetPublicIP = os.Getenv("INTRANET_PUBLIC_IP")
	wireguardPort = os.Getenv("WIREGUARD_PORT")
	sambaServerIP = os.Getenv("SAMBA_SERVER_IP")
	sambaUser = os.Getenv("SAMBA_USER")
	sambaPassword = os.Getenv("SAMBA_PASSWORD")
}

func checkInternetConnection() string {
	cmd := exec.Command("ping", "-c", "1", internetTestIP)
	if os.PathSeparator == '\\' {
		cmd = exec.Command("ping", "-n", "1", internetTestIP)
	}

	err := cmd.Run()
	if err != nil {
		return "❌ Internet Access: Failed"
	}
	return "✅ Internet Access: Successful"
}

func checkPublicIPPing() string {
	cmd := exec.Command("ping", "-c", "1", intranetPublicIP)
	if os.PathSeparator == '\\' {
		cmd = exec.Command("ping", "-n", "1", intranetPublicIP)
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("❌ Public IP (%s): Ping Failed", intranetPublicIP)
	}
	return fmt.Sprintf("✅ Public IP (%s): Ping Successful", intranetPublicIP)
}

func checkWireguardPortUDP() string {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", intranetPublicIP, wireguardPort))
	if err != nil {
		return fmt.Sprintf("❌ WireGuard Port (UDP %s) on Public IP (%s): Not Accessible", wireguardPort, intranetPublicIP)
	}
	conn.Close()
	return fmt.Sprintf("✅ WireGuard Port (UDP %s) on Public IP (%s): Accessible", wireguardPort, intranetPublicIP)
}

func checkSambaServerPing() string {
	cmd := exec.Command("ping", "-c", "1", sambaServerIP)
	if os.PathSeparator == '\\' {
		cmd = exec.Command("ping", "-n", "1", sambaServerIP)
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("❌ Ping to (%s): Failed", sambaServerIP)
	}
	return fmt.Sprintf("✅ Ping to (%s): Successful", sambaServerIP)
}

func checkSambaLogin() string {
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:445", sambaServerIP))
	if err != nil {
		return fmt.Sprintf("❌ Samba Login (%s): Failed (Connection Error: %v)", sambaServerIP, err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     sambaUser,
			Password: sambaPassword,
		},
	}

	client, err := d.Dial(conn)
	if err != nil {
		return fmt.Sprintf("❌ Samba Login (%s): Failed (Authentication Error: %v)", sambaServerIP, err)
	}
	defer client.Logoff()

	return fmt.Sprintf("✅ Samba Login (%s): Successful", sambaServerIP)
}

func main() {
	loadEnvVariables()

	a := app.New()
	w := a.NewWindow("Intranet Check")

	labels := []*widget.Label{
		widget.NewLabel("⌛ Internet Access: Running..."),
		widget.NewLabel("⌛ Public IP Ping: Running..."),
		widget.NewLabel("⌛ WireGuard Port: Running..."),
		widget.NewLabel(fmt.Sprintf("⌛ Ping to %s: Running...", sambaServerIP)),
		widget.NewLabel("⌛ Samba Login: Running..."),
	}

	var objects []fyne.CanvasObject
	for _, label := range labels {
		objects = append(objects, label)
	}

	content := container.NewVBox(objects...)
	w.SetContent(content)
	w.Resize(fyne.NewSize(400, 300))
	w.Show()

	go func() {
		checks := []struct {
			checker func() string
			label   *widget.Label
		}{
			{checkInternetConnection, labels[0]},
			{checkPublicIPPing, labels[1]},
			{checkWireguardPortUDP, labels[2]},
			{checkSambaServerPing, labels[3]},
			{checkSambaLogin, labels[4]},
		}

		for _, check := range checks {
			time.Sleep(checkExecutionDelay)
			result := check.checker()
			check.label.SetText(result)
		}
	}()

	a.Run()
}
