// src/api/forum.js
import axios from 'axios';

export default {
  // 获取热门文章
  getHotArticles() {
    return axios.get('/api/forum/articles/hot');
  },
  
  // 获取热门标签
  getHotTags() {
    return axios.get('/api/forum/tags/hot');
  },
  
  // 添加其他论坛相关API方法
  getArticles(params) {
    return axios.get('/api/forum/articles', { params });
  },
  
  getArticleDetail(id) {
    return axios.get(`/api/forum/articles/${id}`);
  }
};