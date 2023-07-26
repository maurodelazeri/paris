from locust import HttpUser, TaskSet, task, between
import json
import random
import sys

# JSON-RPC request validation
def is_json_rpc(message):
    try:
        payload = json.loads(message)
    except ValueError:
        return False

    if not isinstance(payload, dict):
        return False

    if "jsonrpc" not in payload or payload["jsonrpc"] != "2.0":
        return False

    if "method" not in payload or not isinstance(payload["method"], str):
        return False

    if "params" not in payload or not isinstance(payload["params"], list):
        return False

    if "id" not in payload or not isinstance(payload["id"], (int, str)):
        return False

    return True

# Locust code
def run_request(locust, data, name):
    data["id"] = random.randint(0, sys.maxsize)
    with locust.client.post("", name=name, json=data, catch_response=True) as response:
        try:
            response_data = json.loads(response.content)
        except json.decoder.JSONDecodeError:
            response.failure(f"Invalid JSON received. Response code: {response.status_code}, content: {response.text}")
        else:
            if response_data.get("error", False):
                response.failure(f"Payload error: {response_data.get('error')}. Sent payload: {data}")
            # elif response_data.get("id") != data["id"]:
            #     response.failure(f"Mismatched IDs. Sent ID: {data['id']}, Received ID: {response_data.get('id')}")
            else:
                response.success()
        return response

class RunFullTest(TaskSet):
    @task
    def send_rpc_request(self):
        # This will be your JSON-RPC payload
        payload_json = {"method": "eth_getBlockByNumber", "params": ["latest", True], "id": 1, "jsonrpc": "2.0"}
        run_request(self, payload_json, name=payload_json["method"])

class WebsiteUser(HttpUser):
    tasks = [RunFullTest]
    wait_time = between(1, 5)
