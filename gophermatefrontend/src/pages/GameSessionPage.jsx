import { useState, useEffect, useRef } from 'react';
import React from 'react';
import { useParams } from 'react-router-dom';

import Board, {InitializeBoard} from '../chess/board';
import './GameSessionPage.css';


const GameSessionPage = () => {
    const { id } = useParams();
    const [boardState, setBoardState] = useState(InitializeBoard());
    const [selected, setSelected] = useState(null); // For click-based selection
    const dragStart = useRef(null); // For mousedown/mouseup drag

    function postMove(piece, from, to) {
        console.log('Posting move:', piece, from, to);
        // Send the move to the server
        fetch('/api/games/move', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            // sends in the format {'id': 'sa65s4165a4', 'piece': 'black-pawn', 'from': {'row': 1, 'col': 7}, 'to': {'row': 3, 'col': 7}}
            body: JSON.stringify({'session': id, 'piece': piece, 'from': from, 'to': to}),
        })
        .then(response => {
            if (!response.ok) {
                console.log('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            console.log('Move posted successfully:', data);
            return true;
        })
        .catch(error => {
            console.error('Error posting move:', error);
        });
        return false;
    }

    function onMove(action, row, col) {
        if (action === 'click') {
            if (!selected) {
                // First click: select piece if present
                if (boardState[row][col]) {
                    setSelected({ row, col });
                }
            } else {
                // Second click: move piece
                const from = selected;
                const to = { row, col };
                if (from.row !== to.row || from.col !== to.col) {
                    // Post the move to the server
                    const moveOk = postMove(boardState[from.row][from.col], from, to);
                    if (moveOk) {
                        const newBoard = boardState.map(r => r.slice());
                        newBoard[to.row][to.col] = newBoard[from.row][from.col];
                        newBoard[from.row][from.col] = null;
                        setBoardState(newBoard);
                    }
                }
                setSelected(null);
            }
        } else if (action === 'mousedown') {
            // Start drag
            if (boardState[row][col]) {
                dragStart.current = { row, col };
            } else {
                dragStart.current = null;
            }
        } else if (action === 'mouseup') {
            // End drag and move
            const from = dragStart.current;
            if (from && (from.row !== row || from.col !== col)) {
                const moveOk = postMove(boardState[from.row][from.col], from, to);
                if (moveOk) {
                    const newBoard = boardState.map(r => r.slice());
                    newBoard[row][col] = newBoard[from.row][from.col];
                    newBoard[from.row][from.col] = null;
                    setBoardState(newBoard);
                }
            }
            dragStart.current = null;
        }
    }

    return (
        <div className="game-session-page">
            <h1 className="title">Game Session</h1>
            <div className="content">
                <div className="game-board">
                    <Board boardState={boardState} onMove={onMove} />
                </div>
                <div className="side-actions">
                    <div className="game-controls">
                        <button onClick={() => alert('Resign')}>Resign</button>
                        <button onClick={() => alert('Offer Draw')}>Offer Draw</button>
                    </div>
                    <div className="chat-box">
                        <h2>Chat</h2>
                        <div className="messages">
                            {/* Chat messages will be displayed here */}
                        </div>
                        <input type="text" placeholder="Type a message..." />
                        <button onClick={() => alert('Send message')}>Send</button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default GameSessionPage;