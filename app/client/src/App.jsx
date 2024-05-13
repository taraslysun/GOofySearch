import { Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import SearchResult from "./pages/SearchResult";
import CrawlerSystem from "./pages/CrawlerSystem";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/search/:query" element={<SearchResult />} />
      <Route path="/crawler_system" element={<CrawlerSystem />}/>
    </Routes>
  );
}

export default App;
