const Path = require('path');
const merge = require('webpack-merge');
const common = require('./webpack.common.js');
const Webpack = require('webpack');

module.exports = merge(common, {
    output: {
        path: Path.resolve(__dirname, 'build'),
        filename: '[name].js'
    },
    devtool: 'inline-source-map',
    devServer: {
        contentBase: './build'
    },
    plugins: [
        new Webpack.NamedModulesPlugin(),        
        new Webpack.DefinePlugin({
            'process.env.NODE_ENV': JSON.stringify('development')
        })
    ]
});
