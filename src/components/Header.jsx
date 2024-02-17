import Search from "./Search";

function Header({ query }) {
  return (
    <div className="header-container">
      <Search styleType="header" />
      <h1 className="header-results">Search Result for: {query}</h1>
    </div>
  );
}

export default Header;
