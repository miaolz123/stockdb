import React from 'react';
import { scaleTime } from 'd3-scale';
import { ChartCanvas, Chart, series, axes, helper } from 'react-stockcharts';

const { AreaSeries } = series;
const { XAxis, YAxis } = axes;
const { fitWidth, TypeChooser } = helper;

class StockChart extends React.Component {
  render() {
    const { data, type, width, ratio } = this.props;
    return (
      <TypeChooser type="hybrid">
        {type => <ChartCanvas ratio={ratio} width={width} height={400}
            margin={{ left: 50, right: 50, top: 10, bottom: 30 }}
            seriesName="MSFT"
            data={data} type={type}
            xAccessor={d => d.date}
            xScale={scaleTime()}
            xExtents={[new Date(2016, 11, 2), new Date(2016, 11, 3)]}>
          <Chart id={0} yExtents={d => d.close}>
            <XAxis axisAt="bottom" orient="bottom" ticks={6}/>
            <YAxis axisAt="left" orient="left" />
            <AreaSeries yAccessor={(d) => d.close}/>
          </Chart>
        </ChartCanvas>}
      </TypeChooser>
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
