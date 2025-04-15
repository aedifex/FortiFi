import joblib
import pandas as pd
import os
import shutil
import time
import json
import requests
import numpy as np
from collections import defaultdict
from collections import Counter
from scipy.stats import entropy
from datetime import datetime, timedelta
from sklearn.preprocessing import StandardScaler, LabelEncoder

#API_URL = "http://192.168.200.150:3000/NotifyIntrusion"
API_URL = "http://3.101.107.33:3000/NotifyIntrusion"

CONFIG_DIR = "/opt/fortifi/config"
TSHARK_DIR = "/opt/fortifi/packet_capture/tshark"
PINEAPPLE_DIR = "/opt/fortifi/packet_capture/pineapple"
PROCESSED_TSHARK_DIR = os.path.join(TSHARK_DIR, "processed")
PROCESSED_PINEAPPLE_DIR = os.path.join(PINEAPPLE_DIR, "processed")
ASSET_TRACKING_DIR = '/opt/fortifi/models/model-output/asset-tracking'
IOT_JSON_PATH = '/opt/fortifi/models/model-output/iot_devices.json'

IOT_LABEL_MAPPING = {
    "Smart Things": "Smart Things",
    "Amazon Echo": "Speaker",
    "Netatmo Welcome": "Camera",
    "TP-Link Day Night Cloud Camera": "Camera",
    "Samsung SmartCam": "Camera",
    "Dropcam": "Camera",
    "Insteon Camera": "Camera",
    "Withings Smart Baby Monitor": "Baby Monitor",
    "Belkin Wemo Switch": "Switch",
    "TP-Link Smart plug": "Plug",
    "iHome": "iHome",
    "Belkin Wemo Motion Sensor": "Motion Sensor",
    "NEST Protect Smoke Alarm": "Smoke Alarm",
    "Netatmo Weather Station": "Weather Station",
    "Withings Aura Smart Sleep Sensor": "Sleep Sensor",
    "Light Bulbs LiFX Smart Bulbs": "Lightbulb",
    "Triby Speaker": "Speaker",
    "PIX-STAR Photo-frame": "Photo-Frame",
    "HP Printer": "Printer",
    "Nest Dropcam": "Camera",
    "TPLink Router Bridge LAN (Gateway)": "Router"
}

# Ensure necessary directories exist
os.makedirs(PROCESSED_TSHARK_DIR, exist_ok=True)
os.makedirs(PROCESSED_PINEAPPLE_DIR, exist_ok=True)
os.makedirs(ASSET_TRACKING_DIR, exist_ok=True)

TOKEN_STORAGE_PATH = os.path.join(CONFIG_DIR, "tokens.json")
CUMULATIVE_JSON_PATH = '/opt/fortifi/models/model-output/all_predictions.json'

if os.path.isfile(CUMULATIVE_JSON_PATH):
    with open(CUMULATIVE_JSON_PATH, 'r') as file:
        cumulative_predictions = json.load(file)
else:
    cumulative_predictions = {"0": 0, "1": 0, "2": 0}

# Load IoT device tracking JSON
if os.path.isfile(IOT_JSON_PATH):
    with open(IOT_JSON_PATH, 'r') as file:
        iot_devices = set(json.load(file))  # Convert list to set
else:
    iot_devices = set()

print("Loading models...")
iot_model = joblib.load('/opt/fortifi/models/device_identification_model.pkl')
intrusion_model = joblib.load('/opt/fortifi/models/packet_detection_model.pkl')
print("Models loaded successfully.")

def save_iot_devices():
    """Save unique IoT device IPs to JSON."""
    try:
        with open(IOT_JSON_PATH, 'w') as file:
            json.dump(list(iot_devices), file, indent=4)  # Convert set to list
    except Exception as e:
        print(f"❌ Failed to save IoT device list: {e}")
        
def load_tokens():
    if os.path.exists(TOKEN_STORAGE_PATH):
        with open(TOKEN_STORAGE_PATH, "r") as file:
            tokens = json.load(file)
            return tokens.get("jwt", "")
    return ""

def refresh_token():
    refresh_url = "http://localhost:5000/refresh_token"
    try:
        response = requests.get(refresh_url)
        if response.status_code == 200:
            new_token_data = response.json()
            return new_token_data.get("jwt", "")
    except Exception as e:
        print(f"Error refreshing token: {e}")
    return load_tokens()

def update_weekly_distribution():
    JWT_TOKEN = refresh_token()

    try:
        with open(CUMULATIVE_JSON_PATH, 'r') as file:
            cumulative_predictions = json.load(file)
    except Exception as e:
        print(f"❌ Failed to load updated cumulative predictions: {e}")
        cumulative_predictions = {"0": 0, "1": 0, "2": 0}

    data = {
        "benign": cumulative_predictions.get("0", 0),
        "port_scan": cumulative_predictions.get("1", 0),
        "ddos": cumulative_predictions.get("2", 0)
    }

    API_URL = "http://3.101.107.33:3000/UpdateWeeklyDistribution"
    headers = {"Authorization": f"Bearer {JWT_TOKEN}", "Content-Type": "application/json"}

    try:
        print(f"[→] Sending weekly distribution to {API_URL}")
        print(f"[→] Request payload: {json.dumps(data, indent=2)}")
        response = requests.post(API_URL, headers=headers, json=data)
        print(f"[←] Response Status: {response.status_code}")
        print(f"[←] Response Body: {response.text}")
        print(f"✅ Weekly Distribution Updated.")
    except Exception as e:
        print(f"❌ Failed to update weekly distribution: {e}")

def safe_stats(arr):
    if not arr:
        return {
            "min": 0, "max": 0, "mean": 0, "std": 0, "sum": 0, "count": 0
        }
    arr_np = np.array(arr)
    return {
        "min": float(np.min(arr_np)),
        "max": float(np.max(arr_np)),
        "mean": float(np.mean(arr_np)),
        "std": float(np.std(arr_np)),
        "sum": float(np.sum(arr_np)),
        "count": len(arr_np)
    }

def get_application_info(protocol_str):
    mapping = {
        "bittorrent": ("BitTorrent", "Download"),
        "dhcp": ("DHCP", "Network"),
        "bootp": ("DHCP", "Network"),
        "dhcpv6": ("DHCPV6", "Network"),
        "dns": ("DNS", "Network"),
        "dropbox": ("Dropbox", "Cloud"),
        "http": ("HTTP", "Web"),
        "icmp": ("ICMP", "Network"),
        "icmpv6": ("ICMPV6", "Network"),
        "igmp": ("IGMP", "Network"),
        "imaps": ("IMAPS", "Email"),
        "imo": ("IMO", "VoIP"),
        "ipsec": ("IPSec", "VPN"),
        "llmnr": ("LLMNR", "Network"),
        "mdns": ("MDNS", "Network"),
        "munin": ("Munin", "System"),
        "nat-pmp": ("NAT-PMP", "Network"),
        "ntp": ("NTP", "System"),
        "quic": ("QUIC", "Web"),
        "sip": ("SIP", "VoIP"),
        "smtps": ("SMTPS", "Email"),
        "ssdp": ("SSDP", "System"),
        "ssh": ("SSH", "RemoteAccess"),
        "stun": ("STUN", "Network"),
        "tls": ("TLS", "Media"),
        "ssl": ("TLS", "Media"),
        "ubntac2": ("UBNTAC2", "Network"),
        "unknown": ("Unknown", "Unspecified"),
        "viber": ("Viber", "VoIP"),
        "whois-das": ("Whois-DAS", "Network")
    }
    if protocol_str in mapping:
        return mapping[protocol_str]
    for key in mapping:
        if key in protocol_str:
            return mapping[key]
    return ("Unknown", "Unspecified")

def convert_to_builtin_type(obj):
    """Recursively convert numpy types in a structure to Python built-in types."""
    if isinstance(obj, dict):
        return {k: convert_to_builtin_type(v) for k, v in obj.items()}
    elif isinstance(obj, list):
        return [convert_to_builtin_type(i) for i in obj]
    elif isinstance(obj, (np.integer, np.int64)):
        return int(obj)
    elif isinstance(obj, (np.floating, np.float64)):
        return float(obj)
    else:
        return obj

def notify_new_device(name, ip_address, mac_address):
    JWT_TOKEN = refresh_token()
    headers = {
        "Authorization": f"Bearer {JWT_TOKEN}",
        "Content-Type": "application/json"
    }
    device_data = {
        "name": name,
        "ip_address": ip_address,
        "mac_address": mac_address
    }

    try:
        #response = requests.post("http://192.168.200.150:3000/AddDevice", headers=headers, json=device_data)
        response = requests.post("http://3.101.107.33:3000/AddDevice", headers=headers, json=device_data)
        print(f"[→] Adding new IoT device at /AddDevice")
        print(f"[→] Device Payload: {json.dumps(device_data, indent=2)}")
        print(f"[←] Response Status: {response.status_code}")
        print(f"[←] Response Body: {response.text}")

        if response.status_code == 200:
            print("✅ Device Successfully Added.")
        else:
            print(f"❌ Device Add failed with status {response.status_code}: {response.text}")
    except Exception as e:
        print(f"❌ Failed to send device data to API: {e}")

def notify_intrusion(event_data):
    JWT_TOKEN = refresh_token()
    headers = {"Authorization": f"Bearer {JWT_TOKEN}", "Content-Type": "application/json"}

    try:
        response = requests.post(API_URL, headers=headers, json=event_data)
        print(f"[→] Sending intrusion notification to {API_URL}")
        print(f"[→] Request payload: {json.dumps(event_data, indent=2)}")
        print(f"[←] Response Status: {response.status_code}")
        print(f"[←] Response Body: {response.text}")

        if response.status_code == 200:
            print("✅ Intrusion Event Successfully Sent.")
        else:
            print(f"❌ Notification failed with status {response.status_code}: {response.text}")
    except TypeError as te:
        print("❌ TypeError: Failed to serialize event_data. Here's the payload:")
        print(event_data)
        raise
    except Exception as e:
        print(f"❌ Failed to send notification: {e}")

def preprocess_and_predict(json_csv_path, conn_csv_path):
    from joblib import load

    try:
        with open(IOT_JSON_PATH, 'r') as f:
            known_iot_ips = set(json.load(f))
    except Exception:
        known_iot_ips = set()

    try:
        full_json_df = pd.read_csv(json_csv_path)
        if 'src_ip' not in full_json_df.columns or 'dst_ip' not in full_json_df.columns:
            print("❌ 'src_ip' or 'dst_ip' column missing from JSON CSV.")
            return [], pd.DataFrame(), []
    except Exception as e:
        print(f"❌ Failed to read JSON CSV: {e}")
        return [], pd.DataFrame(), []

    original_src_ips = full_json_df['src_ip'].astype(str).tolist()
    original_macs = full_json_df.get('src_mac', pd.Series(['00:00:00:00:00:00'] * len(full_json_df))).astype(str).tolist()

    cat_cols = ['src_ip', 'dst_ip', 'src_mac', 'dst_mac', 'application_name', 'application_category_name']
    le = LabelEncoder()
    for col in cat_cols:
        if col in full_json_df.columns:
            try:
                full_json_df[col] = le.fit_transform(full_json_df[col].astype(str))
            except:
                full_json_df[col] = 0
        else:
            full_json_df[col] = 0

    full_json_df = full_json_df.fillna(0)
    for col in full_json_df.columns:
        if full_json_df[col].dtype == 'object':
            try:
                full_json_df[col] = full_json_df[col].astype(float)
            except:
                full_json_df[col] = 0.0

    drop_cols = ['src_ip', 'dst_ip', 'src_mac', 'dst_mac']
    full_json_df = full_json_df.drop(columns=[c for c in drop_cols if c in full_json_df.columns], errors='ignore')

    # Print traffic before scaling
    print("=== IoT Traffic Before Scaling ===")
    print(full_json_df.head())

    try:
        scaler = StandardScaler()
        X_scaled = scaler.fit_transform(full_json_df)
        X = pd.DataFrame(X_scaled, columns=full_json_df.columns)
    except Exception as e:
        print(f"❌ Scaling failed: {e}")
        return [], pd.DataFrame(), []

    rows_to_predict = [i for i, ip in enumerate(original_src_ips) if ip not in known_iot_ips]
    if not rows_to_predict:
        print("No new IPs to classify.")
    else:
        X_pred = X.iloc[rows_to_predict]
        to_predict_ips = [original_src_ips[i] for i in rows_to_predict]
        to_predict_macs = [original_macs[i] for i in rows_to_predict]

        try:
            preds = iot_model.predict(X_pred)
            print("=== IoT Model Predictions ===")
            for ip, mac, pred in zip(to_predict_ips, to_predict_macs, preds):
                print(f"IP: {ip}, MAC: {mac}, Predicted Class: {pred}")
        except Exception as e:
            print(f"❌ Model prediction failed: {e}")
            return [], pd.DataFrame(), []

        for ip, mac, pred in zip(to_predict_ips, to_predict_macs, preds):
            if pred in IOT_LABEL_MAPPING:
                if ip not in known_iot_ips:
                    known_iot_ips.add(ip)
                    label = IOT_LABEL_MAPPING[pred]
                    notify_new_device(label, ip, mac)

        with open(IOT_JSON_PATH, 'w') as f:
            json.dump(sorted(known_iot_ips), f, indent=4)

    # === INTRUSION DETECTION ON CONN CSV ===
    try:
        conn_df = pd.read_csv(conn_csv_path)
    except Exception as e:
        print(f"❌ Failed to read conn CSV: {e}")
        return [], pd.DataFrame(), []

    required_columns = [
        'ts', 'uid', 'id.orig_h', 'id.orig_p', 'id.resp_h', 'id.resp_p',
        'proto', 'conn_state', 'missed_bytes', 'orig_pkts',
        'orig_ip_bytes', 'resp_pkts', 'resp_ip_bytes'
    ]

    if not all(col in conn_df.columns for col in required_columns):
        print("❌ Missing one or more required columns in conn CSV.")
        return [], pd.DataFrame(), []

    original_src_ips_conn = conn_df['id.orig_h'].astype(str).tolist()
    original_dst_ips_conn = conn_df['id.resp_h'].astype(str).tolist()
    original_uids_conn = conn_df['uid'].astype(str).tolist()

    rows_to_predict_indices = [
        i for i, (src, dst) in enumerate(zip(original_src_ips_conn, original_dst_ips_conn))
        if src in known_iot_ips or dst in known_iot_ips
    ]

    if not rows_to_predict_indices:
        return [], conn_df, []

    label_columns = ['uid', 'id.orig_h', 'id.resp_h', 'proto', 'conn_state']
    label_encoder = LabelEncoder()
    for col in label_columns:
        try:
            conn_df[col] = label_encoder.fit_transform(conn_df[col].astype(str))
        except Exception as e:
            print(f"⚠️ Failed to label encode column '{col}': {e}")
            conn_df[col] = 0

    conn_df = conn_df.fillna(0)
    for col in conn_df.columns:
        if conn_df[col].dtype == 'object':
            try:
                conn_df[col] = conn_df[col].astype(float)
            except Exception:
                conn_df[col] = 0.0

    # Print conn lines before scaling
    print("=== Conn Lines BEFORE Label Encoding & Scaling ===")
    print(conn_df.loc[rows_to_predict_indices, required_columns].head())

    try:
        conn_scaler = StandardScaler()
        scaled_conn = conn_scaler.fit_transform(conn_df)
    except Exception as e:
        print(f"❌ Failed to scale conn data: {e}")
        return [], conn_df, []

    to_predict = [scaled_conn[i] for i in rows_to_predict_indices]
    processed_df = conn_df.iloc[rows_to_predict_indices].reset_index(drop=True)
    processed_df['original_uid'] = [original_uids_conn[i] for i in rows_to_predict_indices]

    try:
        conn_predictions = intrusion_model.predict(to_predict)
        print("=== Intrusion Model Predictions ===")
        for i, pred in enumerate(conn_predictions):
            src = original_src_ips_conn[rows_to_predict_indices[i]]
            dst = original_dst_ips_conn[rows_to_predict_indices[i]]
            print(f"Packet {i}: {src} → {dst} | Predicted: {pred}")

        malicious_ips_path = os.path.join(CONFIG_DIR, "malicious_ips.txt")
        try:
            with open(malicious_ips_path, "r") as f:
                malicious_ips = set(line.strip() for line in f if line.strip())
        except Exception as e:
            print(f"❌ Failed to load malicious IP list: {e}")
            malicious_ips = set()

        malicious_flags = []
        EXFIL_THRESHOLD = 10 * 1024 * 1024  # 10 MB

        for idx in rows_to_predict_indices:
            src_ip = original_src_ips_conn[idx]
            dst_ip = original_dst_ips_conn[idx]
            is_malicious_ip = src_ip in malicious_ips or dst_ip in malicious_ips

            try:
                bytes_sent = int(conn_df.loc[idx, 'orig_ip_bytes'])
            except Exception:
                bytes_sent = 0

            is_internal_src = src_ip.startswith("192.168.") or src_ip.startswith("10.")
            is_external_dst = not (dst_ip.startswith("192.168.") or dst_ip.startswith("10."))
            is_exfil = is_internal_src and is_external_dst and bytes_sent > EXFIL_THRESHOLD

            if is_malicious_ip:
                malicious_flags.append(1)
            elif is_exfil:
                malicious_flags.append(2)
            else:
                malicious_flags.append(0)

        processed_df['malicious'] = malicious_flags

    except Exception as e:
        print(f"❌ Failed to load or run detection model: {e}")
        conn_predictions = []

    processed_df['id.orig_h'] = [original_src_ips_conn[i] for i in rows_to_predict_indices]
    processed_df['id.resp_h'] = [original_dst_ips_conn[i] for i in rows_to_predict_indices]

    return conn_predictions, processed_df, malicious_flags

   
def track_ip_events(predictions, original_packets, malicious_flags):
    prediction_counts = Counter()

    for idx, pred in enumerate(predictions):
        try:
            row = original_packets.iloc[idx]
            src_ip = str(row.get('id.orig_h', 'unknown')).strip()
            dst_ip = str(row.get('id.resp_h', 'unknown')).strip()
            uid = str(row.get('original_uid', f"packet_{idx}"))

            prediction_counts[str(pred)] += 1

            if pred in [1, 2]:
                event_data = {
                    "event": {
                        "id": uid,
                        "details": "Intrusion detected on your network.",
                        "ts": datetime.utcnow().strftime('%Y-%m-%d %H:%M:%S'),
                        "expires": (datetime.utcnow() + timedelta(hours=1)).strftime('%Y-%m-%d %H:%M:%S'),
                        "type": str(pred),
                        "src": src_ip,
                        "dst": dst_ip,
                        "confidence": 80
                    }
                }

                converted_event_data = convert_to_builtin_type(event_data)

                ip_dir = os.path.join(ASSET_TRACKING_DIR, src_ip)
                os.makedirs(ip_dir, exist_ok=True)

                json_filename = f"event_{datetime.utcnow().strftime('%Y%m%d_%H%M%S')}.json"
                json_path = os.path.join(ip_dir, json_filename)

                try:
                    with open(json_path, "w") as json_file:
                        json.dump(converted_event_data, json_file, indent=4)
                        json_file.flush()
                        os.fsync(json_file.fileno())
                except Exception as e:
                    print(f"❌ Failed to write event JSON: {e}")
                    continue

                notify_intrusion(converted_event_data)

        except Exception as e:
            print(f"❌ Error tracking event at index {idx}: {e}")

    for idx, flag in enumerate(malicious_flags):
        try:
            if flag == 0:
                continue

            row = original_packets.iloc[idx]
            src_ip = str(row.get('id.orig_h', 'unknown')).strip()
            dst_ip = str(row.get('id.resp_h', 'unknown')).strip()
            uid = str(row.get('original_uid', f"malicious_{idx}"))

            if flag == 1:
                details = "Traffic with known malicious IP detected."
                pred_val = 3
                filename_prefix = "malicious"
            elif flag == 2:
                details = "Potential data exfiltration detected: large outbound data transfer."
                pred_val = 4
                filename_prefix = "exfil"
            else:
                continue  # future-proofing

            event_data = {
                "event": {
                    "id": uid,
                    "details": details,
                    "ts": datetime.utcnow().strftime('%Y-%m-%d %H:%M:%S'),
                    "expires": (datetime.utcnow() + timedelta(hours=1)).strftime('%Y-%m-%d %H:%M:%S'),
                    "type": str(pred_val),
                    "src": src_ip,
                    "dst": dst_ip,
                    "confidence": 80
                }
            }

            converted_event_data = convert_to_builtin_type(event_data)

            ip_dir = os.path.join(ASSET_TRACKING_DIR, src_ip)
            os.makedirs(ip_dir, exist_ok=True)

            json_filename = f"{filename_prefix}_{datetime.utcnow().strftime('%Y%m%d_%H%M%S')}.json"
            json_path = os.path.join(ip_dir, json_filename)

            try:
                with open(json_path, "w") as json_file:
                    json.dump(converted_event_data, json_file, indent=4)
                    json_file.flush()
                    os.fsync(json_file.fileno())
            except Exception as e:
                print(f"❌ Failed to write alert JSON: {e}")
                continue

            notify_intrusion(converted_event_data)
            prediction_counts[str(pred_val)] += 1

        except Exception as e:
            print(f"❌ Error tracking flag event at index {idx}: {e}")

    # Save cumulative predictions
    try:
        if os.path.isfile(CUMULATIVE_JSON_PATH):
            with open(CUMULATIVE_JSON_PATH, 'r') as file:
                existing_predictions = json.load(file)
        else:
            existing_predictions = {"0": 0, "1": 0, "2": 0, "3": 0, "4": 0}
    except Exception:
        existing_predictions = {"0": 0, "1": 0, "2": 0, "3": 0, "4": 0}

    for key in ["0", "1", "2", "3", "4"]:
        existing_predictions[key] = existing_predictions.get(key, 0) + prediction_counts.get(key, 0)

    try:
        with open(CUMULATIVE_JSON_PATH, 'w') as file:
            json.dump(existing_predictions, file, indent=4)
    except Exception as e:
        print(f"❌ Failed to save cumulative predictions: {e}")

def convert_json_to_csv(session_dir):
    json_path = os.path.join(session_dir, "capture.json")
    output_path = os.path.join(session_dir, "json_converted.csv")

    try:
        with open(json_path, "r") as f:
            packets = json.load(f)
    except Exception as e:
        print(f"❌ Failed to load JSON from {json_path}: {e}")
        return output_path

    sessions = defaultdict(lambda: {
        "src_mac": None,
        "dst_mac": None,
        "src_ip": None,
        "dst_ip": None,
        "src_port": None,
        "dst_port": None,
        "src2dst_bytes": 0,
        "dst2src_bytes": 0,
        "src2dst_sizes": [],
        "dst2src_sizes": [],
        "ack_packets": 0,
        "rst_packets": 0,
        "total_packets": 0,
        "application_name": None,
        "application_category_name": None
    })

    for pkt in packets:
        layers = pkt.get("_source", {}).get("layers", {})
        frame = layers.get("frame", {})
        eth = layers.get("eth", {})
        ip = layers.get("ip", {})
        udp = layers.get("udp", {})
        tcp = layers.get("tcp", {})

        if not ip:
            continue

        src_ip = ip.get("ip.src")
        dst_ip = ip.get("ip.dst")
        src_mac = eth.get("eth.src")
        dst_mac = eth.get("eth.dst")
        src_port = udp.get("udp.srcport") or tcp.get("tcp.srcport") or "0"
        dst_port = udp.get("udp.dstport") or tcp.get("tcp.dstport") or "0"
        length = int(frame.get("frame.len", 0))
        protocols = frame.get("frame.protocols", "")
        app_name, app_cat = get_application_info(protocols)

        key = (src_ip, dst_ip, src_port, dst_port)
        rev_key = (dst_ip, src_ip, dst_port, src_port)

        if key not in sessions and rev_key in sessions:
            key = rev_key

        session = sessions[key]
        if session["src_ip"] is None:
            session["src_ip"] = src_ip
            session["dst_ip"] = dst_ip
            session["src_mac"] = src_mac
            session["dst_mac"] = dst_mac
            session["src_port"] = int(src_port)
            session["dst_port"] = int(dst_port)
            session["application_name"] = app_name
            session["application_category_name"] = app_cat

        direction = "src2dst" if (src_ip == session["src_ip"] and src_port == str(session["src_port"])) else "dst2src"
        if direction == "src2dst":
            session["src2dst_bytes"] += length
            session["src2dst_sizes"].append(length)
        else:
            session["dst2src_bytes"] += length
            session["dst2src_sizes"].append(length)

        if tcp:
            flags_tree = tcp.get("tcp.flags_tree", {})
            if flags_tree.get("tcp.flags.ack") == "1":
                session["ack_packets"] += 1
            if flags_tree.get("tcp.flags.reset") == "1" and direction == "dst2src":
                session["rst_packets"] += 1

        session["total_packets"] += 1

    results = []
    for sess in sessions.values():
        bidirectional_bytes = sess["src2dst_bytes"] + sess["dst2src_bytes"]
        src2dst_stats = safe_stats(sess["src2dst_sizes"])
        dst2src_stats = safe_stats(sess["dst2src_sizes"])
        total_stats = safe_stats(sess["src2dst_sizes"] + sess["dst2src_sizes"])

        results.append({
            "src_ip": sess["src_ip"],
            "src_mac": sess["src_mac"],
            "dst_ip": sess["dst_ip"],
            "dst_mac": sess["dst_mac"],
            "src_port": sess["src_port"],
            "dst_port": sess["dst_port"],
            "bidirectional_bytes": bidirectional_bytes,
            "src2dst_bytes": sess["src2dst_bytes"],
            "dst2src_bytes": sess["dst2src_bytes"],
            "bidirectional_min_ps": total_stats["min"],
            "bidirectional_mean_ps": total_stats["mean"],
            "bidirectional_stddev_ps": total_stats["std"],
            "bidirectional_max_ps": total_stats["max"],
            "src2dst_min_ps": src2dst_stats["min"],
            "src2dst_mean_ps": src2dst_stats["mean"],
            "src2dst_max_ps": src2dst_stats["max"],
            "dst2src_min_ps": dst2src_stats["min"],
            "dst2src_mean_ps": dst2src_stats["mean"],
            "dst2src_max_ps": dst2src_stats["max"],
            "bidirectional_packets": sess["total_packets"],
            "bidirectional_ack_packets": sess["ack_packets"],
            "dst2src_rst_packets": sess["rst_packets"],
            "application_name": sess["application_name"],
            "application_category_name": sess["application_category_name"]
        })

    try:
        df = pd.DataFrame(results)
        df.to_csv(output_path, index=False)
    except Exception as e:
        print(f"❌ Failed to write CSV: {e}")

    return output_path

def convert_connlog_to_csv(session_dir):
    """Convert Zeek conn.log to CSV with selected headers."""
    connlog_path = os.path.join(session_dir, "conn.log")
    output_csv_path = os.path.join(session_dir, "connlog_converted.csv")

    if not os.path.exists(connlog_path):
        print(f"❌ conn.log not found in {session_dir}")
        return output_csv_path

    headers = []
    records = []
    with open(connlog_path, 'r') as f:
        for line in f:
            if line.startswith("#fields"):
                headers = line.strip().split('	')[1:]
            elif line.startswith("#"):
                continue
            else:
                fields = line.strip().split('	')
                if len(fields) == len(headers):
                    records.append(dict(zip(headers, fields)))

    required_columns = [
        "ts", "uid", "id.orig_h", "id.orig_p", "id.resp_h", "id.resp_p",
        "proto", "conn_state", "missed_bytes", "orig_pkts", "orig_ip_bytes",
        "resp_pkts", "resp_ip_bytes"
    ]

    df = pd.DataFrame(records)
    for col in required_columns:
        if col not in df.columns:
            df[col] = ""

    df = df[required_columns]
    df.to_csv(output_csv_path, index=False)
    return output_csv_path

def monitor_and_process():
    while True:
        for base_dir, processed_dir in [(TSHARK_DIR, PROCESSED_TSHARK_DIR), (PINEAPPLE_DIR, PROCESSED_PINEAPPLE_DIR)]:
            capture_dirs = sorted([
                os.path.join(base_dir, d) for d in os.listdir(base_dir)
                if os.path.isdir(os.path.join(base_dir, d)) and d.startswith("capture_")
            ], key=os.path.getmtime)

            if len(capture_dirs) < 2:
                continue  # Not enough to process

            for session_dir in capture_dirs[:-1]:  # Skip the most recent (active) directory
                json_csv = convert_json_to_csv(session_dir)
                conn_csv = convert_connlog_to_csv(session_dir)

                if os.path.exists(json_csv) and os.path.exists(conn_csv):
                    conn_predictions, processed_df, malicious_flags = preprocess_and_predict(json_csv, conn_csv)


                    if len(conn_predictions) > 0:
                        track_ip_events(conn_predictions, processed_df, malicious_flags)

                    # ✅ Move processed session directory immediately after processing
                    try:
                        shutil.move(session_dir, processed_dir)
                    except Exception as e:
                        print(f"❌ Failed to move {session_dir} to {processed_dir}: {e}")

        try:
            update_weekly_distribution()
        except Exception as e:
            print(f"❌ Failed to update weekly distribution: {e}")

        time.sleep(30)

if __name__ == "__main__":
    monitor_and_process()
