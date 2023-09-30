#!/usr/bin/env bash

{ # This ensures the entire script is downloaded. #

if [[ "$EUID" -ne 0 ]]; then
  echo "Please run as root"
  exit 1
fi


latest_version=$(curl -s https://api.github.com/repos/toledompm/cf-ddns/releases?per_page=1 | grep tag_name | cut -d '"' -f 4)

os_name=$(uname -s)
if [[ $os_name == "Darwin" ]]; then
  os_name="darwin"
elif [[ $os_name == "Linux" ]]; then
  os_name="linux"
else
  echo "Unsupported OS"
  exit 1
fi

arch_name=$(uname -m)
if [[ $arch_name == "x86_64" ]]; then
  arch_name="amd64"
elif [[ $arch_name == "aarch64" ]] || [[ $arch_name == "arm64" ]]; then
  arch_name="arm64"
else
  echo "Unsupported architecture"
  exit 1
fi

# Download and install the latest binary
curl -s https://api.github.com/repos/toledompm/cf-ddns/releases/tags/${latest_version} \
  | grep "browser_download_url.*${os_name}-${arch_name}" \
  | cut -d '"' -f 4 \
  | xargs -n 1 curl -Ls -o cf-ddns

chmod +x cf-ddns

mv cf-ddns /usr/local/bin/cf-ddns

# Create config file
mkdir -p /etc/cf-ddns

read -p "Enter your Cloudflare API token: " token

read -p "Enter your Cloudflare zoneId: " zone_id

read -p "Enable IPv6? (y/N): " ipv6_enabled
ipv6_enabled=${ipv6_enabled:-"n"}
if [[ $ipv6_enabled == "y" ]]; then
  ipv6_enabled=true
else
  ipv6_enabled=false
fi

read -p "Enable IPv4? (y/N): " ipv4_enabled
ipv4_enabled=${ipv4_enabled:-"n"}
if [[ $ipv4_enabled == "y" ]]; then
  ipv4_enabled=true
else
  ipv4_enabled=false
fi

read -p "Enter the address to use for fetching IPv6 (default: https://myexternalip.com/json): " ipv6_fetch_address
ipv6_fetch_address=${ipv6_fetch_address:-"https://myexternalip.com/json"}

read -p "Enter the address to use for fetching IPv4 (default: https://myexternalip.com/json): " ipv4_fetch_address
ipv4_fetch_address=${ipv4_fetch_address:-"https://myexternalip.com/json"}

cat > /etc/cf-ddns/config.json <<EOF
  {
    "cloudflare": {
      "token": "${token}",
      "zoneId": "${zone_id}"
    },
    "records": [
      {
        "name": "change-me.example.com",
        "proxy": false
      }
    ],
    "ipv6": {
      "enabled": ${ipv6_enabled},
      "fetchAddress": "${ipv6_fetch_address}"
    },
    "ipv4": {
      "enabled": ${ipv4_enabled},
      "fetchAddress": "${ipv4_fetch_address}"
    }
  }
EOF

read -p "Enable systemd service? (y/N)" systemd_enabled
systemd_enabled=${systemd_enabled:-"n"}
if [[ $systemd_enabled == "y" ]]; then
  systemd_enabled=true
else
  systemd_enabled=false
fi

if [[ $systemd_enabled ]]; then
  curl -s https://raw.githubusercontent.com/toledompm/cf-ddns/main/systemd/cf-ddns.service > /etc/systemd/system/cf-ddns.service
  curl -s https://raw.githubusercontent.com/toledompm/cf-ddns/main/systemd/cf-ddns.timer > /etc/systemd/system/cf-ddns.timer
fi

echo "Change the records in /etc/cf-ddns/config.json and run 'systemctl start cf-ddns' to start the service."

}
