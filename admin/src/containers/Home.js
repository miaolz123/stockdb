import { getStats } from '../actions';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Table } from 'antd';

class Home extends Component {
  constructor(props) {
    super(props);

    this.refresh = this.refresh.bind(this);
  }

  componentWillMount() {
    this.refresh();
  }

  refresh() {
    const { dispatch } = this.props;

    dispatch(getStats());
  }

  render() {
    const { client } = this.props;
    const columns = [{
      title: 'Market',
      dataIndex: 'Market',
    }, {
      title: 'Disk',
      dataIndex: 'Disk',
      render: (t) => `${(t / 1024).toFixed(2)} KB`,
    }, {
      title: 'Record',
      dataIndex: 'Record',
    }];

    return (
      <Table
        dataSource={client.stats}
        pagination={false}
        columns={columns}
      />
    );
  }
}

const mapStateToProps = (state) => ({
  client: state.client,
});

export default connect(mapStateToProps)(Home);
