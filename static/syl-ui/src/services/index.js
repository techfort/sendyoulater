import { get } from 'axios';

const to = promise => promise.then(data => ({ error: null, data }))
  .catch(error => ({ error, data: null }));

export default (baseUrl) => {
  const session = {
    check: async () => to(get(`${baseUrl}/check`)),
    getUserData: async () => to(get(`${baseUrl}/user`)),
  }
  return {
    session,
  };
};