import React from 'react';
import { format } from 'd3-format';
import { timeFormat } from 'd3-time-format';
import { ChartCanvas, Chart } from 'react-stockcharts';
import { CandlestickSeries, BarSeries } from 'react-stockcharts/lib/series';
import { discontinuousTimeScaleProvider } from 'react-stockcharts/lib/scale';
import { CrossHairCursor, MouseCoordinateX, MouseCoordinateY } from 'react-stockcharts/lib/coordinates';
import { Label } from 'react-stockcharts/lib/annotation';
import { OHLCTooltip } from 'react-stockcharts/lib/tooltip';
import { XAxis, YAxis } from 'react-stockcharts/lib/axes';
import { fitWidth } from 'react-stockcharts/lib/helper';

class StockChart extends React.Component {
  render() {
    const { data, type, width, ratio } = this.props;
    const margin = { left: 60, right: 60, top: 10, bottom: 20 };
    const accessor = d => { d.date = new Date(d.time * 1000); return d; };
    const dateFormat = timeFormat('%Y-%m-%d %H:%M');
    const height = 650;
    const title = 'OKCOIN.CN - BTC/CNY';

    return (
      <ChartCanvas ratio={ratio} width={width} height={height}
          margin={margin} type={type}
          seriesName="MSFT"
          data={data}
          xAccessor={d => new Date(d.time * 1000)} xScaleProvider={discontinuousTimeScaleProvider}
          xExtents={[new Date(2015, 12, 1), new Date(2016, 12, 30)]}>
        <Label
          x={() => (width - margin.left - margin.right) / 2}
          y={() => (height - margin.top - margin.bottom) / 2}
          fontSize={36} text={title} opacity={0.2} />
        <Chart id={1} yExtents={[d => [d.high, d.low]]} padding={20}>
          <XAxis axisAt="bottom" orient="bottom"/>
          <YAxis axisAt="left" orient="left" ticks={9} />
          <MouseCoordinateY
            at="right"
            orient="right"
            displayFormat={format('.2f')} />
          <CandlestickSeries fill={(d) => d.close > d.open ? '#6BA583' : '#D75442'} opacity={1} />
          <OHLCTooltip accessor={accessor} xDisplayFormat={dateFormat} forChart={1} origin={[-40, 0]}/>
        </Chart>
        <Chart id={2} height={90}
            yExtents={d => d.volume}
            origin={(w, h) => [0, h - 90]}>
          <MouseCoordinateX
            at="bottom"
            orient="bottom"
            rectWidth={120}
            displayFormat={dateFormat} />
          <BarSeries
            yAccessor={d => d.volume}
            fill={(d) => d.close > d.open ? '#6BA583' : '#D75442'}
            opacity={0.2}
            stroke={false} />
        </Chart>
        <CrossHairCursor />
      </ChartCanvas>
    );
  }
}

StockChart.propTypes = {
  data: React.PropTypes.array.isRequired,
  width: React.PropTypes.number.isRequired,
  ratio: React.PropTypes.number.isRequired,
  type: React.PropTypes.oneOf(['svg', 'hybrid']).isRequired,
};
StockChart.defaultProps = {
  type: 'svg',
};
StockChart = fitWidth(StockChart);

export default StockChart;
