import axios from 'axios';

export const API_URL = 'http://localhost:8080';

export const registerUser = async (userData) => {
  try {
    const response = await axios.post(API_URL + '/api/register', userData);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : new Error('Network error');
  }
};

export const loginUser = async (credentials) => {
  try {
    const response = await axios.post(API_URL + '/api/login', credentials);
    return response.data;
  } catch (error) {
    throw error.response ? error.response.data : new Error('Network error');
  }
};