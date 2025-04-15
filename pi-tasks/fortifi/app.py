import os
import subprocess
import json
import tempfile
import requests  # Ensure requests is imported at the top of your file
from flask import Flask, render_template, request, redirect, url_for, jsonify

app = Flask(__name__)

FORTIFI_AUTH_URL = "http://3.101.107.33:3000/PiInit"  # Replace with the actual API URL
FORTIFI_REFRESH_URL = "http://3.101.107.33:3000/RefreshPi"  # Refresh token endpoint

#FORTIFI_AUTH_URL = "http://192.168.200.150:3000/PiInit"  # Replace with the actual API URL
#FORTIFI_REFRESH_URL = "http://192.168.200.150:3000/RefreshPi"  # Refresh token endpoint

#CONFIG_DIR = r"/home/fortifi/FortiFi-Frontend/config"
CONFIG_DIR = r"/opt/fortifi/config"
SETUP_STATE_PATH = os.path.join(CONFIG_DIR, "setup_state.json")
WIFI_CREDENTIALS_PATH = os.path.join(CONFIG_DIR, "wifi_credentials.txt")
TOKEN_STORAGE_PATH = os.path.join(CONFIG_DIR, "tokens.json")
OUTPUT_JSON = '/opt/fortifi/models/model-output/all_predictions.json'

os.makedirs(CONFIG_DIR, exist_ok=True)

def refresh_token():
    """Call /refresh_token before making API requests."""
    refresh_url = "http://localhost:5000/refresh_token"
    try:
        response = requests.get(refresh_url)
        if response.status_code == 200:
            new_token_data = response.json()
            new_jwt = new_token_data.get("jwt", "")
            if new_jwt:
                return new_jwt
        print(f"Failed to refresh token: {response.status_code}, {response.text}")
    except Exception as e:
        print(f"Error refreshing token: {e}")
    return load_tokens()

def read_setup_state():
    if os.path.exists(SETUP_STATE_PATH):
        with open(SETUP_STATE_PATH, "r") as file:
            return json.load(file)
    else:
        # Explicitly write "not_started" if no file exists
        setup_state = {"setup_step": "not_started"}
        with open(SETUP_STATE_PATH, "w") as file:
            json.dump(setup_state, file)
        return setup_state

def update_setup_state(step):
    setup_state = read_setup_state()
    setup_state["setup_step"] = step
    with open(SETUP_STATE_PATH, "w") as file:
        json.dump(setup_state, file)

def get_system_uuid():
    """Retrieve the system UUID using dmidecode."""
    try:
        return subprocess.run(
            ['sudo', '/usr/sbin/dmidecode', '-s', 'system-uuid'],
            capture_output=True, text=True, check=True
        ).stdout.strip()
    except subprocess.CalledProcessError as e:
        return None  # Return None if the UUID retrieval fails

# Update your reset_setup function to remove cron jobs when resetting the device
def reset_setup():
    if os.path.exists(SETUP_STATE_PATH):
        os.remove(SETUP_STATE_PATH)
    if os.path.exists(WIFI_CREDENTIALS_PATH):
        os.remove(WIFI_CREDENTIALS_PATH)
    if os.path.exists(TOKEN_STORAGE_PATH):
        os.remove(TOKEN_STORAGE_PATH)

    # Remove cron jobs
    modify_cron(add_cron=False)

    return jsonify({"message": "Device reset successfully"}), 200

def modify_cron(add_cron=True):
    cron_job_fortifi = "* * * * * /opt/fortifi/scripts/run_fortifi.sh >> /var/log/run_models.log 2>&1"
    cron_job_token_refresh = "0 0 * * 0 /usr/bin/curl -X GET http://localhost:5000/reset_weekly_distribution >> /var/log/reset_weekly_distribution.log 2>&1"

    # Get current crontab lines
    process = subprocess.run(['sudo', 'crontab', '-u', 'root', '-l'], capture_output=True, text=True)
    cron_lines = process.stdout.splitlines() if process.returncode == 0 else []

    # Remove existing entries related to Fortifi or token refresh
    cron_lines = [line for line in cron_lines if 'run_fortifi.sh' not in line and 'reset_weekly_distribution' not in line]

    # Add back cron jobs only if add_cron is True
    if add_cron:
        cron_lines.extend([cron_job_fortifi, cron_job_token_refresh])

    # Write modified crontab to a temp file and load it
    with tempfile.NamedTemporaryFile(mode='w+', delete=False) as tmp_cron:
        tmp_cron.write('\n'.join(cron_lines) + '\n')
        tmp_cron_path = tmp_cron.name

    subprocess.run(['sudo', 'crontab', '-u', 'root', tmp_cron_path])
    os.remove(tmp_cron_path)

def calculate_week_total():
    """Calculate the total number of predictions from all_predictions.json."""
    if os.path.exists(OUTPUT_JSON):
        with open(OUTPUT_JSON, "r") as file:
            cumulative_predictions = json.load(file)
        return sum(cumulative_predictions.values())
    return 0

def load_tokens():
    """Load stored JWT and refresh tokens."""
    if os.path.exists(TOKEN_STORAGE_PATH):
        with open(TOKEN_STORAGE_PATH, "r") as file:
            return json.load(file)
    return {}

def save_tokens(jwt_token, refresh_token):
    """Save JWT and refresh tokens."""
    with open(TOKEN_STORAGE_PATH, "w") as file:
        json.dump({"jwt": jwt_token, "refresh": refresh_token}, file)


@app.route('/')
def home():
    setup_state = read_setup_state()
    if setup_state.get("setup_step") == "completed":
        return render_template('index.html', show_reset=True, show_resume=False)
    elif setup_state.get("setup_step") in ["wifi_setup", "internet_connected"]:
        return render_template('index.html', show_reset=True, show_resume=True, setup_step=setup_state["setup_step"])
    return render_template('index.html', show_reset=False, show_resume=False)


@app.route('/wifi', methods=['GET', 'POST'])
def wifi():
    if request.method == 'POST':
        ssid = request.form.get('ssid')
        password = request.form.get('password')

        if ssid and password:
            with open(WIFI_CREDENTIALS_PATH, 'w') as f:
                f.write(f"SSID: {ssid}\nPassword: {password}")

            return "", 200  # Success

        return "Error: Missing WiFi credentials", 400

    return render_template('wifi.html')

@app.route('/register')
def register():
    return render_template('register.html')

@app.route('/get_setup_state')
def get_setup_state():
    return jsonify(read_setup_state())

@app.route('/update_setup_state', methods=["POST"])
def update_setup():
    data = request.json
    if "setup_step" in data:
        update_setup_state(data["setup_step"])
        
        # Verify the state was updated
        return jsonify({"message": "Setup state updated"}), 200
    
    return jsonify({"error": "Invalid request"}), 400

@app.route('/refresh_token', methods=['GET'])
def refresh_token():
    tokens = load_tokens()
    refresh_token = tokens.get("refresh")

    try:
        device_uuid = get_system_uuid()
    except subprocess.CalledProcessError as e:
        return jsonify({"error": f"Failed to retrieve system UUID: {str(e)}"}), 500


    headers = {"Refresh": refresh_token if refresh_token else tokens.get("jwt", "")}
    try:
        response = requests.get(
            FORTIFI_REFRESH_URL,
            headers=headers,
            params={"id": device_uuid}
        )

        if response.status_code == 200:
            new_jwt = response.headers.get("Jwt")
            new_refresh = response.headers.get("Refresh")

            if new_jwt and new_refresh:
                save_tokens(new_jwt, new_refresh)
                return jsonify({"message": "Token refreshed successfully", "jwt": new_jwt}), 200
            else:
                return jsonify({"error": "Missing tokens in response"}), 500
        else:
            return jsonify({"error": f"Failed to refresh token: {response.status_code}, {response.text}"}), response.status_code

    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500

@app.route("/reset_weekly_distribution", methods=["POST"])
def reset_weekly_distribution():
    """Reset the weekly distribution by sending a request to the external API."""
    jwt_token = refresh_token()
    if not jwt_token:
        return jsonify({"error": "JWT token missing"}), 401

    week_total = calculate_week_total()
    data = {"week_total": week_total}
    headers = {"Authorization": f"Bearer {jwt_token}", "Content-Type": "application/json"}
    
    api_url = "http://3.101.107.33:3000/ResetWeeklyDistribution"
    #api_url = "http://192.168.200.150/ResetWeeklyDistribution"
    response = requests.post(api_url, json=data, headers=headers)

    if response.status_code == 200:
        return jsonify({"message": "Weekly distribution reset successfully."}), 200
    elif response.status_code == 401:
        return jsonify({"error": "Unauthorized - check JWT token."}), 401
    elif response.status_code == 404:
        return jsonify({"error": "User not found in database."}), 404
    elif response.status_code == 405:
        return jsonify({"error": "Method not allowed - check request method."}), 405
    elif response.status_code == 500:
        return jsonify({"error": "Internal server error - check logs."}), 500
    else:
        return jsonify({"error": f"Unexpected error: {response.status_code}, {response.text}"}), response.status_code

@app.route('/check_internet')
def check_internet():
    try:
        response = requests.get('https://www.google.com', timeout=5)
        if response.status_code == 200:
            return jsonify({"status": "connected"}), 200
        else:
            return jsonify({"status": "not connected"}), 503
    except requests.RequestException:
        return jsonify({"status": "not connected"}), 503


# Update your authenticate function to start cron jobs on successful authentication
@app.route('/authenticate', methods=['POST'])
def authenticate():
    try:
        device_id = get_system_uuid()
    except subprocess.CalledProcessError as e:
        return jsonify({"error": f"Failed to retrieve system UUID: {str(e)}"}), 500
    if not device_id:
        return jsonify({"error": "Missing device ID"}), 400

    headers = {"Content-type": "application/json", "Accept": "*/*"}

    try:
        response = requests.post(
            "http://3.101.107.33:3000/PiInit",
            #"http://192.168.200.150:3000/PiInit",
            json={"id": device_id},
            headers=headers
        )

        if response.status_code == 200:
            jwt_token = response.headers.get("Jwt")
            refresh_token = response.headers.get("Refresh")

            if jwt_token and refresh_token:
                with open(TOKEN_STORAGE_PATH, "w") as file:
                    json.dump({"jwt": jwt_token, "refresh": refresh_token}, file)

                return jsonify({"jwt": jwt_token, "refresh": refresh_token}), 200
            else:
                return jsonify({"error": "Missing tokens in response"}), 500

        return jsonify({"error": f"Authentication failed: {response.text}"}), response.status_code

    except requests.exceptions.RequestException as e:
        return jsonify({"error": str(e)}), 500
    
@app.route('/get_device_uuid')
def get_device_uuid():
    try:
        # Run dmidecode command to get system UUID
        device_uuid = get_system_uuid()
        return jsonify({"uuid": device_uuid}), 200
    except subprocess.CalledProcessError as e:
        return jsonify({"error": f"Failed to retrieve system UUID: {str(e)}"}), 500
    
@app.route('/start_services', methods=['POST'])
def start_services():
    """Start cron jobs after successful user registration."""
    try:
        modify_cron(add_cron=True)
        return jsonify({"message": "Services started successfully"}), 200
    except Exception as e:
        return jsonify({"error": f"Failed to start services: {str(e)}"}), 500

@app.route('/reset_device', methods=['POST'])
def reset_device():
    reset_setup()
    return redirect(url_for('home'))

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=True)
