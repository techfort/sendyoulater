import { post, get } from 'axios';
import to from '../helpers';

export default (baseUrl) => {
  const session = {
    check: async () => to(get(`${baseUrl}/check`)),
    getUserData: async () => to(get(`${baseUrl}/user`)),
    login: async (email) => to(post(`${baseUrl}/loginfromfe`, 
      email 
    )),
    loadData: async () => to(get(`${baseUrl}/loadData`)),
  }
  return {
    session,
  };
};