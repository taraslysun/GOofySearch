import { useParams } from "react-router-dom";
import Header from "../components/Header";

function SearchResult() {
  const { query } = useParams();
  return (
    <div>
      <Header query={query} />
    </div>
  );
}

export default SearchResult;
