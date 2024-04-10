import logo from "../assets/goofle_logo.png";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

import "./styles.css";

function Search({ styleType }) {
  const [searchQuery, setSearchQuery] = useState("");
  const navigate = useNavigate();

  const handleSubmit = function (event) {
    event.preventDefault();

    if (searchQuery) {
      navigate(`/search/${searchQuery}`);
    }
    setSearchQuery("");
  };

  const handleChange = function (event) {
    event.preventDefault();
    setSearchQuery(event.target.value);
  };

  return (
    <div className={styleType == "main" ? "main-search" : "header-search"}>
      <img
        className={styleType == "main" ? "main-logo" : "header-logo"}
        src={logo}
        alt="Logo"
      />

      <form
        onSubmit={handleSubmit}
        className={
          styleType == "main"
            ? "main-search-container"
            : "header-search-container"
        }
      >
        <input
          type="text"
          placeholder="Search..."
          value={searchQuery}
          onChange={handleChange}
        />
      </form>
    </div>
  );
}

export default Search;
