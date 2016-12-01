import { getOHLCs } from '../actions';
import StockChart from '../components/StockChart';
import React, { Component } from 'react';
import { connect } from 'react-redux';

class Home extends Component {
  componentWillMount() {
    const { dispatch } = this.props;

    dispatch(getOHLCs());
  }

  render() {
    const { client } = this.props;

    return (
      <div className="container">
        {client.data.length > 0 ? <StockChart data={client.data} /> : 'HOME'}
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  client: state.client,
});

export default connect(mapStateToProps)(Home);
