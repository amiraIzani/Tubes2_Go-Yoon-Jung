import React from "react";
import SearchStats from "../components/SearchStats";
import "../styles/ResultsPage.css";
import { useNavigate, useLocation } from "react-router-dom";
import TreeVisualizer from "../components/TreeVisualizer";

const ResultsPage = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const { calculatedResult } = location.state || {};

  const nodes = calculatedResult?.nodes_visited ?? 0;
  const searchTime = calculatedResult?.time_us ?? 0;
  const recipeTree = calculatedResult?.recipes ?? [];

  const handleResetClick = () => {
    navigate("/");
  };

  return (
    <div className="results-page">
      <div className="results-page__content">
        <TreeVisualizer tree={recipeTree} />
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
