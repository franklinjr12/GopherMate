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
    fetch(`/api/games/${id}/join`, { method: 'POST' })
      .then((response) => response.json())
      .then((data) => alert(data.message));
  };

  return (
    <div>
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
  );
};

export default GamesPage;