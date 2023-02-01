const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin')

module.exports = function override(config, env) {
  config.plugins.push(
    new MonacoWebpackPlugin({
      languages: ['yaml'],
      globalAPI: true,
    })
  )

  config.optimization = {
    splitChunks: {
      chunks: 'all',
      cacheGroups: {
        giraffe: {
          test: /Giraffe[\\/]/,
          name: 'giraffe',
        },
      },
    },
  }

  return config
}
