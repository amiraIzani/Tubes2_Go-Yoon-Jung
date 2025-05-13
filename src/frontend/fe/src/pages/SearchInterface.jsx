import React from "react";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import SearchBar from "../components/SearchBar";
import ControlPanel from "../components/ControlPanel";
import "../styles/SearchInterface.css";

const SearchInterface = () => {
  // result, //element yg dipilih bentuk string 
  // (contoh: {
  //"name": "Coconut", "recipes": [{ "elements": ["Palm", "Fruit"] },
  const [result, setResult] = useState(null);
  const [algorithm, setAlgorithm] = useState('BFS');
  const [recipeMode, setRecipeMode] = useState('one');
  const [maxParam, setMaxParam] = useState('');

  const navigate = useNavigate();

  const handleResultClick = () => {
    if (!algorithm || !recipeMode || !maxParam || maxParam < 1 || !result) {
      return;
    }

    const calculatedResult = {
      ///dummy
      nodes: 1,
      searchTime: 90,
    };



    navigate("/results", {
      // state: {
      //   calculatedResult
      // },
    });



  };


  return (
    <div className="search-interface">
      <div className="search-interface__container">
        <div className="search-interface__content">
          <div className="search-interface__left-column">
            <SearchBar
              result={result}
              setResult={setResult}
            />
          </div>
          <div className="search-interface__right-column">
            <ControlPanel
              algorithm={algorithm}
              setAlgorithm={setAlgorithm}
              recipeMode={recipeMode}
              setRecipeMode={setRecipeMode}
              maxParam={maxParam}
              setMaxParam={setMaxParam}
            />
            <div className="control-panel__result" onClick={handleResultClick}>
              RESULT
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SearchInterface; 
