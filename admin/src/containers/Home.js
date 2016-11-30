import StockChart from '../components/StockChart';
import React, { Component } from 'react';
import { connect } from 'react-redux';

class Home extends Component {
  render() {
    const data = [
      {
        date: new Date(2016, 11, 1),
        open: 123456,
        high: 123456,
        low: 123456,
        close: 110,
        volume: 123456,
      },
      {
        date: new Date(2016, 11, 2),
        open: 123456,
        high: 123456,
        low: 123456,
        close: 120,
        volume: 123456,
      },
      {
        date: new Date(2016, 11, 3),
        open: 123456,
        high: 123456,
        low: 123456,
        close: 130,
        volume: 123456,
      },
      {
        date: new Date(2016, 11, 4),
        open: 123456,
        high: 123456,
        low: 123456,
        close: 140,
        volume: 123456,
      }
    ];

    return (
      <div className="container">
          <StockChart data={data} />
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  client: state.client,
});

export default connect(mapStateToProps)(Home);
