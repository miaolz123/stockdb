import { getOHLCs } from '../actions';
import StockChart from '../components/StockChart';
import React, { Component } from 'react';
import { connect } from 'react-redux';

class Home extends Component {
  constructor(props) {
    super(props);

    this.state = {
      innerHeight: window.innerHeight > 500 ? window.innerHeight : 500,
    };
  }

  componentWillMount() {
    const { dispatch } = this.props;

    dispatch(getOHLCs());
  }

  render() {
    const { client } = this.props;

    return (
      <div className="container">
        {client.data.length > 0 ? <StockChart data={client.data} height={this.state.innerHeight - 100} /> : 'HOME'}
      </div>
    );
  }
}

const mapStateToProps = (state) => ({
  client: state.client,
});

export default connect(mapStateToProps)(Home);
