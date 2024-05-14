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
      console.log(result)
      let hits = result;

      hits = hits.filter((hit) => hit["_source"]["title"] != "");

      setData(hits);
    };

    fetchData();
  }, [query]);

  for (let i = 0; i < data.length; i++) {
    console.log(data[i]["_id"]);
    console.log(data[i]["_source"]["title"]);
  }

  return (
    <div className="results">
      {data.map((hit) => (
        <a
          key={hit["_id"]}
          href={hit["_source"]["url"]}
          className="result-item"
        >
          {hit["_source"]["title"]}
        </a>
      ))}
    </div>
  );
}

export default Results;
