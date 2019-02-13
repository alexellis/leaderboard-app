module.exports = {
  devServer: {
    proxy: {
      '^/api/*': {
        target: 'https://alexellis.o6s.io',
        changeOrigin: true,
        secure: false,
        pathRewrite: {
          '/api': ''
        }
      }
    }
  }
};
