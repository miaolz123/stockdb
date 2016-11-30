import * as actions from '../constants/actions';
import { Base64 } from 'js-base64';
import { Client } from 'hprose-html5/dist/hprose-html5';

// resetError
export function resetError() {
  return { type: actions.RESET_ERROR };
}

// login
export function login(server, username, password) {
  const token = Base64.encode(`${username}:${password}`);
  return { type: actions.LOGIN, server, token };
}

// requestFailure
function requestFailure(message) {
  return { type: actions.REQUEST_FAILURE, message };
}

// getTimeRange
export function getTimeRange(market, symbol) {
  return (dispatch, getState) => {
    const server = localStorage.getItem('server');
    const token = localStorage.getItem('token');

    dispatch(getTimeRangeRequest());
    if (!server || !token) {
      dispatch(logout());
      return;
    }

    const client = Client.create(server, ['GetTimeRange']);

    client.setHeader('Authorization', `Basic ${token}`);
    client.GetTimeRange(null, (resp) => {
      if (resp.success) {
        dispatch(getTimeRangeSuccess());
      } else {
        dispatch(requestFailure(resp.message));
      }
    }, (name, err) => {
      dispatch(requestFailure('Server error'));
    });
  };
}

function getTimeRangeRequest() {
  return { type: actions.GET_TIME_RANGE_REQUEST };
}

function getTimeRangeSuccess() {
  return { type: actions.GET_TIME_RANGE_SUCCESS };
}

// Logout
export function logout() {
  return { type: actions.LOGOUT };
}
