import { useState, useEffect } from "react";

function Results({ query }) {
  const [data, setData] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      const requestOptions = {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ query: query }),
      };
      const response = await fetch(
        "http://127.0.0.1:3000/api/search",
        requestOptions
      );
      const result = await response.json();
      let hits = result.hits.hits;

      hits = hits.filter((hit) => hit["_source"]["title"] != "");

      setData(hits);
    };

    fetchData();
  }, [query]);

  for (let i = 0; i < data.length; i++) {
    console.log(data[i]["_id"]);
    console.log(data[i]["_source"]);
    // console.log(data[i]);
  }

  return (
    <div className="results">
      {data.map((hit) => (
        <div key={hit["_id"]} className="results-container">
          <h3 className="result-title">{hit["_source"]["title"]}</h3>
          <a
            key={hit["_id"]}
            href={hit["_source"]["link"]}
            className="result-item"
          >
            {hit["_source"]["link"]}
          </a>
        </div>
      ))}
    </div>
  );
}

export default Results;