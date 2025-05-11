import React, { useState } from "react";
import elements from "../assets/recipes.json";
import search from '../assets/Search.png';

const SearchBar = ({ result, setResult }) => {
  const [searchText, setSearchText] = useState("");


  const handleSearch = () => {
    const found = elements.find(
      (el) => el.name.toLowerCase() === searchText.toLowerCase()
    );
    setResult(found || null);
  };

  return (
    <div className="search-bar__wrapper">
      <div className="search-bar">
        <input
          type="text"
          className="search-bar__input"
          placeholder="Search Element..."
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
        />
        <div className="search-bar__button" onClick={handleSearch}>
          <img
            src={search}
            alt="Search"
            className="search-bar__icon"
          />
        </div>
      </div>


      {result && (
        <div className="search-bar__element">
          <div className="search-bar__element-icon"></div>
          <div className="search-bar__element-text">{result.name}</div>
          {result.recipes.map((r, index) => (
            <div key={index} className="search-bar__recipe" >
              {r.elements.join(" + ")}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default SearchBar;
