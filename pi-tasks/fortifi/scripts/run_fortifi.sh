#!/bin/bash

LOCKFILE="/var/run/run_fortifi.lock"
BSSID_FILE="/tmp/bssid_channel.txt"
CRED_FILE="/opt/fortifi/config/static-wifi-credentials.txt"

# Function to clean up background processes when script is terminated
cleanup() {
    echo "[!] Terminating script. Cleaning up running processes..."
    
    echo "[*] Stopping packet capture on WiFi Pineapple..."
    ssh -i "$PINEAPPLE_SSH_KEY" "$PINEAPPLE_USER@$PINEAPPLE_IP" "pkill -f airodump-ng"

    echo "[*] Stopping local TShark capture..."
    pkill -f "tshark"

    rm -f "$LOCKFILE"
    echo "[!] Cleanup complete. Exiting."
    exit 0
}
trap cleanup SIGINT SIGTERM

# Check if the lock file exists and the process is running
if [ -f "$LOCKFILE" ]; then
  PID=$(cat "$LOCKFILE")
  if [ -d "/proc/$PID" ]; then
    echo "Another instance is already running (PID $PID). Exiting."
    exit 1
  fi
fi

# Write the current PID to the lock file
echo $$ > "$LOCKFILE"

sleep 5
export PATH="/usr/bin:/usr/sbin:/sbin:/bin"

CAPTURE_DIR="/opt/fortifi/packet_capture"
PINEAPPLE_CAPTURE_DIR="$CAPTURE_DIR/pineapple"
TSHARK_CAPTURE_DIR="$CAPTURE_DIR/tshark"
mkdir -p "$CAPTURE_DIR" "$PINEAPPLE_CAPTURE_DIR" "$TSHARK_CAPTURE_DIR"

PINEAPPLE_IP="172.16.42.1"
PINEAPPLE_USER="root"
PINEAPPLE_SSH_KEY="$HOME/.ssh/pineapple_key"

if [[ ! -f "$CRED_FILE" ]]; then
    echo "[!] Credential file not found: $CRED_FILE"
    exit 1
fi

TARGET_SSID=$(grep -i "^SSID:" "$CRED_FILE" | awk -F ':' '{print $2}' | xargs)
WIFI_PASSWORD=$(grep -i "^Password:" "$CRED_FILE" | awk -F ':' '{print $2}' | xargs)

if [[ -z "$TARGET_SSID" || -z "$WIFI_PASSWORD" ]]; then
    echo "[!] Failed to extract SSID or Password from $CRED_FILE"
    exit 1
fi

if ! pgrep -f "python3 fortifi.py" > /dev/null; then
    echo "Starting fortifi.py..."
    nohup python3 -u /opt/fortifi/scripts/fortifi.py >> /var/log/fortifi.log 2>&1 &
fi

echo "[*] Configuring WiFi Pineapple for packet capture..."
ssh -i "$PINEAPPLE_SSH_KEY" "$PINEAPPLE_USER@$PINEAPPLE_IP" << EOF
    airmon-ng check kill
    airmon-ng start wlan1

    echo "[*] Scanning for target SSID: '$TARGET_SSID' for 10 seconds..."
    rm -f /tmp/scan-01.csv
    airodump-ng --write /tmp/scan --output-format csv wlan1mon &
    PID=\$!
    sleep 10
    kill \$PID

    TARGET_LINE=\$(grep "$TARGET_SSID" /tmp/scan-01.csv | head -n 1)

    if [[ -z "\$TARGET_LINE" ]]; then
        echo "[!] Failed to find target SSID: $TARGET_SSID"
        exit 1
    fi

    BSSID=\$(echo "\$TARGET_LINE" | awk -F ',' '{print \$1}' | tr -d ' ')
    CHANNEL=\$(echo "\$TARGET_LINE" | awk -F ',' '{print \$4}' | tr -d ' ')

    if [[ -z "\$BSSID" || -z "\$CHANNEL" ]]; then
        echo "[!] Error extracting BSSID or Channel."
        exit 1
    fi

    echo "[+] Found target SSID: $TARGET_SSID"
    echo "[+] BSSID: \$BSSID"
    echo "[+] Channel: \$CHANNEL"
    echo "\$BSSID \$CHANNEL" > $BSSID_FILE

    echo "[*] Setting WiFi adapter to channel \$CHANNEL..."
    iwconfig wlan1mon channel "\$CHANNEL"
EOF

# Capture Loop
while true; do
  TIMESTAMP=$(date '+%Y-%m-%d_%H-%M-%S')
  OUTPUT_DIR="$TSHARK_CAPTURE_DIR/capture_$TIMESTAMP"
  PINEAPPLE_OUTPUT_DIR="$PINEAPPLE_CAPTURE_DIR/capture_$TIMESTAMP"
  mkdir -p "$OUTPUT_DIR" "$PINEAPPLE_OUTPUT_DIR"

  REMOTE_CAPTURE_PATH="/tmp/capture_$TIMESTAMP"

  echo "[*] Starting TShark capture..."
  timeout 30 tshark -i br0 -w "$OUTPUT_DIR/capture.pcap" &
  TSHARK_PID=$!

  echo "[*] Restarting WiFi Pineapple packet capture..."
  ssh -i "$PINEAPPLE_SSH_KEY" "$PINEAPPLE_USER@$PINEAPPLE_IP" << EOF
    if [[ ! -f $BSSID_FILE ]]; then
        echo "[!] BSSID and Channel file not found. Exiting capture."
        exit 1
    fi

    read BSSID CHANNEL < $BSSID_FILE

    echo "[*] Using stored values - BSSID: \$BSSID, Channel: \$CHANNEL"
    airodump-ng --bssid "\$BSSID" -c "\$CHANNEL" -w "$REMOTE_CAPTURE_PATH" wlan1mon &
    CAP_PID=\$!
    sleep 5
    aireplay-ng --deauth 10 -a "\$BSSID" wlan1mon
    sleep 25
    kill \$CAP_PID

    echo "[*] Attempting decryption with airdecap-ng..."
    airdecap-ng -e "$TARGET_SSID" -p "$WIFI_PASSWORD" "${REMOTE_CAPTURE_PATH}-01.cap"

    if [[ -f "${REMOTE_CAPTURE_PATH}-01-dec.cap" ]]; then
        echo "[+] Decryption successful. Decrypted file is ready."
    else
        echo "[!] Decryption failed or no handshake captured."
        exit 1
    fi
EOF

    wait $TSHARK_PID

    echo "[*] Fetching decrypted capture from WiFi Pineapple..."
    scp -i "$PINEAPPLE_SSH_KEY" "$PINEAPPLE_USER@$PINEAPPLE_IP:${REMOTE_CAPTURE_PATH}-01-dec.cap" "$PINEAPPLE_OUTPUT_DIR/capture.pcap"

    if [[ $? -eq 0 ]]; then
        echo "[+] Pineapple capture copied to: $PINEAPPLE_OUTPUT_DIR/capture.pcap"

        echo "[*] Cleaning up capture files on WiFi Pineapple..."
        TIMESTAMP_PREFIX=$(basename "$REMOTE_CAPTURE_PATH")  # e.g., capture_2025-02-29_22-30-00 → capture_2025-02-29_22
        TIMESTAMP_PREFIX_SHORT=$(echo "$TIMESTAMP_PREFIX" | cut -d'_' -f1-3 | tr '_' '-')

        ssh -i "$PINEAPPLE_SSH_KEY" "$PINEAPPLE_USER@$PINEAPPLE_IP" "rm -f /tmp/${TIMESTAMP_PREFIX}*"
        echo "[+] Removed files matching: /tmp/${TIMESTAMP_PREFIX}*"
    else
        echo "[!] Failed to retrieve decrypted Pineapple capture."
        continue
    fi

  echo "[*] Converting both captures to JSON..."
  tshark -r "$OUTPUT_DIR/capture.pcap" -T json > "$OUTPUT_DIR/capture.json"
  tshark -r "$PINEAPPLE_OUTPUT_DIR/capture.pcap" -T json > "$PINEAPPLE_OUTPUT_DIR/capture.json"

  echo "[*] Running Zeek on both captures..."
  (
    cd "$OUTPUT_DIR"
    /opt/zeek/bin/zeek -r capture.pcap || echo "[!] Zeek failed on TShark capture"
  )
  (
    cd "$PINEAPPLE_OUTPUT_DIR"
    /opt/zeek/bin/zeek -r capture.pcap || echo "[!] Zeek failed on Pineapple capture"
  )

  echo "[*] Sleeping for 1 second before next capture..."
  sleep 1
done

rm -f "$LOCKFILE"
