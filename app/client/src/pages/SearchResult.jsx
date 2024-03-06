import { useParams } from "react-router-dom";
import Header from "../components/Header";
import Results from "../components/Results";

function SearchResult() {
  const { query } = useParams();
  return (
    <div className="searchRes">
      <Header query={query} />

      <Results query={query} />
    </div>
  );
}

export default SearchResult;
