var cooking = require('cooking');
var path = require('path');

cooking.set({
  entry: {
    app: './src/app.js'
  },
  dist: './dist',
  template: 'src/index.tpl',

  // development
  devServer: {
    port: 8080,
    publicPath: '/'
  },

  // production
  clean: true,
  hash: true,
  chunk: true,
  publicPath: '/dist/',
  assetsPath: 'static',
  sourceMap: true,
  extractCSS: true,
  urlLoaderLimit: 10000,
  postcss: [
    // require('postcss-cssnext')
  ],
  alias: {
    'src': path.join(__dirname, 'src')
  },
  extends: ['react', 'lint']
});

module.exports = cooking.resolve();
