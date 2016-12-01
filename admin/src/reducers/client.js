import * as actions from '../constants/actions';
import assign from 'lodash/assign';

const CLIENT_INIT = {
  loading: false,
  message: '',
  data: []
};

function user(state = CLIENT_INIT, action) {
  switch (action.type) {
    case actions.RESET_ERROR:
      return assign({}, state, {
        loading: false,
        message: '',
      });
    case actions.LOGIN:
      localStorage.setItem('server', action.server);
      localStorage.setItem('token', action.token);
      return state;
    case actions.REQUEST_FAILURE:
      if (action.message === 'Unauthorized') {
        localStorage.removeItem('token');
      }
      return assign({}, state, {
        loading: false,
        message: action.message,
      });
    case actions.GET_TIME_RANGE_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.GET_TIME_RANGE_SUCCESS:
      return assign({}, state, {
        loading: false,
      });
    case actions.GET_OHLCS_REQUEST:
      return assign({}, state, {
        loading: true,
      });
    case actions.GET_OHLCS_SUCCESS:
      return assign({}, state, {
        loading: false,
        data: action.data,
      });
    case actions.LOGOUT:
      localStorage.removeItem('token');
      return CLIENT_INIT;
    default:
      return state;
  }
}

export default user;
