import axios from 'axios';
import { API_URL } from './authService';


export const postMove = async ({ session, user, piece, from, to }) => {
  try {
    const response = await axios.post(`${API_URL}/api/games/move`, {
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
