<template>
  <div class="answer" style="padding: 20px; position: relative">
    <el-button
      @click="$router.push('/askrelease')"
      type="primary"
      plain
      class="answer-button"
    >
      发布内容
    </el-button>
    <el-tabs v-model="answerActive">
      <el-tab-pane label="我的帖子" name="first">
        <None v-if="!aclList || aclList.length < 1" tips="空空如也" />
        <div class="answer-content-box">
          <ul class="lists-box">
            <li
              class="lists-box-item clearfix"
              v-for="item in aclList"
              :key="item.id"
            >
              <a
                @click="$router.push(`/forumCenter/detail/${item.id}`)"
                class="ques-list-title fl"
                target="_blank"
              >
                {{ item.title }}
                <span v-for="tag in item.tagList" :key="tag.id" class="topic-t">
                  #{{ tag.name }}#
                </span>
              </a>
              <div class="ques-list-tool fr">
                <div v-html="item.content"></div>
                <p class="ques-list-time">{{ item.createTime }}</p>
                <el-button type="danger" size="small" @click="deleteAcl(item.id)">删除</el-button>
              </div>
            </li>
          </ul>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
export default {
  data() {
    return {
      answerActive: 'first',
      aclList: []
    }
  }
}
</script>

<style lang="scss">
.answer .el-tabs__item { font-size: 16px; }
.answer-button { position: absolute; right: 20px; z-index: 1; }
.lists-box-item {
  padding: 20px 0;
  border-bottom: 1px solid #eee;
  .ques-list-title { font-size: 16px; color: #333; }
  .topic-t { color: #409eff; margin-left: 10px; }
}
</style>