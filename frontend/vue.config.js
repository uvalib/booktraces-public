// NOTES on this are found here:
//    https://cli.vuejs.org/config/#devserver
//    https://github.com/chimurai/http-proxy-middleware#proxycontext-config
module.exports = {
  devServer: {
    public: process.env.BASE_URL,
    proxy: {
      '/api': {
        target: process.env.API_HOST, // or 'http://localhost:8095',
        changeOrigin: true,
        logLevel: 'debug'
      },
      '/authenticate': {
        target: process.env.API_HOST, // or 'http://localhost:8095',
        changeOrigin: true,
        logLevel: 'debug'
      },
      '/uploads': {
        target: process.env.API_HOST, // or 'http://localhost:8095',
        changeOrigin: true,
        logLevel: 'debug'
      }
    }
  },
  configureWebpack: {
    performance: {
      // bump max sizes to 512k
      maxEntrypointSize: 512000,
      maxAssetSize: 512000
    }
  }
}
