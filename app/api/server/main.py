from flask import Flask, request, jsonify
from flask_cors import CORS
from elasticsearch import Elasticsearch
import subprocess
import shlex


app = Flask(__name__)
CORS(app)


@app.route('/api/api_credentials', methods=['POST'])
def api_credentials():
    data = request.json
    CLOUD_ID = data['cloud_id']
    API_KEY = data['api_key']
    global es
    es = Elasticsearch(cloud_id=CLOUD_ID, api_key=API_KEY, node_class='requests')
    print(CLOUD_ID, API_KEY)
    return jsonify({"message": "API credentials set"})


@app.route('/api/search', methods=['POST', 'GET'])
def search():
    query = request.json.get('query')
    res = es.search(index="test", body={
        "query": {
            "bool": {
                "should": [
                    {"match": {"text": {"query": query, "boost": 1}}},
                    {"match": {"title": {"query": query, "boost": 3}}},
                ]
            },
        },
    })
    return jsonify(res['hits']['hits'])

@app.route('/api/execute_ssh', methods=['POST'])
def execute_ssh():
    data = request.json
    ip = data['ip']
    username = data['username']
    password = data['password']
    path_to_crawler = data['path_to_crawler']
    path_to_task_manager = data['path_to_task_manager']
    is_host = data['is_host']
    host_ip = data['host_ip']
    worker_num = data['worker_num']

    if is_host:
        tm_commands = f"cd {path_to_task_manager} && nohup go run . &"
        cr_commands = f"cd {path_to_crawler} && go run main.go master {host_ip} {worker_num}"
        full_command = f"{tm_commands} {cr_commands}"
    else:
        full_command = f"cd {path_to_crawler} && go run main.go worker {host_ip}"

    ssh_command = f"sshpass -p {shlex.quote(password)} ssh -o StrictHostKeyChecking=no \
{shlex.quote(username)}@{shlex.quote(ip)} '{full_command}'"

    process = subprocess.Popen(ssh_command, shell=True)
    return jsonify({"message": "SSH command executed", "pid": process.pid})

if __name__ == '__main__':
    app.run(debug=True, port=3000)
