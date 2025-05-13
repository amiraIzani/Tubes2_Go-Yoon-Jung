import React from "react";

const SearchStats = ({ nodes = 0, searchTime = 0 }) => {
  return (
    <div className="search-stats">
      <div className="search-stats__metric">
        <div className="search-stats__label">Number of Nodes</div>
        <div className="search-stats__value">{nodes} nodes</div>
      </div>
      <div className="search-stats__metric">
        <div className="search-stats__label">Search Time</div>
        <div className="search-stats__value">{searchTime} ms</div>
      </div>
    </div>
  );
};

export default SearchStats;
