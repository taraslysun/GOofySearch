from flask import Flask, request, jsonify
from flask_cors import CORS
from elasticsearch import Elasticsearch
import subprocess
import shlex
from pagerank.pageRank import pagerank
import numpy as np
import requests
from bs4 import BeautifulSoup
import time
import re


app = Flask(__name__)
CORS(app)
CLOUD_ID = "mini-google:ZXVyb3BlLXdlc3Q5LmdjcC5lbGFzdGljLWNsb3VkLmNvbTo0NDMkZWRhMWY0MTkyZmJiNGM3YjhiNDQ2ODk4NjBiNGMyNTckOTUzOThlMjVmNjdmNDA4MzhiYzJhOTE4ODAyZDZjYmI="
API_KEY = "NFBZcGFJOEJ1WDg3RXdUSUlaX2o6M0hsZkdsWGlSeEtZc1M0NGpqUXkzZw=="


es = Elasticsearch(cloud_id=CLOUD_ID, api_key=API_KEY, node_class='requests')


@app.route('/api/api_credentials', methods=['POST', 'GET'])
def api_credentials():
    data = request.json
    CLOUD_ID = data['cloud_id']
    API_KEY = data['api_key']
    global es
    es = Elasticsearch(cloud_id=CLOUD_ID, api_key=API_KEY, node_class='requests')
    print(CLOUD_ID, API_KEY)
    return jsonify({"message": "API credentials set"})


def build_adjacency_matrix(urls):
    n = len(urls)
    adj_matrix = np.zeros((n, n))
    url_to_index = {url: idx for idx, url in enumerate(urls)}

    for idx, url in enumerate(urls):
        # try:
        #     response = requests.get(url)
        #     soup = BeautifulSoup(response.content, 'html.parser')
        #     links = soup.find_all('a', href=True)
        #     for link in links:
        #         href = link['href']
        #         if href in url_to_index:
        #             adj_matrix[idx][url_to_index[href]] += 1
        # except requests.RequestException:
        #     continue
        try:
            response = requests.get(url)
            links = re.findall(r'href=[\'"]?([^\'" >]+)', response.text)
            for link in links:
                if link in url_to_index:
                    adj_matrix[idx][url_to_index[link]] += 1
        except requests.RequestException:
            continue

    return adj_matrix

@app.route('/api/search_pagerank', methods=['POST', 'GET'])
def search_and_rank():
    query = request.json.get('query')
    res = es.search(index="final_data", body={
        "query": {
            "bool": {
                "should": [
                    {"match": {"text": {"query": query, "boost": 1}}},
                    {"match": {"title": {"query": query, "boost": 3}}},
                ]
            },
        },
    })
    start = time.time()    
    urls = [hit['_source']['link'] for hit in res['hits']['hits']]
    adj_matrix = build_adjacency_matrix(urls)
    print("adj_matrix built in", time.time() - start)
    site_names = urls
    site_categories = [0] * len(urls)
    ordered_sites = pagerank(adj_matrix, site_categories, site_names)
    print(ordered_sites)
    print("ordered in", time.time() - start)
    ranked_results = {site: score for site, score, _ in ordered_sites}
    print("ranked in", time.time() - start)
    sorted_hits = sorted(res['hits']['hits'], key=lambda x: ranked_results[x['_source']['link']], reverse=True)
    print('sorted in', time.time() - start)
    return jsonify(sorted_hits)


@app.route('/api/search', methods=['POST', 'GET'])
def search():
    query = request.json.get('query')
    print(query)
    res = es.search(index="final_data", body={
        "query": {
            "bool": {
                "should": [
                    {"match": {"text": {"query": query, "boost": 1}}},
                    {"match": {"title": {"query": query, "boost": 3}}},
                ]
            },
        },
    })
    res1 = es.search(index="dmytro_table", body={
        "query": {
            "bool": {
                "should": [
                    {"match": {"text": {"query": query, "boost": 1}}},
                    {"match": {"title": {"query": query, "boost": 3}}},
                ]
            },
        },
    })
    res2 = es.search(index="sviat", body={
        "query": {
            "bool": {
                "should": [
                    {"match": {"text": {"query": query, "boost": 1}}},
                    {"match": {"title": {"query": query, "boost": 3}}},
                ]
            },
        },
    })
    small = list(res['hits']['hits'])
    small.extend(list(res1['hits']['hits']))
    small.extend(list(res2['hits']['hits']))
    return jsonify(small)

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
        tm_commands = f"cd {path_to_task_manager} && go run . &"
        cr_commands = f"cd {path_to_crawler} && go run main.go master {host_ip} {worker_num}"
        full_command = f"{tm_commands} {cr_commands}"
    else:
        full_command = f"cd {path_to_crawler} && go run main.go worker {host_ip}"

    ssh_command = f"sshpass -p {shlex.quote(password)} ssh -o StrictHostKeyChecking=no \
{shlex.quote(username)}@{shlex.quote(ip)} '{full_command}'"

    process = subprocess.Popen(ssh_command, shell=True)
    return jsonify({"message": "SSH command executed", "pid": process.pid})

if __name__ == '__main__':
    app.run(debug=True, port=3000, host = "0.0.0.0")
