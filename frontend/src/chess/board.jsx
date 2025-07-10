import React, { useRef } from 'react';
import './board.css';

export function InitializeBoard() {
    const initialBoard = Array(8).fill(null).map(() => Array(8).fill(null));
    // Initialize pieces on the board. Black is on top and White is on bottom.
    for (let i = 0; i < 8; i++) {
        initialBoard[1][i] = 'black-pawn';
        initialBoard[6][i] = 'white-pawn';
    }
    // white pieces
    initialBoard[7][0] = 'white-rook';
    initialBoard[7][1] = 'white-knight';
    initialBoard[7][2] = 'white-bishop';
    initialBoard[7][3] = 'white-queen';
    initialBoard[7][4] = 'white-king';
    initialBoard[7][5] = 'white-bishop';
    initialBoard[7][6] = 'white-knight';
    initialBoard[7][7] = 'white-rook';
    // black pieces
    initialBoard[0][0] = 'black-rook';
    initialBoard[0][1] = 'black-knight';
    initialBoard[0][2] = 'black-bishop';
    initialBoard[0][3] = 'black-queen';
    initialBoard[0][4] = 'black-king';
    initialBoard[0][5] = 'black-bishop';
    initialBoard[0][6] = 'black-knight';
    initialBoard[0][7] = 'black-rook';
    return initialBoard;
};

const Board = ({ boardState, onMove }) => {
    // Timing and state refs
    const mouseDownInfo = useRef({ time: 0, row: null, col: null, triggered: false });
    const CLICK_THRESHOLD = 200; // ms

    const handleMouseDown = (row, col) => {
        mouseDownInfo.current = {
            time: Date.now(),
            row,
            col,
            triggered: false
        };
    };

    const handleMouseUp = (row, col) => {
        const now = Date.now();
        const { time, row: downRow, col: downCol, triggered } = mouseDownInfo.current;
        if (triggered) return; // Already handled
        if (downRow === row && downCol === col && now - time < CLICK_THRESHOLD) {
            // Treat as click
            onMove('click', row, col);
        } else {
            // Treat as hold
            onMove('mousedown', downRow, downCol);
            onMove('mouseup', row, col);
        }
        mouseDownInfo.current.triggered = true;
    };

    const handleClick = (row, col) => {
        // Prevent default click if already handled by timing logic
        if (mouseDownInfo.current.triggered) {
            mouseDownInfo.current.triggered = false;
            return;
        }
        onMove('click', row, col);
    };

    const renderSquare = (row, col) => {
        const isLightSquare = (row + col) % 2 === 0;
        const squareClass = isLightSquare ? 'light-square' : 'dark-square';
        const piece = boardState[row][col]; // Get the piece at this position

        return (
            <div
                key={`${row}-${col}`}
                className={`square ${squareClass}`}
                onClick={() => handleClick(row, col)}
                onMouseDown={() => handleMouseDown(row, col)}
                onMouseUp={() => handleMouseUp(row, col)}
            >
                {piece && <div className={`piece ${piece}`}></div>}
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