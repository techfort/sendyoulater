import { get } from 'axios';

export default (baseUrl) => {
  const session = {
    check: async () => get(`${baseUrl}/check`)
  }
  return {
    session,
  };
};