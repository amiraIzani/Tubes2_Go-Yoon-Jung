import React, { useState } from "react";
import elements from "../assets/recipes.json";
import images from '../assets/images.json';
import search from '../assets/Search.png';

const SearchBar = ({ result, setResult }) => {
  const [searchText, setSearchText] = useState("");


  const handleSearch = () => {
    const found = elements.find(
      (el) => el.name.toLowerCase() === searchText.toLowerCase()
    );
    setResult(found || null);
  };

  const getImageLink = (elementName) => {
    const imageObj = images.find(
      (img) => img.name.toLowerCase() === elementName.toLowerCase()
    );

    return imageObj?.imageUrl;
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
          <img
            src={getImageLink(result.name)}
            alt={result.name}
            className="search-bar__element-icon"
          />
          <div className="search-bar__element-text">{result.name}</div>

        </div>
      )}
    </div>
  );
};

export default SearchBar;
