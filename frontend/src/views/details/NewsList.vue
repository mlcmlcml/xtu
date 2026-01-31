<template>
  <div class="news">
    <div class="center-content">
      <header>
        <div class="header-search">
          <div class="search_left">
            <h2>新闻资讯</h2>
          </div>
        </div>
      </header>
      <div class="news-list clearfix">
        <None v-if="!newsList || !newsList.length" tips="空空如也" />
        <ul class="newsList_box">
          <li
            class="clearfix"
            v-for="item in newsList"
            :key="item.id"
            @click="$router.push(`/news/detail/${item.id}`)"
          >
            <div class="img_left">
              <img :src="item.cover" alt="" />
            </div>
            <div class="rightMes">
              <div class="topMes clearfix">
                <p class="title">{{ item.title }}</p>
                <p class="time">
                  <i class="el-icon el-icon-time"></i>
                  {{ item.createTime }}
                </p>
              </div>
              <div class="text" v-html="item.content"></div>
            </div>
          </li>
        </ul>
        <div style="text-align: center" v-if="total > 10">
          <el-pagination
            v-model:current-page="current"
            :page-size="10"
            layout="total, prev, pager, next"
            :total="total"
            @current-change="getnewsList"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      newsList: [],
      total: 0,
      current: 1
    }
  }
}
</script>

<style lang="scss" scoped>
.news-list {
  background: #fff;
  min-height: 500px;
}
.newsList_box {
  padding: 33px 30px 49px 16px;
  li {
    padding-bottom: 14px;
    border-bottom: 1px solid #f2f6fc;
    margin-bottom: 40px;
    cursor: pointer;
    .img_left {
      float: left;
      width: 224px;
      height: 126px;
      img { width: 100%; height: 100%; border-radius: 4px; }
    }
    .rightMes {
      float: right;
      width: 900px;
      .topMes {
        margin-bottom: 15px;
        .title { font-size: 16px; font-weight: 500; }
        .time { color: #909399; float: right; }
      }
      .text {
        color: #606266;
        font-size: 14px;
        line-height: 22px;
        height: 44px;
        overflow: hidden;
      }
    }
  }
}
</style>