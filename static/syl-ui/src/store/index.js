import Vuex from 'vuex'
import Vue from 'vue'

Vue.use(Vuex)


const m = {
  USER: 'USER',
  EVENTS: 'EVENTS',
  EVENT: 'EVENT',
  ADD_EVENTS: 'ADD_EVENTS',
  REMOVE_EVENT: 'REMOVE_EVENT',
  SET_USER: 'SET_USER',
  RESET_USER: 'RESET_USER'
};

const state = {
  user: null,
  events: null,
};

const mutations = {
  [m.SET_USER] (state, user) {
    state.user = user;
    return state;
  },
  [m.RESET_USER] (state) {
    state.user = {};
    return state;
  },
  [m.ADD_EVENT] (state, event) {
    state.events[event.id] = event;
    return state;
  },
  [m.REMOVE_EVENT] (state, event) {
    state.events[event.id] = null;
    return state;
  }
};

const actions = {
  setUser({ commit }, user) {
    return commit(m.SET_USER, user);
  },
  resetUser({ commit }) {
    return commit(m.RESET_USER)
  }
};

const getters = {
  [m.USER]: (state) => state.user,
  [m.EVENTS]: (state) => state.events,
  [m.EVENT]: (state, id) => state.events[id],
};

export default new Vuex.Store({
  state,
  mutations,
  actions,
  getters,
});