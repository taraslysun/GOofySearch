"Main server module"

from flask import Flask, request, jsonify
from elasticsearch import Elasticsearch
from flask_cors import CORS

# Elasticsearch Cloud credentials
CLOUD_ID = "mini-google:ZXVyb3BlLXdlc3QzLmdjcC5jbG91ZC5lcy5pbzo0NDMkZWRmNDdlMTMzNWQzNGEyMGFkMWFiMDg2Mjc5ODZkNWEkYjRkZTRjNjZlNDkyNGI1NDhhMjNkNWYyNTE5ZTNhZDk="
API_KEY = "VVpJRnlZNEJfRzZqUW1QVnVESFI6ek44cEdIVFNTMVdGNldTVVhRY0V2Zw=="

# Initialize Elasticsearch client
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

if __name__ == '__main__':
    app.run(debug=True, port=3000)
