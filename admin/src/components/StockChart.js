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
    const { data, height, type, width, ratio } = this.props;
    const margin = { left: 60, right: 0, top: 10, bottom: 20 };
    const accessor = d => { d.date = new Date(d.time * 1000); return d; };
    const dateFormat = timeFormat('%Y-%m-%d %H:%M');
    const title = 'OKCOIN.CN(BTC/CNY), HOUR';
    const showGrid = true;
    const gridHeight = height - margin.top - margin.bottom;
    const gridWidth = width - margin.left - margin.right;
    const yGrid = showGrid ? {
      innerTickSize: -1 * gridWidth,
      tickStrokeDasharray: 'Solid',
      tickStrokeOpacity: 0.1,
      tickStrokeWidth: 1,
    } : {};
    const xGrid = showGrid ? {
      innerTickSize: -1 * gridHeight,
      tickStrokeDasharray: 'Solid',
      tickStrokeOpacity: 0.1,
      tickStrokeWidth: 1,
    } : {};

    return (
      <ChartCanvas
        ratio={ratio}
        width={width}
        height={height}
        margin={margin}
        type={type}
        seriesName="MSFT"
        data={data}
        xAccessor={d => new Date(d.time * 1000)}
        xScaleProvider={discontinuousTimeScaleProvider}
        xExtents={[new Date(2015, 12, 1), new Date(2016, 12, 30)]}
      >
        <Label
          fontSize={36}
          text={title}
          opacity={0.2}
          x={() => (width - margin.left - margin.right) / 2}
          y={() => (height - margin.top - margin.bottom) / 2}
        />
        <Chart
          id={1}
          padding={20}
          yExtents={[d => [d.high, d.low]]}
        >
          <XAxis
            axisAt="bottom"
            orient="bottom"
            {...xGrid}
          />
          <YAxis
            axisAt="left"
            orient="left"
            ticks={9}
            {...yGrid}
          />
          <CandlestickSeries
            opacity={1}
            stroke="#ABAEB7"
            wickStroke="#ABAEB7"
            fill={d => d.close > d.open ? '#00CD7A' : '#CD0027'}
          />
          <MouseCoordinateY
            dx={-60}
            at="right"
            orient="right"
            rectWidth={50}
            arrowWidth={10}
            displayFormat={format('.2f')}
          />
          <OHLCTooltip
            accessor={accessor}
            xDisplayFormat={dateFormat}
            forChart={1} origin={[-40, 0]}
          />
        </Chart>
        <Chart
          id={2}
          height={90}
          yExtents={d => d.volume} origin={(w, h) => [0, h - 90]}
        >
          <MouseCoordinateX
            at="bottom"
            orient="bottom"
            rectWidth={120}
            displayFormat={dateFormat}
          />
          <BarSeries
            opacity={0.2}
            stroke={false}
            yAccessor={d => d.volume}
            fill={(d) => d.close > d.open ? '#00CD7A' : '#CD0027'}
          />
        </Chart>
        <CrossHairCursor />
      </ChartCanvas>
    );
  }
}

StockChart.propTypes = {
  data: React.PropTypes.array.isRequired,
  height: React.PropTypes.number.isRequired,
  width: React.PropTypes.number.isRequired,
  ratio: React.PropTypes.number.isRequired,
  type: React.PropTypes.oneOf(['svg', 'hybrid']).isRequired,
};
StockChart.defaultProps = {
  height: 600,
  type: 'svg',
};
StockChart = fitWidth(StockChart);

export default StockChart;
