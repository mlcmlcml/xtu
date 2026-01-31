<template>
  <div>
    <div class="forum-box">
      <div class="forum-box-ul">
        <el-tabs v-model="forumActive" @tab-click="handleClick">
          <el-tab-pane
            :label="item.name"
            :name="item.id + ''"
            v-for="item in forumTablist"
            :key="item.id"
          >
            <div class="forum-tab-content">
              <div
                class="forum-content-item"
                v-for="item in AclList"
                :key="item.id"
              >
                <div class="forum-item-title">
                  <div class="community-sign" v-if="item.isTop">置顶</div>
                  <div class="community-sign-essence" v-if="item.isEss">
                    精华
                  </div>
                  <a
                    @click="$router.push(`/forumCenter/detail/${item.id}`)"
                    class="title-text ellipsis"
                    target="_blank"
                    >{{ item.title }}</a
                  >
                  <a
                    v-if="item.tagList && item.tagList.length"
                    v-for="tag in item.tagList"
                    :key="tag.id"
                    class="topic-t"
                  >
                    <span>#{{ tag.name }}#</span>
                  </a>
                </div>
                <div class="forum-item-text">
                  <!-- 显示内容预览 -->
                  <div v-html="getContentPreview(item.content)"></div>
                </div>
                <div class="forum-item-info clearfix">
                  <div class="forum-item-tag fl">
                    
                    <span class="icon-box">
                      <i class="el-icon el-icon-s-comment"></i>
                      {{ item.commentCount || 0 }}
                    </span>
                    <span class="icon-box">
                      <i class="el-icon el-icon-view"></i>
                      {{ item.viewCount || 0 }}
                    </span>
                  </div>
                  <div class="forum-text-date fr">
                    <span class="info-item">
                      {{ item.stuName }}
                      <span>发布于{{ formatDate(item.updateTime) }}</span>
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
        <div style="text-align: center; padding: 10px 0" v-if="total > 0">
          <el-pagination
            :current-page.sync="currentPage"
            :page-size="pageSize"
            background
            layout="total, prev, pager, next"
            :total="total"
            @current-change="handlePageChange"
          >
          </el-pagination>
        </div>
        <div v-else class="no-data">
          <p>暂无数据</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  data() {
    return {
      forumActive: "0",
      currentPage: 1,
      pageSize: 10,
      forumTablist: [],
      AclList: [],
      total: 0,
      // 添加缓存
      categoryCache: null,
      categoryCacheTime: 0
    };
  },
  created() {
    this.getAclCateList();
    this.getAclList();
  },
  methods: {
    handleClick() {
      this.currentPage = 1;
      this.getAclList();
    },
    
    handlePageChange(page) {
      this.getAclList();
    },
    
    formatDate(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`;
    },
    
    getContentPreview(content) {
      if (!content) return '';
      // 移除HTML标签，获取纯文本
      const plainText = content.replace(/<[^>]*>/g, '');
      // 截取前100个字符
      return plainText.length > 100 ? plainText.substring(0, 100) + '...' : plainText;
    },
    
    // 获取文章列表
    async getAclList() {
      try {
        this.$loading({ text: '加载中...' });
        
        const res = await axios.get('/api/forum/articles', {
          params: {
            page: this.currentPage,
            pageSize: this.pageSize,
            cateId: this.forumActive === "0" ? 0 : parseInt(this.forumActive)
          }
        });
        
        if (res.data.code === 20000) {
          this.AclList = res.data.data.items;
          this.total = res.data.data.total;
        }
      } catch (error) {
        console.error('获取文章列表失败:', error);
        this.$message.error('获取文章列表失败');
      } finally {
        this.$loading().close();
      }
    },
    
    // 获取分类列表
    async getAclCateList() {
      // 检查缓存（5分钟内有效）
      const now = Date.now();
      if (this.categoryCache && now - this.categoryCacheTime < 5 * 60 * 1000) {
        this.forumTablist = this.categoryCache;
        return;
      }
      
      try {
        const res = await axios.get('/api/forum/categories');
        if (res.data.code === 20000) {
          this.forumTablist = res.data.data;
          // 设置缓存
          this.categoryCache = res.data.data;
          this.categoryCacheTime = now;
        }
      } catch (error) {
        console.error('获取分类列表失败:', error);
        // 失败时使用默认分类
        this.forumTablist = [
          { name: "全部", id: 0 },
          { name: "讨论", id: 1 },
          { name: "问答", id: 2 }
        ];
      }
    }
  }
};
</script>


<style lang="scss" scoped>
.forum-box {
  padding: 20px 30px;
  background: #fff;
  border-radius: 10px;
  position: relative;
  .forum-box-ul {
    .forum-tab-content {
      .forum-content-item {
        width: 100%;
        border-bottom: 1px solid #ededed;
        padding: 22px 0;
        letter-spacing: normal;
        .forum-item-title {
          .community-sign {
            display: inline-block;
            width: 42px;
            height: 20px;
            line-height: 20px;
            font-size: 12px;
            font-weight: 400;
            color: #fff;
            border-radius: 10px;
            text-align: center;
            background: $theme-color-icon2;
            margin-right: 10px;
          }
          .community-sign-essence {
            display: inline-block;
            width: 42px;
            height: 20px;
            line-height: 20px;
            font-size: 12px;
            font-weight: 400;
            color: #fff;
            border-radius: 10px;
            text-align: center;
            background: $theme-color-icon;
            margin-right: 10px;
          }
          .title-text {
            font-size: 16px;
            color: #333;
            line-height: 21px;
            max-width: 450px;
            margin-right: 20px px;
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
            font-weight: 700;
          }
          .topic-t {
            margin-left: 12px;
            line-height: 20px;
            font-size: 14px;
            font-weight: 400;
            color: #ff9e3f;
            cursor: pointer;
          }
        }
        .forum-item-text {
          font-weight: 400;
          margin: 16px 0 26px;
          color: #888;
          line-height: 22px;
          font-size: 14px;
          text-overflow: -o-ellipsis-lastline;
          overflow: hidden;
          text-overflow: ellipsis;
          display: -webkit-box;
          -webkit-line-clamp: 2;
          line-clamp: 2;
          -webkit-box-orient: vertical;
        }
        .forum-item-info {
          .forum-item-tag {
            .icon-box {
              margin-right: 20px;
              color: #e2e2e2;
            }
          }
          .forum-text-date {
            line-height: 26px;
            color: #999;
            font-size: 12px;
          }
        }
      }
    }
  }
}
</style>
