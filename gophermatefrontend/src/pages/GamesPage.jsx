import React, { useEffect, useState } from 'react';

const GamesPage = () => {
  const [games, setGames] = useState([]);

  useEffect(() => {
    fetch('/api/games')
      .then((response) => response.json())
      .then((data) => {
        const gamesWithStatus = data.map((game) => ({
          ...game,
          status: game.playerBlack ? 'In Progress' : 'Open',
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
            Game ID: {game.id} - Status: {game.status}
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