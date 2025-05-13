const ControlPanel = ({
  algorithm,
  setAlgorithm,
  recipeMode,
  setRecipeMode,
  maxParam,
  setMaxParam
}) => {
  const handleAlgorithmChange = (e) => {
    setAlgorithm(e.target.value);
  };

  const handleRecipeModeChange = (e) => {
    setRecipeMode(e.target.value);
    console.log('Recipe mode selected:', e.target.value);
  };

  const handleMaxParamChange = (e) => {
    setMaxParam(e.target.value);
  };

  return (
    <div className="control-panel">
      <div className="control-panel__settings">
        <div className="control-panel__section">
          <div className="control-panel__label">Choose Algorithm</div>
          <div className="control-panel__value">
            <select id="algorithm-selector" className="selector" value={algorithm} onChange={handleAlgorithmChange}>
              <option value="BFS">BFS</option>
              <option value="DFS">DFS</option>
            </select>
          </div>
        </div>
        <div className="control-panel__section">
          <div className="control-panel__label">How Many Recipes</div>
          <div className="control-panel__value">
            <select id="Recipes-selector" class="selector" value={recipeMode} onChange={handleRecipeModeChange}>
              <option value="one">One Recipe</option>
              <option value="multiple">Multiple Recipes</option>
            </select>
          </div>
        </div>
        <div className="control-panel__section">
          <div className="control-panel__label">Input Max Parameter</div>
          <div className="control-panel__input">
            <input
              type="number"
              value={maxParam}
              onChange={handleMaxParamChange}
              placeholder="Enter Max Param"
              min="1"
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default ControlPanel;

