import { post, get } from 'axios';
import to from '../helpers';

export default (baseUrl) => {
  const session = {
    check: async () => to(get(`${baseUrl}/check`)),
    getUserData: async () => to(get(`${baseUrl}/user`)),
    login: async (email) => to(post(`${baseUrl}/loginfromfe`, 
      email 
    )),
    loadData: async (user) => to(get(`${baseUrl}/loadData?user=${user}`)),
    saveEmailAction: async (data) => to(post(`${baseUrl}/action/email/save`, data)),
  }
  return {
    session,
  };
};