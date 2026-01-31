<template>
  <div>
      <div class="banner">
      <div class="banner-container center-content">
        <div class="banner-wrapper">
          <el-carousel :height="bannerHeight + 'px'">
            <el-carousel-item v-for="item in banner" :key="item.id">
              <div class="banner-image-container">
                <!-- 背景模糊层 -->
                <div class="blur-background">
                  <img :src="item.imageUrl" :alt="item.title" />
                </div>
                <!-- 前景清晰图片 -->
                <div class="foreground-image">
                  <img :src="item.imageUrl" :alt="item.title" />
                </div>
                <!-- 左右渐变遮罩 -->
                <div class="gradient-overlay left-overlay"></div>
                <div class="gradient-overlay right-overlay"></div>
              </div>
            </el-carousel-item>
          </el-carousel>
        </div>
      </div>
    </div>
    <div class="friendlinks homemok">
      <div class="flink-content center-content">
        <div
          class="flink-item fadeInUp"
          v-for="(item, key) in blogroll"
          :key="item.id"
          :style="{ animationDelay: key * 0.1 + 's' }"
          @click="goLinkClick(item.link)"
        >
          <div class="flink-item-img">
            <img :src="item.cover" alt="" />
          </div>
          <div class="flink-item-text">
            {{ item.title }}
          </div>
          <div class="flink-item-desc">
            {{ item.desc }}
          </div>
        </div>
      </div>
    </div>
    <div class="recommend homemok">
      <div class="recommend-content center-content">
        <div class="left-notice recommend-item">
          <div class="notice-title">
            <span class="title-text">
              <svg class="icon-font">
                <use xlink:href="#icon-gonggao"></use>
              </svg>
              最新公告
            </span>
            <a
              class="title-right"
              target="_blank"
              @click="$router.push(`/news/list`)"
            >
              查看更多
              <i class="el-icon-d-arrow-right"></i
            ></a>
          </div>
          <div class="notice-list">
            <div
              class="notice-list-item"
              v-for="item in affiche"
              :key="item.id"
              @click="$router.push(`/news/detail/${item.id}`)"
            >
              <a class="noticeitem-text ellipsis" :title="item.title"
                ><i class="iconfont icon-yuandian"></i>
                {{ item.title }}
              </a>
              <div class="noticeitem-time">
                {{ new Date(item.updateTime).format() }}
              </div>
            </div>
          </div>
        </div>
        <div class="right-forum recommend-item">
          <div class="notice-title">
            <span class="title-text">
              <svg class="icon-font">
                <use xlink:href="#icon-wendang1"></use>
              </svg>
              热门问答
            </span>
            <a
              class="title-right"
              @click="$router.push('/forumCenter')"
              target="_blank"
            >
              查看更多
              <i class="el-icon-d-arrow-right"></i
            ></a>
          </div>
          <div class="notice-list">
            <div
              class="notice-list-item"
              v-for="item in aclhot"
              :key="item.id"
            >
              <a
                class="noticeitem-text"
                @click="$router.push(`/forumCenter/detail/${item.id}`)"
                :title="item.title"
                ><i class="iconfont icon-yuandian"></i>
                {{ item.title }}
              </a>
              <div class="noticeitem-time">
                {{ new Date(item.createTime).format() }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="hotCourse homemok">
      <div class="hotCourse-content center-content">
        <h3 class="mok-title">热门学习资源</h3>
        <div class="hotCourse-list clearfix">
          <div
            class="hotCourse-list-item fadeInUp"
            v-for="(item, key) in course"
            :key="item.id"
            :style="{ animationDelay: key * 0.1 + 's' }"
            @click="$router.push(`/courseDetail?id=${item.id}`)"
          >
            <div class="item-img">
              <img :src="item.cover" :alt="item.title" />
            </div>
            <h2 class="item-text ellipsis">
              {{ item.title }}
            </h2>
          </div>
        </div>
        <div class="mok-more">
          <a @click="$router.push('/know')">更多
            <i class="el-icon el-icon-d-arrow-right"></i>
          </a>
        </div>
        
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapActions } from 'vuex';
import axios from 'axios';
import Know from './know.vue';

// 简单的内存缓存实现
const cache = {
  data: {},
  set(key, data, ttl = 300000) { // 默认5分钟
    this.data[key] = {
      data,
      expiry: Date.now() + ttl
    };
  },
  get(key) {
    const item = this.data[key];
    if (!item) return null;
    
    if (Date.now() > item.expiry) {
      delete this.data[key];
      return null;
    }
    
    return item.data;
  }
};

export default {
  components: {
    Know 
  },
  data() {
    return {
       bannerHeight: 0, 
      showModal: false,
      affiche: [
      ],
      banner: [
        { id: 1, imageUrl: 'http://localhost:3000/api/pdfs/1.jpg', title: 'Image 1' },
        { id: 2, imageUrl: 'http://localhost:3000/api/pdfs/2.jpg', title: 'Image 2' },
        { id: 3, imageUrl: 'http://localhost:3000/api/pdfs/3.jpg', title: 'Image 3' },
        { id: 4, imageUrl: 'http://localhost:3000/api/pdfs/4.jpg', title: 'Image 4' },
        { id: 5, imageUrl: 'http://localhost:3000/api/pdfs/5.jpg', title: 'Image 5' },
      ],
      blogroll: [
        {
          id: 1,
          cover: 'http://localhost:3000/api/pdfs/logo1.png', 
          title: '湘潭大学计算机学院',
          desc: '湘潭大学计算机学院官方网站，提供学院最新动态、专业建设、师资力量、科研成果、招生就业等全方位信息服务',
          link: 'https://jwxy.xtu.edu.cn/index.htm'
        },
        {
          id: 2,
          cover: 'http://localhost:3000/api/pdfs/1602730889.jpg',
          title: '湖南省网安基地',
          desc: '湖南省网络安全人才培养和技术研发的重要基地，致力于网络安全教育培训、技术研究和产业孵化',
          link: 'http://www.cybersecbase.com/' 
        },
        {
          id: 3,
          cover: 'http://localhost:3000/api/pdfs/1.png', 
          title: '牛客',
          desc: '专业的IT求职学习平台，提供海量笔试面试题库、技术课程、企业真题和求职经验分享',
          link: 'https://www.nowcoder.com/'
        },
        {
          id: 4,
          cover: 'http://localhost:3000/api/pdfs/logo.png', 
          title: '黑马程序员',
          desc: '国内知名IT职业培训机构，提供Java、前端、Python、大数据等高质量实战课程',
          link: 'https://www.itheima.com/' 
        },
      ],
      course: [], // 初始化为空数组，将从后端获取
      aclhot: [], // 初始化为空数组，将从后端获取
    };
  },
  async mounted() {
    this.calculateBannerHeight();
    window.addEventListener('resize', this.calculateBannerHeight);
    await this.fetchHotCourses();
    await this.fetchHotArticles();
  },
  methods: {
    goLinkClick(link) {
      window.open(link, '_blank'); 
    },
      calculateBannerHeight() {
      // 获取banner容器的宽度（1200px或实际显示宽度）
      const bannerContainer = document.querySelector('.banner-container');
      if (!bannerContainer) return;
      
      const width = bannerContainer.clientWidth;
      // 根据16:9比例计算高度
      


      this.bannerHeight = (width * 9) / 16;
      
      // 设置最大高度限制
      if (this.bannerHeight > 500) {
        this.bannerHeight = 550;
      }
      
      // 设置最小高度限制
      if (this.bannerHeight < 250) {
        this.bannerHeight = 250;
      }
    },
    // 获取热门课程
    async fetchHotCourses() {
  const cacheKey = 'courses_list';
  const cachedData = cache.get(cacheKey);
  
  if (cachedData) {
    this.course = cachedData.slice(0, 5); // 取前5个
    return;
  }
  
  try {
    const response = await axios.get('/api/courses', {
      params: {
        page: 1,
        pageSize: 12, // 获取足够多的课程
        order: 1 // 按最新排序
      }
    });
    
    if (response.data.code === 20000) {
      const allCourses = response.data.data.courseList;
      // 取前5个课程，如果不足5个则取全部
      this.course = allCourses.slice(0, 5);
      // 缓存全部课程，有效期5分钟
      cache.set(cacheKey, allCourses, 300000);
    }
  } catch (error) {
    console.error('获取课程列表失败:', error);
    // 失败时使用备用数据
    this.course = [
      {
        id: 1,
        title: '网络安全知识基础',
        cover: require('@/assets/img/presentation.png'), 
      },
      {
        id: 2,
        title: '网络安全与法律',
        cover: require('@/assets/img/loyal.png'), 
      },
      {
        id: 3,
        title: '网络诈骗处理方法',
        cover: require('@/assets/img/123.png'), 
      },
      {
        id: 4,
        title: '网络安全入校园',
        cover: require('@/assets/img/567.png'), 
      },
      {
        id: 5,
        title: '网络安全入校园',
        cover: require('@/assets/img/567.png'), 
      },
    ];
  }
},

    
    // 获取热门文章
    async fetchHotArticles() {
      const cacheKey = 'hot_articles';
      const cachedData = cache.get(cacheKey);
      
      if (cachedData) {
        this.aclhot = cachedData;
        return;
      }
      
      try {
        const response = await axios.get('/api/forum/articles/hot');
        
        if (response.data.code === 20000) {
          this.aclhot = response.data.data.slice(0, 6); // 限制为6篇文章
          // 缓存结果，有效期5分钟
          cache.set(cacheKey, this.aclhot, 300000);
        }
      } catch (error) {
        console.error('获取热门文章失败:', error);
        // 失败时使用备用数据
        this.aclhot = [
          { id: 1, title: '这是一个热点', createTime: '2024-10-16' },
          { id: 2, title: '网络安全基础问题', createTime: '2024-10-15' },
          { id: 3, title: '如何防范网络钓鱼', createTime: '2024-10-14' },
          { id: 4, title: '数据加密技术解析', createTime: '2024-10-13' },
          { id: 5, title: '最新网络安全法规', createTime: '2024-10-12' },
          { id: 6, title: '企业安全防护策略', createTime: '2024-10-11' },
        ];
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.homemok {
  padding: 78px 0;
}
.mok-more {
  text-align: center;
  a {
    font-size: 22px;
    letter-spacing: 2px;
    font-weight: 600;
    .el-icon {
      display: inline-block;
      -webkit-transform: rotate(90deg);
      transform: rotate(90deg);
      font-weight: 400;
    }
  }
}
.mok-title {
  text-align: center;
  line-height: 37px;
  letter-spacing: 3px;
  font-weight: 600;
  font-size: 26px;
  padding-bottom: 52px;
}
.center-content {
  width: 1200px;
  margin: auto;
}
.banner {
  width: 100%;
  display: flex;
  justify-content: center;
  padding: 20px 0;
  background: #f8f9fa;
  
  .banner-container {
    width: 100%;
    margin: 0 auto;
    
    .banner-wrapper {
      position: relative;
      width: 100%;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
    
    .el-carousel {
      width: 100%;
      background: #fff;
      
      ::v-deep .el-carousel__container {
        height: 100%;
      }
      
      ::v-deep .el-carousel__arrow {
        background-color: rgba(255, 255, 255, 0.8);
        color: #333;
        z-index: 10;
        
        &:hover {
          background-color: rgba(255, 255, 255, 1);
        }
      }
    }
    
    .banner-image-container {
      position: relative;
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      overflow: hidden;
      
      // 背景模糊层
      .blur-background {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: 1;
        
        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          filter: blur(15px) brightness(1.1);
          transform: scale(1.05); // 稍微放大以避免边缘模糊不足
        }
      }
      
      // 前景清晰图片
      .foreground-image {
        position: relative;
        z-index: 2;
        max-width: 100%;
        max-height: 100%;
        
        img {
          width: 960px;
          height: 540px;
          max-width: 100%;
          max-height: 100%;
          object-fit: contain;
          display: block;
        }
      }
      
      // 渐变遮罩
      .gradient-overlay {
        position: absolute;
        top: 0;
        height: 100%;
        width: 200px;
        z-index: 3;
        pointer-events: none;
        
        &.left-overlay {
          left: 0;
          background: linear-gradient(90deg, 
            rgba(248, 249, 250, 0.95) 0%,
            rgba(248, 249, 250, 0.7) 50%,
            rgba(248, 249, 250, 0) 100%);
        }
        
        &.right-overlay {
          right: 0;
          background: linear-gradient(270deg, 
            rgba(248, 249, 250, 0.95) 0%,
            rgba(248, 249, 250, 0.7) 50%,
            rgba(248, 249, 250, 0) 100%);
        }
      }
    }
  }
}

// 响应式设计
@media (max-width: 1240px) {
  .banner .banner-container {
    width: 90%;
    max-width: 1200px;
    
    .banner-image-container {
      .gradient-overlay {
        width: 150px;
      }
    }
  }
}

@media (max-width: 1024px) {
  .banner {
    .banner-container {
      .banner-image-container {
        .gradient-overlay {
          width: 100px;
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .banner {
    padding: 15px 0;
    
    .banner-container {
      width: 95%;
      
      .banner-image-container {
        .gradient-overlay {
          width: 60px;
        }
      }
    }
  }
}

@media (max-width: 480px) {
  .banner {
    .banner-container {
      .banner-image-container {
        .gradient-overlay {
          width: 30px;
        }
      }
    }
  }
}

.fadeInUp {
  animation-name: fadeInUp;
  animation-duration: 0.5s;
  animation-delay: 0.3s;
}
@keyframes fadeInUp {
  from {
    opacity: 0;
    -webkit-transform: translate3d(0, 100%, 0);
    transform: translate3d(0, 100%, 0);
  }
  to {
    opacity: 1;
    -webkit-transform: translateZ(0);
    transform: translateZ(0);
  }
}
.friendlinks {
  .flink-content {
    width: 1230px;
    display: flex;
    .flink-item {
      cursor: pointer;
      width: 254px;
      height: 310px;
      border: 1px solid #eef3f7;
      background: #f5f7fa;
      -webkit-box-shadow: 4px 4px 33px 0 #dae8f0;
      box-shadow: 4px 4px 33px 0 #dae8f0;
      -webkit-transition: all 0.2s;
      transition: all 0.2s;
      padding-top: 40px;
      -webkit-box-sizing: border-box;
      box-sizing: border-box;
      // animation-delay: 0.3s;
      // animation-duration: 1s;
      .flink-item-img {
        text-align: center;
        img {
          width: 60%;
        }
      }
      .flink-item-text {
        width: 163px;
        padding-top: 27px;
        margin: auto;
        font-size: 16px;
        font-weight: 500;
        line-height: 22px;
        letter-spacing: 2px;
      }
      .flink-item-desc {
        width: 163px;
        padding-top: 27px;
        margin: auto;
        font-size: 16px;
        font-weight: 500;
        line-height: 22px;
        letter-spacing: 2px;
      }
      &:hover {
        background: #fff;
      }

      &:not(:first-child) {
        margin-left: 68px;
      }
    }
  }
}
.recommend {
  .recommend-content {
    margin: auto;
    display: flex;
    justify-content: space-around;
    .recommend-item {
      overflow: hidden;
      width: 567px;
      height: 360px;
      border-radius: 6px;
      background-color: #fff;
      padding: 24px 30px 30px 20px;
      -webkit-box-sizing: border-box;
      box-sizing: border-box;
      -webkit-transition: all 0.2s;
      transition: all 0.2s;
      &:hover {
        box-shadow: 4px 4px 5px 4px #e4e8eb;
      }
    }
    .left-notice,
    .right-forum {
      .notice-title {
        border-bottom: 1px solid #ebeef5;
        padding-bottom: 26px;
        .title-text {
          color: #303133;
          font-weight: 600;
          line-height: 22px;
          letter-spacing: 2px;
          font-size: 16px;
        }
        .title-right {
          color: #606266;
          font-size: 14px;
          letter-spacing: 1px;
          float: right;
        }
      }
      .notice-list {
        padding-top: 27px;
        .notice-list-item {
          display: flex;
          font-size: 14px;
          color: #606266;
          .noticeitem-text {
            text-align: left;
            line-height: 20px;
            letter-spacing: 1px;
            -webkit-box-flex: 1;
            -ms-flex: 1;
            flex: 1;
            -webkit-transition: all 0.2s;
            transition: all 0.2s;
            color: inherit;
            &:hover {
              color: $theme-color-font;
            }
          }
          .noticeitem-time {
            width: 110px;
            text-align: right;
          }
          &:not(:first-child) {
            margin-top: 14px;
          }
        }
      }
    }
  }
}
.hotCourse {
  .hotCourse-content {
    .hotCourse-list {
      .hotCourse-list-item {
        float: left;
        border-radius: 6px;
        overflow: hidden;
        box-shadow: rgba(215, 219, 221, 0.5) 8px 8px 18px 0;
        width: 226px;
        height: 224px;
        margin-bottom: 36px;
        -webkit-transition: all 0.2s;
        transition: all 0.2s;
        cursor: pointer;
        .item-img {
          width: 100%;
          height: 166px;
          border-radius: 0 0 6px 6px;
          overflow: hidden;
          img {
            width: 100%;
            height: 100%;
          }
        }
        .item-text {
          line-height: 58px;
          text-align: center;
          color: #303133;
          font-size: 14px;
          padding: 0 10px;
        }
        &:hover {
          transform: translateY(-5px);
          box-shadow: 10px 10px 16px 0 #b6cdd8;
        }
        &:not(:nth-child(5n)) {
          margin-right: 17px;
        }
      }
    }
  }
}
.modal {

display: flex;

position: fixed;

z-index: 1000;

left: 0;

top: 0;

width: 100%;

height: 100%;

overflow: auto;

background-color: rgba(0, 0, 0, 0.5); /* 背景半透明 */

}



.modal-content {
  position: fixed; 
  margin: auto; 
  width: 100%; 
  max-width: 1200px; 
  height: 100%; 
  max-height: 600px; 
  border: 1px solid #888; 
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2); 
  transform: translateY(-50%); 
  left: 50%;
  top: 40%;
  transform: translate(-50%, -50%); 
  box-sizing: border-box; /* 确保padding不影响整体尺寸 */
}.close {
 color: #aaa;
  font-size: 28px;
  font-weight: bold;
  z-index: 99;
  /* 调整为块级元素，单独占一行 */
  display: block;
  text-align: right; /* 保持右对齐但在单独一行 */
  margin-bottom: 15px; /* 与下方内容保持距离 */
  cursor: pointer;

}



.close:hover,

.close:focus {

color: black;

text-decoration: none;

cursor: pointer;

}

</style>
