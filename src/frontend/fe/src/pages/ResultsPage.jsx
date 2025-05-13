import React from "react";
import SearchStats from "../components/SearchStats";
import "../styles/ResultsPage.css";
import { useNavigate } from "react-router-dom";
import TreeVisualizer from "../components/TreeVisualizer";

const ResultsPage = ({ nodes = 0, searchTime = 0, onTryMore }) => {
  const navigate = useNavigate();

  const handleResetClick = () => {
    navigate("/");
  };

  return (
    <div className="results-page">
      <div className="results-page__content">
        <TreeVisualizer />
      </div>
      <div className="results-page__sidebar">
        <SearchStats nodes={nodes} searchTime={searchTime} />
        <div className="try-more-button" onClick={handleResetClick}>
          TRY MORE RECIPE
        </div>
      </div>
    </div>
  );
};

export default ResultsPage;
