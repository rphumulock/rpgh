const Path = require('path');
const merge = require('webpack-merge');
const UglifyJSPlugin = require('uglifyjs-webpack-plugin');
const common = require('./webpack.common.js');
const Webpack = require('webpack');

module.exports = merge(common, {
    output: {
        path: Path.resolve(__dirname, 'build'),
        filename: '[name].[chunkhash].js',
        publicPath: './'
    },
    plugins: [
        new UglifyJSPlugin(),
        new Webpack.HashedModuleIdsPlugin(),
        new Webpack.DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify('production')
        })
    ]
});
