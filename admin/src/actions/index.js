import * as actions from '../constants/actions';
import StockDB from 'stockdb';

// resetError
export function resetError() {
  return { type: actions.RESET_ERROR };
}

// login
export function login(server, username, password) {
  const token = window.btoa(`${username}:${password}`);
  return { type: actions.LOGIN, server, token };
}

// requestFailure
function requestFailure(message) {
  return { type: actions.REQUEST_FAILURE, message };
}

// getStats
export function getStats() {
  return (dispatch, getState) => {
    const server = localStorage.getItem('server');
    const token = localStorage.getItem('token');

    dispatch(getStatsRequest());
    if (!server || !token) {
      dispatch(logout());
      return;
    }

    const client = StockDB.New(server, window.atob(token));

    client.GetStats((resp) => {
      if (resp.Success) {
        dispatch(getStatsSuccess(resp.Data));
      } else {
        dispatch(requestFailure(resp.Message));
      }
    }, (name, err) => {
      dispatch(requestFailure('Server error'));
    });
  };
}

function getStatsRequest() {
  return { type: actions.GET_STATS_REQUEST };
}

function getStatsSuccess(stats) {
  return { type: actions.GET_STATS_SUCCESS, stats };
}

// getSymbols
export function getSymbols() {
  return (dispatch, getState) => {
    const server = localStorage.getItem('server');
    const token = localStorage.getItem('token');

    dispatch(getSymbolsRequest());
    if (!server || !token) {
      dispatch(logout());
      return;
    }

    const client = StockDB.New(server, window.atob(token));

    client.GetMarkets((resp) => {
      if (resp.Success) {
        dispatch(getMarketsSuccess(resp.Data));
        resp.Data.forEach((m, i) => {
          client.GetSymbols(m, (resp) => {
            if (resp.Success) {
              dispatch(getSymbolsSuccess(i, resp.Data));
            } else {
              dispatch(requestFailure(resp.Message));
            }
          }, (name, err) => {
            dispatch(requestFailure('Server error'));
          });
        });
      } else {
        dispatch(requestFailure(resp.Message));
      }
    }, (name, err) => {
      dispatch(requestFailure('Server error'));
    });
  };
}

function getSymbolsRequest() {
  return { type: actions.GET_SYMBOLS_REQUEST };
}

function getMarketsSuccess(markets) {
  return { type: actions.GET_MARKETS_SUCCESS, markets };
}

function getSymbolsSuccess(index, symbols) {
  return { type: actions.GET_SYMBOLS_SUCCESS, index, symbols };
}

// getTimeRange
function getTimeRange(opt) {
  return (dispatch, getState) => {
    const server = localStorage.getItem('server');
    const token = localStorage.getItem('token');

    dispatch(getTimeRangeRequest());
    if (!server || !token) {
      dispatch(logout());
      return;
    }

    const client = StockDB.New(server, window.atob(token));

    client.GetTimeRange(opt, (resp) => {
      if (resp.Success) {
        dispatch(getTimeRangeSuccess(resp.Data));
      } else {
        dispatch(requestFailure(resp.Message));
      }
    }, (name, err) => {
      dispatch(requestFailure('Server error'));
    });
  };
}

function getTimeRangeRequest() {
  return { type: actions.GET_TIME_RANGE_REQUEST };
}

function getTimeRangeSuccess(timeRange) {
  return { type: actions.GET_TIME_RANGE_SUCCESS, timeRange };
}

// getPeriodRange
export function getPeriodRange(symbol) {
  return (dispatch, getState) => {
    const server = localStorage.getItem('server');
    const token = localStorage.getItem('token');

    dispatch(getPeriodRangeRequest());
    if (!server || !token) {
      dispatch(logout());
      return;
    }

    const client = StockDB.New(server, window.atob(token));
    const opt = { Market: symbol[0], Symbol: symbol[1] };

    client.GetPeriodRange(opt, (resp) => {
      if (resp.Success) {
        dispatch(getPeriodRangeSuccess(resp.Data));
      } else {
        dispatch(requestFailure(resp.Message));
      }
    }, (name, err) => {
      dispatch(requestFailure('Server error'));
    });
  };
}

function getPeriodRangeRequest() {
  return { type: actions.GET_PERIOD_RANGE_REQUEST };
}

function getPeriodRangeSuccess(periodRange) {
  return { type: actions.GET_PERIOD_RANGE_SUCCESS, periodRange };
}

// getOHLCs
export function getOHLCs(symbol, period) {
  return (dispatch, getState) => {
    const server = localStorage.getItem('server');
    const token = localStorage.getItem('token');

    dispatch(getOHLCsRequest());
    if (!server || !token) {
      dispatch(logout());
      return;
    }

    const client = StockDB.New(server, window.atob(token));
    const opt = { Market: symbol[0], Symbol: symbol[1], Period: period };

    if (opt.Market !== '' && period > 0) {
      client.GetOHLCs(opt, (resp) => {
        if (resp.Success) {
          dispatch(getTimeRange(opt));
          dispatch(getOHLCsSuccess(resp.Data));
        } else {
          dispatch(requestFailure(resp.Message));
        }
      }, (name, err) => {
        dispatch(requestFailure('Server error'));
      });
    }
  };
}

function getOHLCsRequest() {
  return { type: actions.GET_OHLCS_REQUEST };
}

function getOHLCsSuccess(data) {
  return { type: actions.GET_OHLCS_SUCCESS, data };
}

// Logout
export function logout() {
  return { type: actions.LOGOUT };
}
