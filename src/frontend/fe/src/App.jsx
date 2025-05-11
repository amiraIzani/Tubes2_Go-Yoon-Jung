import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import SearchInterface from "./pages/SearchInterface";
import ResultsPage from "./pages/ResultsPage";
import "./App.css";

function App() {

  return (
    <Router>
      <Routes>
        <Route path="/" element={<SearchInterface />} />
        <Route path="/results" element={<ResultsPage nodes={42} searchTime={150} />} />
      </Routes>
    </Router>
  );

}

export default App;
