import { useState, useEffect } from 'react';
import React from 'react';

import Board from '../chess/board';
import './GameSessionPage.css';

const GameSessionPage = () => {
    const [boardState, setBoardState] = useState([]);

    useEffect(() => {
        // Initialize board state or fetch game data
    }, []);

    function onMove(row, col) {
        // Handle the move logic here
        console.log(`Move made to row: ${row}, col: ${col}`);
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