import React from 'react';
import './board.css';

const Board = ({ boardState, onMove }) => {
    const renderSquare = (row, col) => {
        const isLightSquare = (row + col) % 2 === 0;
        const squareClass = isLightSquare ? 'light-square' : 'dark-square';
        const piece = boardState[row][col]; // Get the piece at this position

        return (
            <div
                key={`${row}-${col}`}
                className={`square ${squareClass}`}
                onClick={() => onMove(row, col)}
            >
                {piece && <div className={`piece ${piece}`}></div>} {/* Render piece if present */}
            </div>
        );
    };

    return (
        <div className="board">
            {boardState.map((row, rowIndex) => (
                <div key={rowIndex} className="row">
                    {row.map((_, colIndex) => renderSquare(rowIndex, colIndex))}
                </div>
            ))}
        </div>
    );
};

export default Board;