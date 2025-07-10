import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

export const postMove = async ({ session, user, piece, from, to }) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/games/move`, {
      session,
      user,
      piece,
      from,
      to,
    });
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : new Error('Network error');
  }
};
