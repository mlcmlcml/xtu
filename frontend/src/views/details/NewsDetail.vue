<template>
  <div class="news">
    <div class="center-content">
      <div class="newsDetail">
        <h2 class="title">
          {{ AfficheInfo.title }}
        </h2>
        <div class="time">
          <span>{{ AfficheInfo.author }}</span
          ><span>{{ AfficheInfo.createTime }}</span>
        </div>
        <div class="content" v-html="AfficheInfo.content"></div>
      </div>
    </div>
  </div>
</template>

<script>
// 注意：如果编译报 eduApi 找不到，请确保你在 script 顶部 import 了它
// import eduApi from '@/api/edu' 

export default {
  data() {
    return {
      AfficheInfo: {},
    };
  },
  created() {
    if (this.$route.params && this.$route.params.id) {
      this.getAfficheDetail(this.$route.params.id);
    }
  },
  methods: {
    getAfficheDetail(id) {
      // 假设 eduApi 已在全局或当前页面可用
      eduApi.getAffichebyId(id).then((res) => {
        if (res.code == 20000) {
          this.AfficheInfo = res.data.item;
        }
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.news {
  width: 100%;
  background-color: #f5f7fa;
  overflow: hidden;
  padding-bottom: 80px;
  padding: 21px 0;
  min-height: 800px;
}
.newsDetail {
  min-height: 800px;
  background: #fff;
  padding: 34px;
  text-align: center;
  .title {
    font-size: 24px;
    font-weight: 500;
    color: #303133;
    line-height: 24px;
    margin-bottom: 20px;
  }
  .time {
    color: #909399;
    font-size: 14px;
    margin-bottom: 40px;
    span {
      margin: 0 10px;
    }
  }
  .content {
    text-align: left;
    line-height: 1.8;
    font-size: 16px;
    color: #606266;
  }
}
</style>