import { useParams } from "react-router-dom";
import Header from "../components/Header";
import Results from "../components/Results";
import { useState } from "react";

function SearchResult() {
  const { query } = useParams();
  const [sortByPagerank, setSortByPagerank] = useState(false);
  return (
    <div className="searchRes">
      <Header query={query} />

      <div className="sort">
        <button onClick={() => {
              if (sortByPagerank) {
                setSortByPagerank(false);
              }
              else {
                setSortByPagerank(true);  
              }
            }}
            style={
              // sortByPagerank
              //   ? { backgroundColor: "blue", color: "white" }
              //   : { backgroundColor: "white", color: "black" }
              {
                backgroundColor: "white",
                color:"black",
                borderRadius: "25px",
                fontWeight: "bold",
                padding: "5px",
                margin: "10px",
                border: sortByPagerank ? "4px solid green" : "4px solid gray"
              }
              
            }>PageRank sorted</button>
      </div>


      <Results query={query} sortByPagerank={sortByPagerank} />
    </div>
  );
}

export default SearchResult;
