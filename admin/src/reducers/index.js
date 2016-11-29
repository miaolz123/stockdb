import client from './client';
import { combineReducers } from 'redux';
import { routerReducer as routing } from 'react-router-redux';

const rootReducer = combineReducers({
  client,
  routing,
});

export default rootReducer;
