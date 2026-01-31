const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8002,
    // Webpack 5 使用 allowedHosts 代替 disableHostCheck
    allowedHosts: 'all', 
    historyApiFallback: true,
    proxy: {
      // 1. 转发所有 API 请求到 Go 后端
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
      // 2. 转发课程图片请求 (对应后端日志里的 /img/course/)
      '/img': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
      // 3. 转发通用图片请求 (对应后端日志里的 /images/)
      '/images': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      }
    },
  },
  configureWebpack: {
    resolve: {
      alias: {
        "assets": "@/assets",
      },
    },
    module: {
      rules: [
        // 处理媒体文件 (代替原来的 url-loader)
        {
          test: /\.(mp4|webm|ogg|mp3|wav|flac|aac)$/,
          type: 'asset',
          parser: {
            dataUrlCondition: {
              maxSize: 10 * 1024 // 10KB
            }
          },
          generator: {
            filename: 'videos/[name].[hash][ext]'
          }
        },
        // 处理 PDF 文件 (代替原来的 file-loader)
        {
          test: /\.(pdf)(\?.*)?$/,
          type: 'asset/resource',
          generator: {
            filename: 'assets/pdf/[name].[hash:8][ext]'
          }
        }
      ],
    },
  },
  css: {
    loaderOptions: {
      scss: {
        // sass-loader 高版本中使用 additionalData
        additionalData: `@import "@/assets/css/common.scss";`
      }
    }
  }
})