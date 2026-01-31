
<template>
  <div class="videodetail">
    <div class="center-content">
      <header>
        <el-breadcrumb separator-class="el-icon-arrow-right">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/courses' }">课程</el-breadcrumb-item>
          <el-breadcrumb-item>视频详情</el-breadcrumb-item>
          <el-breadcrumb-item>{{ video ? video.title : '加载中...' }}</el-breadcrumb-item>
        </el-breadcrumb>
      </header>

      <div v-loading="loading" class="play-page">
        <div v-if="video" class="video-container">
          <div id="dplayer" class="dplayer-container"></div>
          
          <div class="video-info">
            <h2>{{ video.title }}</h2>
            <p>{{ video.description || '暂无描述' }}</p>
            <div class="meta">
              <span>时长: {{ formatDuration(video.duration) }}</span>
              <span>上传时间: {{ formatDate(video.createdAt) }}</span>
            </div>
          </div>
          
          <div class="video-controls">
            <el-button @click="toggleNotes" type="primary" class="notes-button">
              {{ showNotes ? '隐藏笔记' : '显示笔记' }}
            </el-button>
            <el-button @click="downloadVideo" type="success" class="download-button">
              <i class="el-icon-download"></i> 下载视频
            </el-button>
          </div>
        </div>
        
        <div v-if="!video" class="empty-state">
          <el-empty description="视频加载失败">
            <el-button type="primary" @click="retry">重新加载</el-button>
          </el-empty>
        </div>
        
        <div v-if="showNotes && video" class="notes-section">
          <h3>学习笔记</h3>
          <el-input
            type="textarea"
            :rows="6"
            placeholder="记录学习笔记..."
            v-model="noteContent"
          ></el-input>
          <div class="notes-actions">
            <el-button @click="saveNote" type="primary">保存笔记</el-button>
            <el-button @click="clearNote">清空</el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import DPlayer from 'dplayer';
import html2canvas from 'html2canvas';
import axios from 'axios'; // 引入axios

export default {
  data() {
    return {
      video: null,
      loading: false,
      showNotes: false,
      noteContent: '',
      player: null
    };
  },
  mounted() {
    this.loadVideo();
  },
  beforeDestroy() {
    if (this.player) {
      this.player.destroy();
    }
  },
  methods: {
    async loadVideo() {
      try {
        this.loading = true;
        const videoId = this.$route.params.vid;
        
        if (!videoId) {
          this.$message.error('视频ID无效');
          return;
        }
        
        const response = await axios.get(`/api/videos/${videoId}`);
        
        if (response.data.code === 20000) {
          this.video = response.data.data.video;
          
          // 确保URL是相对路径
          const url = new URL(this.video.url);
          this.video.url = url.pathname;
          
          this.$nextTick(() => {
            this.initPlayer();
          });
        } else {
          this.$message.error('获取视频信息失败: ' + response.data.message);
        }
      } catch (error) {
        console.error('获取视频信息出错:', error);
        this.$message.error('视频加载失败');
      } finally {
        this.loading = false;
      }
    },
    
    initPlayer() {
      // 销毁旧的播放器实例
      if (this.player) {
        this.player.destroy();
      }
      
      // 确保容器存在
      const dplayerElement = document.getElementById('dplayer');
      if (!dplayerElement) {
        console.error("DPlayer容器未找到");
        return;
      }
      
      // 清空容器
      dplayerElement.innerHTML = '';
      
      // 初始化新播放器
      this.player = new DPlayer({
        container: dplayerElement,
        autoplay: true,
        theme: '#FADFA3',
        video: {
          url: this.video.url,
          type: 'auto'
        },
        contextmenu: [
          {
            text: '关于播放器',
            link: 'https://github.com/DIYgod/DPlayer'
          }
        ]
      });
      
      // 监听错误事件
      this.player.on('error', () => {
      });
    },
    
    toggleNotes() {
      this.showNotes = !this.showNotes;
    },
    
    saveNote() {
      if (this.noteContent.trim()) {
        // 实际应用中应调用API保存笔记
        this.$message.success('笔记已保存');
      } else {
        this.$message.warning('笔记内容不能为空');
      }
    },
    
    clearNote() {
      this.noteContent = '';
    },
    
    downloadVideo() {
      if (!this.video || !this.video.url) {
        this.$message.warning('无法下载视频');
        return;
      }
      
      // 创建下载链接
      const link = document.createElement('a');
      link.href = this.video.url;
      link.download = `${this.video.title}.mp4`;
      link.click();
    },
    
    formatDuration(seconds) {
      if (!seconds) return '00:00';
      const mins = Math.floor(seconds / 60);
      const secs = Math.floor(seconds % 60);
      return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
    },
    
    formatDate(dateString) {
      if (!dateString) return '';
      const date = new Date(dateString);
      return date.toLocaleDateString();
    },
    
    retry() {
      this.loadVideo();
    }
  }
};
</script>

<style lang="scss" scoped>
.videodetail {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  
  .center-content {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    padding: 20px;
  }
  
  .play-page {
    margin-top: 20px;
    
    .video-container {
      position: relative;
      background: #000;
      border-radius: 8px;
      overflow: hidden;
      
      .dplayer-container {
        width: 100%;
        height: 500px;
        
        @media (max-width: 768px) {
          height: 300px;
        }
      }
      
      .video-info {
        padding: 15px;
        background: #f5f7fa;
        border-top: 1px solid #ebeef5;
        
        h2 {
          margin: 0 0 10px;
          color: #303133;
        }
        
        p {
          margin: 0 0 10px;
          color: #606266;
        }
        
        .meta {
          display: flex;
          gap: 15px;
          font-size: 13px;
          color: #909399;
        }
      }
      
      .video-controls {
        padding: 10px 15px;
        background: #f5f7fa;
        border-top: 1px solid #ebeef5;
        display: flex;
        justify-content: flex-end;
        gap: 10px;
      }
    }
    
    .empty-state {
      margin-top: 50px;
      text-align: center;
    }
    
    .notes-section {
      margin-top: 20px;
      padding: 15px;
      border: 1px solid #ebeef5;
      border-radius: 4px;
      background: #f8f9fa;
      
      h3 {
        margin-top: 0;
        margin-bottom: 15px;
        color: #303133;
      }
      
      .notes-actions {
        margin-top: 15px;
        text-align: right;
      }
    }
  }
}
</style>