import React, { useEffect, useState } from 'react';

const GamesPage = () => {
  const [games, setGames] = useState([]);

  useEffect(() => {
    fetch('http://localhost:8080/api/games')
      .then((response) => response.json())
      .then((data) => {
        const gamesWithStatus = data.map((game) => ({
          ...game,
          status: game.player_black ? 'In Progress' : 'Open',
        }));
        setGames(gamesWithStatus);
        console.log('Fetched games:', gamesWithStatus);
      });
  }, []);

  const joinGame = (id) => {
    fetch(`http://localhost:8080/api/games/${id}/join`, { method: 'POST' })
      .then((response) => response.json())
      .then((data) => {
        if (data.message) {
          // Redirect to the game session page after joining
          window.location.href = `/gamesession/${id}`;
        } else {
          alert(data.error || 'Failed to join game');
        }
      });
  };

  const createGame = () => {
    fetch('http://localhost:8080/api/games', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ player_token: localStorage.getItem('token') || ''}),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.id) {
          console.log('Create game response:', data);
          // Redirect to the new game session page
          window.location.href = `/gamesession/${data.id}`;
        } else {
          alert(data.message || 'Failed to create game');
        }
      })
      .catch(() => alert('Failed to create game'));
  };

  return (
    <div className="games-page">
      <div className="header">
        <button onClick={createGame}>Create Game</button>
        <button onClick={() => window.location.href = '/logout'}>Logout</button>
      </div>
      <div className="games-list">
        <h1>Available Games</h1>
        <ul>
          {games.map((game) => (
            <li key={game.id}>
              Game ID: {game.id} - Status: {game.status} - Player White: {game.player_white} - Player Black: {game.player_black}
              {game.status === 'Open' && (
                <button onClick={() => joinGame(game.id)}>Join</button>
              )}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
};

export default GamesPage;