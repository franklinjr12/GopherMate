import React from 'react';
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