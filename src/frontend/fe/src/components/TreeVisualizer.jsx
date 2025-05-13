import React from 'react';
import Tree from 'react-d3-tree';

const containerStyles = { width: '100%', height: '100vh' };

export default function TreeVisualizer({ tree }) {
  if (!tree || (Array.isArray(tree) && tree.length === 0)) {
    return <div style={{ padding: '1rem' }}>No recipe tree available.</div>;
  }

  const treeArray = Array.isArray(tree) ? tree : [tree];

  return (
    <div style={containerStyles}>
      {treeArray.map((t, idx) => (
        <Tree
          key={idx}
          data={t}
          orientation="vertical"
          translate={{ x: window.innerWidth / 2, y: 50 + idx * 300 }}
        />
      ))}
    </div>
  );
}
