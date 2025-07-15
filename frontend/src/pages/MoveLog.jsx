import React from 'react';
import './MoveLog.css';

const MESSAGE_LIMIT = 50; // Limit the number of messages for protection

const MoveLog = ({ moveLog }) => {
  // Only show the last MESSAGE_LIMIT moves
  const limitedLog = moveLog.slice(-MESSAGE_LIMIT);
  return (
    <div className="move-log">
      <h3>Move Log</h3>
      <div className="move-log-list">
        {limitedLog.length === 0 ? (
          <div className="move-log-empty">No moves yet.</div>
        ) : (
          limitedLog.map((move, idx) => (
            <div key={idx} className="move-log-entry">
              {move}
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default MoveLog;
