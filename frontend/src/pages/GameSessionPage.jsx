import { useState, useEffect, useRef } from 'react';
import React from 'react';
import { useParams } from 'react-router-dom';
import Board, { InitializeBoard } from '../chess/board';
import { postMove as postMoveApi } from '../services/gameService';
import './GameSessionPage.css';


const GameSessionPage = () => {
    const { id } = useParams();
    const [boardState, setBoardState] = useState(InitializeBoard());
    const [selected, setSelected] = useState(null); // For click-based selection
    const dragStart = useRef(null); // For mousedown/mouseup drag
    const userToken = localStorage.getItem('token'); // Assuming user token is stored in localStorage
    const [lastMoveNumber, setLastMoveNumber] = useState(0); // Track last move number
    const [lastMoveNotation, setLastMoveNotation] = useState(''); // Track last move notation
    const [turn, setTurn] = useState('white'); // Track whose turn it is

    // Helper to parse notation like "white-pawn e2->e4" and update boardState
    function applyNotationToBoard(notation) {
        if (!notation) return;
        // Example: "white-pawn e2->e4"
        const match = notation.match(/^(\w+-\w+)\s+([a-h][1-8])->([a-h][1-8])$/);
        if (!match) return;
        const [, piece, from, to] = match;
        // Convert algebraic to board indices
        function algebraicToIndex(square) {
            const col = square.charCodeAt(0) - 'a'.charCodeAt(0);
            const row = 7 - (parseInt(square[1], 10) - 1);
            return { row, col };
        }
        const fromIdx = algebraicToIndex(from);
        const toIdx = algebraicToIndex(to);
        setBoardState(prev => {
            const newBoard = prev.map(r => r.slice());
            newBoard[toIdx.row][toIdx.col] = newBoard[fromIdx.row][fromIdx.col];
            newBoard[fromIdx.row][fromIdx.col] = null;
            return newBoard;
        });
    }

    useEffect(() => {
        let isMounted = true;
        let intervalId = null;
        async function fetchBoard() {
            try {
                const res = await fetch(`http://localhost:8080/api/games/${id}/board`, {
                    headers: {
                        'Authorization': userToken ? `Bearer ${userToken}` : undefined,
                    },
                });
                if (!res.ok) return;
                const data = await res.json();
                if (isMounted) {
                    // format is { number: 1, notation: "white-pawn e2->e4" }
                    if (lastMoveNumber !== data.number) {
                        const turn = data.notation.split(' ')[0].split('-')[0]; // e.g., "white" from "white-pawn e2->e4"
                        setTurn(turn);
                        setLastMoveNumber(data.number);
                        setLastMoveNotation(data.notation);
                        applyNotationToBoard(data.notation);
                    }
                }
            } catch (e) {
                console.log('Error fetching board:', e);
            }
        }
        fetchBoard();
        intervalId = setInterval(fetchBoard, 1000);
        return () => {
            isMounted = false;
            if (intervalId) clearInterval(intervalId);
        };
    }, [id, userToken, lastMoveNumber]);

    async function postMove(piece, from, to) {
        try {
            const data = await postMoveApi({
                session: id,
                user: userToken,
                piece,
                from,
                to,
            });
            setLastMoveNumber(lastMoveNumber + 1);
            return true;
        } catch (error) {
            alert('Invalid move: ' + error.error);
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
            <div style={{ marginBottom: '10px', fontWeight: 'bold', fontSize: '1.2em' }}>
                {turn === 'white' || turn === 'black' ? `Current turn: ${turn.charAt(0).toUpperCase() + turn.slice(1)}` : 'Current turn: Unknown'}
            </div>
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