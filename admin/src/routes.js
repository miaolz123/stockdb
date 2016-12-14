import React from 'react';
import { Route, IndexRoute } from 'react-router';
import App from './containers/App';
import Home from './containers/Home';
import Chart from './containers/Chart';
import Login from './containers/Login';

export default (
  <Route>
    <Route path="/" component={App}>
      <IndexRoute component={Home} />
      <Route path="/chart" component={Chart} />
    </Route>
    <Route path="/login" component={Login} />
  </Route>
);
