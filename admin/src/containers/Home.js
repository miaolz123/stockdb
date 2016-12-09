import { getOHLCs } from '../actions';
import StockChart from '../components/StockChart';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Card, Row, Col, Cascader, Radio, Spin, Tooltip } from 'antd';

class Home extends Component {
  constructor(props) {
    super(props);

    this.state = {
      innerHeight: window.innerHeight > 500 ? window.innerHeight : 500,
      innerWidth: window.innerWidth > 500 ? window.innerWidth : 500,
    };

    this.symbol = ['', ''];
    this.period = 60;

    this.refresh = this.refresh.bind(this);
    this.onSymbolChange = this.onSymbolChange.bind(this);
    this.onPeriodChange = this.onPeriodChange.bind(this);
  }

  componentWillReceiveProps(nextProps) {
    const { client } = nextProps;
    const newSymbol = [];

    if (this.symbol[0] === '' && client.symbols.length > 0) {
      newSymbol.push(client.symbols[0].value);
      if (client.symbols[0].children.length > 0) {
        newSymbol.push(client.symbols[0].children[0].value);
        this.symbol = newSymbol;
        this.refresh();
      }
    }
  }

  componentWillMount() {
    this.refresh();
  }

  refresh() {
    const { dispatch } = this.props;
    const { symbol, period } = this;

    dispatch(getOHLCs(symbol, period));
  }

  onSymbolChange(symbol) {
    this.symbol = symbol;
    this.refresh();
  }

  onPeriodChange(e) {
    this.period = Number(e.target.value);
    this.refresh();
  }

  render() {
    const { client } = this.props;
    const { innerHeight, innerWidth } = this.state;
    const klineAmount = parseInt(innerWidth / 10, 10);
    const { symbol, period } = this;
    const data = [];
    const displayRender = symbol => `${symbol[1]} @ ${symbol[0]}`;

    if (client.data) {
      client.data.forEach(d => {
        data.push({
          time: d.Time,
          open: d.Open,
          high: d.High,
          low: d.Low,
          close: d.Close,
          volume: d.Volume,
        });
      });
    }

    if (client.timeRange[1] > 0 && client.timeRange[1] - klineAmount * period > client.timeRange[0]) {
      client.timeRange[0] = client.timeRange[1] - klineAmount * period;
    }

    return (
      <Card bordered={false}>
        <Row className="chart-header">
          <Col span={6}>
            <Cascader
              size="small"
              value={symbol}
              allowClear={false}
              expandTrigger="hover"
              options={client.symbols}
              displayRender={displayRender}
              onChange={this.onSymbolChange}
            />
          </Col>
          <Col span={18} style={{ textAlign: 'right' }}>
            <Tooltip title="Change Period" placement="left">
              <Radio.Group
                size="small"
                defaultValue={String(period)}
                onChange={this.onPeriodChange}
              >
                <Radio.Button value="60">M</Radio.Button>
                <Radio.Button value="300">5M</Radio.Button>
                <Radio.Button value="900">15M</Radio.Button>
                <Radio.Button value="1800">30M</Radio.Button>
                <Radio.Button value="3600">H</Radio.Button>
                <Radio.Button value="28800">8H</Radio.Button>
                <Radio.Button value="86400">D</Radio.Button>
              </Radio.Group>
            </Tooltip>
          </Col>
        </Row>
        <Spin size="large" spinning={client.loading}>
          {data.length > 0
          ? <StockChart
              data={data}
              symbol={symbol}
              period={period}
              timeRange={client.timeRange}
              height={innerHeight - 128}
            />
          : <div><br /><br /><br /><br /><br /></div>}
        </Spin>
      </Card>
    );
  }
}

const mapStateToProps = (state) => ({
  client: state.client,
});

export default connect(mapStateToProps)(Home);
