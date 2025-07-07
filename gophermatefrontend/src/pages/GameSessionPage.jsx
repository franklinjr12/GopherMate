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
    const userToken = localStorage.getItem('token'); // Assuming user token is stored in localStorage

    async function postMove(piece, from, to) {
        console.log('Posting move:', piece, from, to);
        try {
            const response = await fetch('/api/games/move', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({'session': id, 'user': userToken, 'piece': piece, 'from': from, 'to': to}),
            });
            if (!response.ok) {
                console.log('Network response was not ok');
                return false;
            }
            const data = await response.json();
            console.log('Move posted successfully:', data);
            return true;
        } catch (error) {
            console.error('Error posting move:', error);
            return false;
        }
    }

    async function onMove(action, row, col) {
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
                    const moveOk = await postMove(boardState[from.row][from.col], from, to);
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
                const moveOk = await postMove(boardState[from.row][from.col], from, { row, col });
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