<template>
  <div class="forumdetail">
    <!-- 面包屑导航 -->
    <header class="forum-header">
      <el-breadcrumb separator-class="el-icon-arrow-right">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item>{{ aclInfo.cateName || '未知分类' }}</el-breadcrumb-item>
        <el-breadcrumb-item>{{ aclInfo.title }}</el-breadcrumb-item>
      </el-breadcrumb>
    </header>

    <!-- 文章内容 -->
    <div class="forum-page">
      <h1>{{ aclInfo.title }}</h1>
      <div class="details-user-name-details">
        <div class="community-top">
          <img :src="aclInfo.stuHead" alt="作者头像" />
          <div class="name-box">
            <b>{{ aclInfo.stuName }}</b>
            <div class="name-info">
              <span>发布于{{ formatDate(aclInfo.createTime) }}</span>
              <span>浏览{{ aclInfo.viewCount }}次</span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Markdown内容渲染 -->
      <div class="detailCon">
        <v-md-editor 
          v-if="aclInfo.content" 
          :value="aclInfo.content" 
          mode="preview"
          class="article-content"
        ></v-md-editor>
        <p v-else>内容加载中...</p>
      </div>
    </div>

    <!-- 标签区域 -->
    <div class="forum-details-bottom">
      <div class="bottom-left">
        <span 
          v-for="tag in aclInfo.tagList" 
          :key="tag.id" 
          class="thanks-tips"
        >
          #{{ tag.name }}#
        </span>
      </div>
    </div>

    <!-- 评论区域 -->
    <div class="message-board">
      <div class="message-tit">
        <dt>
          全部评论<span>{{ commentList.length }}</span>
        </dt>
      </div>
      
      <div class="message-list" v-if="commentList.length">
        <div v-for="comment in commentList" :key="comment.id" class="comment-item">
          <div class="message-content-figure-tit">
            <img :src="comment.author_head" alt="头像" class="head">
            <span class="uname">{{ comment.author_name }}</span>
            <span>{{ formatDate(comment.create_time) }}</span>
            <div class="content">{{ comment.content }}</div>
            <div class="message-content-operation-btn">
              <span class="icon-box" @click="likeComment(comment.id)">
                <i class="el-icon-thumb"></i> 点赞({{ comment.like_count }})
              </span>
              <span class="icon-box" @click="replyComment(comment)">
                <i class="el-icon-chat-dot-round"></i> 回复
              </span>
            </div>
          </div>
        </div>
      </div>
      
      <div v-else class="no-comments">
        <p>暂无评论，成为第一个评论的人吧！</p>
      </div>
      
      <!-- 评论编辑器 -->
      <div class="comment-editor">
        <v-md-editor 
          v-model="newComment" 
          class="hide-toolbar"
          height="150px" 
          placeholder="请输入评论内容..."
          ref="commentEditor"
        ></v-md-editor>
        <el-button 
          type="primary" 
          @click="submitComment"
          :disabled="!newComment.trim()"
        >
          发表评论
        </el-button>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import axios from 'axios'
export default {
  data() {
    return {
      aclInfo: {
        tagList: [] // 确保tagList有默认值
      },
      commentList: [],
      newComment: ""
    };
  },
  computed: {
    ...mapState('user', ['userInfo']),
  },
  methods: {
    formatDate(dateString) {
      if (!dateString) return "";
      const date = new Date(dateString);
      return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`;
    },
    
   async getForumDetail(id) {
  try {
    const res = await axios.get(`/api/forum/articles/${id}`);
    if (res.data.code === 20000) {
      this.aclInfo = res.data.data.aclInfo;
      
      // 确保tagList存在
      if (!this.aclInfo.tagList) {
        this.aclInfo.tagList = [];
      }
      // 获取文章详情成功后，再获取评论
      this.getComments();
    }
  } catch (error) {
    console.error('获取文章详情失败:', error);
    this.$message.error('获取文章详情失败');
  }
},
    
    async submitComment() {
      if (!this.newComment.trim()) return;
      if (!this.userInfo.stuId) {
        this.$message.warning('请先登录');
        return;
      }
      
      try {
        const res = await axios.post('/api/forum/comments', {
          articleId: this.aclInfo.id,
          content: this.newComment,
          authorId: this.userInfo.id,
          authorName: this.userInfo.nickName,
          authorHead: this.userInfo.userHead
        });
        
        if (res.data.code === 20000) {
          this.$message.success('评论发表成功');
          
          const newCommentData = {
            id: res.data.data.id,
            article_id: this.aclInfo.id,
            content: this.newComment,
            author_id: this.userInfo.id,
            author_name: this.userInfo.nickName,
            author_head: this.userInfo.userHead,
            create_time: new Date().toISOString(),
            like_count: 0
          };
          
          // 将新评论添加到列表开头
          this.commentList.unshift(newCommentData);
          this.newComment = "";
        }
      } catch (error) {
        console.error('发表评论失败:', error);
        this.$message.error('发表评论失败');
      }
    },
    
    async getComments() {
      try {
        const res = await axios.get(`/api/forum/comments?articleId=${this.aclInfo.id}`);
        if (res.data.code === 20000) {
          this.commentList = res.data.data;
        }
      } catch (error) {
        console.error('获取评论失败:', error);
      }
    },
    
    async likeComment(commentId) {
      try {
        await axios.post(`/api/forum/comments/${commentId}/like`);
        this.$message.success('点赞成功');
        
        // 更新点赞数，无需重新获取全部评论
        const comment = this.commentList.find(c => c.id === commentId);
        if (comment) {
          comment.like_count = (comment.like_count || 0) + 1;
        }
      } catch (error) {
        console.error('点赞失败:', error);
        this.$message.error('点赞失败');
      }
    },
  
    replyComment(comment) {
      // 实现回复功能
      this.newComment = `@${comment.author_name} `;
      this.$refs.commentEditor.focus();
    }
  },
  created() {
    const articleId = this.$route.params.id;
    if (articleId) {
      this.getForumDetail(articleId);
    }
  }
};
</script>


<style lang="scss" scoped>

.hide-toolbar ::v-deep(.v-md-editor__toolbar) {
  display: none;
}


.malr {
  // margin-left: 58px;
  // margin-right: 62px;

  margin-left: 20px;
  margin-right: 20px;
}
.forum-page {
  background: rgb(255, 255, 255);
  padding: 30px;
  h1 {
    font-size: 22px;
    margin-bottom: 30px;
    font-weight: 700;
    line-height: 1;
    color: #333;
  }
  .details-user-name-details {
    width: 100%;
    height: 50px;
    margin-bottom: 30px;
    .community-top {
      display: flex;
      img {
        width: 50px px;
        height: 50px;
        border-radius: 50%;
        margin-right: 10px;
      }
      .name-box {
        padding: 6px px 0 8px;
        width: 702px;
        border-bottom: 1px solid #f6f6f6;
        b {
          display: block;
          margin-bottom: 5px;
          font-size: 16px;
          font-weight: 700;
          color: #333;
          line-height: 16px;
        }
        .name-info span {
          font-size: 14px;
          font-weight: 400;
          color: #c3c3c3;
          line-height: 14px;
          margin-right: 20px;
        }
      }
    }
  }

  // .myEditorTxt {
  //   font-size: 15px;
  //   line-height: 1.8;
  //   img {
  //     width: 100%;
  //   }
  // }
  .detailCon {
    width: 702px;
    margin: 0 auto;
    font-size: 15px;
    line-height: 1.8;
    padding-bottom: 30px;
  }
}

.github-markdown-body {
  padding: 0px !important;
}
.v-md-editor {
  background: transparent !important;
}
.forum-header {
  padding-bottom: 16px;
}
.content-text {
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: monospace;
  background: #f8f8f8;
  padding: 15px;
  border-radius: 4px;
  line-height: 1.6;
  max-height: 500px;
  overflow: auto;
}
.forum-details-bottom {
  padding: 0 34px;
  height: 71px;
  border-top: 1px solid #f6f6f6;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  border-radius: 0 0 8px 8px;
  margin-left: 0;

  .thanks-tips {
    line-height: 20px;
    font-size: 14px;
    font-weight: 400;
    color: #ff9e3f;
    cursor: pointer;
    margin-right: 5px;
  }
  .toolBar {
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    -webkit-box-pack: end;
    -ms-flex-pack: end;
    justify-content: flex-end;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    li {
      text-align: center;
      cursor: pointer;
      margin-right: 20px;
      font-size: 14px;
      span {
        color: #dadada;
        font-size: 16px;
        i {
          font-size: 18px;
          vertical-align: middle;
          margin-right: 5px;
        }
      }
    }
  }
}
.message-board {
  width: 100%;
  min-height: 300px;
  background: #fff;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  padding: 0 30px;
  margin-top: 20px;
  border-radius: 8px;
  overflow: hidden;
  .message-tit {
    width: 100%;
    height: 50px;
    line-height: 50px;
    margin-top: 20px;
  }
  .message-list {
    .comment-item {
      box-sizing: border-box;
      padding: 20px 0 30px 0;
      border-bottom: 1px solid #efefef;
      .message-content-reply-box .message-content-figure-tit {
        .head {
          width: 30px;
          height: 30px;
        }
      }
      .message-content-figure-tit {
        width: 100%;
        // margin-bottom: 20px;
        .head {
          width: 50px;
          height: 50px;
          float: left;
          border-radius: 50%;
          margin-right: 15px;
        }

        .uname {
          color: #333;
          font-size: 15px;
          font-weight: 700;
          display: inline-block;
          padding-top: 8px;
          margin-right: 40px;
        }
        span {
          font-size: 15px;
          font-weight: 400;
          display: block;
          margin-bottom: 16px;
          color: #c3c3c3;
        }
        em {
          color: #333;
          font-size: 15px;
          font-weight: 700;
          display: inline-block;
          padding-top: 8px;
          font-style: normal;
        }
        b {
          margin-left: 40px;
          font-weight: 400;
          color: #c3c3c3;
          font-size: 13px;
        }
        .textCon {
          line-height: 20px;
          margin-bottom: 10px;
          letter-spacing: 0.5px;
          font-family: Microsoft YaHei;
          color: #666;
        }
        .content {
          font-size: 14px;
          font-weight: 400;
          color: #333;
          margin-bottom: 0;
          padding-left: 57px;
          line-height: 20px;
        }
      }
      .message-content-operation-btn {
        width: calc(100% - 58px);
        margin-bottom: 10px;
        span {
          font-weight: 400;
          color: #c3c3c3;
          cursor: pointer;
          font-size: 13px;
          display: inline-block;
          margin-right: 15px;
        }
        .icon-box.active {
          color: #1058fa !important;
        }
      }
      .message-content-reply-box {
        margin-left: 58px;
        margin-right: 62px;
        -webkit-box-sizing: border-box;
        box-sizing: border-box;
        padding: 0 0;
        position: relative;
        overflow: hidden;
        margin-top: 20px;
        background-color: #f6f6f6;
        border-radius: 5px;
      }
    }
  }
}
// 添加新样式
.detailCon {
  margin-top: 20px;
  padding: 20px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.article-content {
  background: transparent !important;
  padding: 0 !important;
}

.comment-editor {
  margin-top: 30px;
  
  .el-button {
    margin-top: 15px;
    float: right;
  }
}

.no-comments {
  text-align: center;
  padding: 40px 0;
  color: #999;
}

.name-info {
  display: flex;
  gap: 15px;
  margin-top: 5px;
  font-size: 14px;
  color: #888;
}

// 其他样式保持不变...
</style>