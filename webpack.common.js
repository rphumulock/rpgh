const Path = require('path');
const ExtractTextPlugin = require("extract-text-webpack-plugin");
const HtmlWebpackPlugin = require('html-webpack-plugin');
const CleanWebpackPlugin = require('clean-webpack-plugin');
const Webpack = require('webpack');

module.exports = {
    entry: {
        index: './src/js/index.js',
        projects: './src/js/projects.js',
        common: [
            'materialize-css/dist/js/materialize.min.js'
        ]
    },
    module: {
        rules: [
            {
                test: /\.jsx?$/,
                exclude: /(node_modules|bower_components)/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['babel-preset-env']
                    }
                }
            },
            {
                test: /\.css/,
                use: ExtractTextPlugin.extract({
                    fallback: "style-loader",
                    use: "css-loader"
                })
            },
            {
                test: /\.(woff|woff2|eot|ttf)$/,
                loader: 'file-loader'
            },
            {
                test: /\.(mp4|webm)$/,
                loader: 'file-loader'
            },
            {
                test: /\.(gif|png|svg|jpe?g)$/,
                use: [
                    {
                        loader: 'url-loader',
                        options: {
                            limit: 40000
                        }
                    },
                    'image-webpack-loader'
                ]
            }
        ]
    },
    plugins: [
        new CleanWebpackPlugin(['docs']),
        new ExtractTextPlugin("styles.css"),
        new Webpack.optimize.CommonsChunkPlugin({
            name: 'common' // Specify the common bundle's name.
        }),
        new Webpack.optimize.CommonsChunkPlugin({
            name: 'runtime'
        }),
        new HtmlWebpackPlugin({
            template: './src/index.html',
            filename: './index.html',
            chunks: ['runtime', 'common', 'index']
        }),
        new HtmlWebpackPlugin({
            template: './src/projects.html',
            filename: './projects.html',
            chunks: ['runtime', 'common', 'index']
        }),
        new Webpack.ProvidePlugin({ // inject ES5 modules as global vars
            $: 'jquery',
            jQuery: 'jquery',
            'window.jQuery': 'jquery',
            Tether: 'tether'
        })
    ]
};



/*               /* use: [
                    {
                        loader: 'url-loader',
                        options: {
                            limit: 40000
                        }
                    },
                    'image-webpack-loader'
                ]
                
                [
                    {
                        loader: 'url-loader',
                        options: {
                            limit: 4000000
                        }
                    },
                    {
                        loader: 'image-webpack-loader',
                        options: {
                            optipng: {
                                enabled: true,
                            }
                        }
                    }
                ]*/