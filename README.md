# CF-DDNS
ddns service for Cloudflare. Currently only supports A and AAAA records.

## Usage

Install with curl
```
curl -s https://raw.githubusercontent.com/toledompm/cf-ddns/install.sh | sudo sh
```

### Configuration
cf-ddns is configured using a json file. The default location is `/etc/cf-ddns/config.json`. You can specify a different location using the `-c` flag.

```json
{
  "cloudflare": {
    "token": "...", // Cloudflare API token, requires Zone.DNS.Edit permissions
    "zoneId": "..."
  },
  "records": [
    {
      "name": "ddns.toledompm.xyz",
      "proxy": true
    },
  ],
  "ipv6": {
    "enabled": true,
    "fetchAddress": "https://myexternalip.com/json" // URL to fetch IPv6 address from, must return a json object with an "ip" field.
  },
  "ipv4": {
    "enabled": true,
    "fetchAddress": "https://myexternalip.com/json"
  }
}
```