import axios from 'axios';

const API_ENDPOINT = process.env.REACT_APP_API_ENDPOINT;
if (API_ENDPOINT == null) throw new Error('No REACT_APP_API_ENDPOINT provided!');

function e(name) {
  return `${API_ENDPOINT}${name}`;
}

export async function getHome() {
  const res = await axios.get(e('home'));
  return res.data;
}
