import requests

# Replace with your actual API key
API_KEY = "99aa356fcd42e1cc4508f99140994c82b859a322eaa21c46d3e35419836b7cade125ba55b7a46f1a"

# Output file path
output_file = "/opt/fortifi/config/malicious_ips.txt"

# AbuseIPDB endpoint
url = "https://api.abuseipdb.com/api/v2/blacklist"

# Parameters for the request
params = {
    "confidenceMinimum": "90",
    "plaintext": ""  # can also use 'Accept: text/plain'
}

# Headers for the request
headers = {
    "Key": API_KEY,
    "Accept": "text/plain"
}

def fetch_and_save_blacklist():
    try:
        response = requests.get(url, headers=headers, params=params)

        if response.status_code == 200:
            ip_list = response.text.strip()

            with open(output_file, "w") as f:
                f.write(ip_list)
            
            print(f"[+] Successfully wrote {len(ip_list.splitlines())} IPs to {output_file}")
        else:
            print(f"[!] Failed to fetch blacklist. Status: {response.status_code}")
            print(response.text)

    except Exception as e:
        print(f"[!] An error occurred: {e}")

if __name__ == "__main__":
    fetch_and_save_blacklist()
