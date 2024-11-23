# 📋 Intranet Check Tool

## 🚀 Overview
The **Intranet Check Tool** is a compact executable for verifying internet connectivity, public IP reachability, WireGuard port status, and Samba server functionality. It features a simple GUI for real-time status updates.

## 🚔 Checks
- **Internet Connectivity**: Tests access by pinging a predefined IP.
- **Public IP Reachability**: Verifies that the public IP responds to pings.
- **WireGuard Port Check**: Confirms UDP port availability on the public IP.
- **Samba Server Validation**:
    - Ping test for server availability.
    - Login test to validate credentials.

## 🤐 `.env`:
   ```plaintext
   INTERNET_TEST_IP=<IP_for_internet_check>
   INTRANET_PUBLIC_IP=<Public_IP_to_ping>
   WIREGUARD_PORT=<WireGuard_Port>
   SAMBA_SERVER_IP=<Samba_Server_IP>
   SAMBA_USER=<Samba_Username>
   SAMBA_PASSWORD=<Samba_Password>
   ```