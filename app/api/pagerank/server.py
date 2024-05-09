from flask import Flask, request, jsonify
from flask_cors import CORS
import requests
import numpy as np
from pageRank import pagerank
import re

app = Flask(__name__)
CORS(app)

def parse_sites(site_link_list: list[str]):
    site_links = {}
    for site in site_link_list:
        site_links[site] = []
    for site in site_link_list:
        response = requests.get(site)
        if response.status_code == 200:
            links = re.findall(r'href=[\'"]?([^\'" >]+)', response.text)
            links = [link for link in links if link.startswith('http')]
            site_links[site] = links
            print()
    return site_links

def build_adj_matrix(site_links: dict):
    site_names = list(site_links.keys())
    N = len(site_names)
    adj_matrix = np.zeros((N, N), dtype=int)
    for i, site in enumerate(site_names):
        for link in site_links[site]:
            if link in site_names:
                j = site_names.index(link)
                adj_matrix[i, j] = 1

    print(adj_matrix)
    return adj_matrix, site_names

    
    


@app.route('/pagerank', methods=['POST'])
def page_rank():
    data = request.json
    site_links = data['site_links']
    
    site_links = parse_sites(site_links)
    adj_matrix, site_names = build_adj_matrix(site_links)
    ordered_sites = pagerank(adj_matrix, site_names)
    return jsonify(ordered_sites)



# Example usage:
# curl -X POST http://localhost:5000/pagerank -H "Content-Type: application/json" -d '{"site_links": ["https://www.google.com", "https://www.yahoo.com"]}'

if __name__ == '__main__':
    app.run(debug=True)
