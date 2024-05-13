"Main server module"

from elasticsearch import Elasticsearch
from flask_cors import CORS
from flask import Flask, request, jsonify
import subprocess
import shlex

CLOUD_ID = ""
API_KEY = ""

es = Elasticsearch(cloud_id=CLOUD_ID, api_key=API_KEY, node_class='requests')

app = Flask(__name__)
CORS(app)



@app.route('/api/search', methods=['POST', 'GET'])
def search():
    "Seach route"
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
    hits = res['hits']['hits']
    return jsonify(hits)

app = Flask(__name__)
CORS(app)

@app.route('/api/execute_ssh', methods=['POST'])
def execute_ssh():
    data = request.json
    ip = data['ip']
    username = data['username']
    password = data['password']

    path_to_crawler = data['path_to_crawler']
    path_to_task_manager = data['path_to_task_manager']
    is_host = data['is_host']

    id = data['id']
    host_ip = data['host_ip']
    worker_num = data['worker_num']

    print(ip, username, password, path_to_crawler, is_host, path_to_task_manager)

    cd_crawler = f"cd {path_to_crawler}"
    if is_host:
        cd_task_manager = f"cd {path_to_task_manager}"
        run_cr = f"go run main.go worker {host_ip} {id}"
    else:    
        run_cr = f"go run main.go host {host_ip} {worker_num}"
    command = f"sshpass -p {password} ssh -o StrictHostKeyChecking=no {username}@{ip} '{cd_crawler}; {run_cr}'"


if __name__ == '__main__':
    app.run(debug=True, port=3000)
