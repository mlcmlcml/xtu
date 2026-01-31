<template>
  <div class="teacher">
    <div class="center-content">
      <header>
        <!-- 搜索部分 -->
        <div class="header-search">
          <div class="search_left">
            <h2>全部讲师</h2>
          </div>
          <div class="search_right">
            <div class="input_word_search">
              <el-input
                placeholder="输入你想搜索的讲师"
                class="input-with-select"
                size="small"
                v-model="searchValue"
                @keyup.enter.native="getTeacherList"
              >
                <el-button
                  slot="append"
                  icon="el-icon-search"
                  @click="getTeacherList"
                ></el-button>
              </el-input>
            </div>
          </div>
        </div>
      </header>
      <div class="teacher-list clearfix">
        <None v-if="!teacherList || !teacherList.length" tips="空空如也" />
        <div
          v-else
          class="teacher-list-item no-avatar"
          v-for="item in teacherList"
          :key="item.id"
          @click="$router.push(`/teacherDetail?id=${item.id}`)"
        >
          <div class="teacher-item-content">
            <div class="teacher-username">{{ item.teacherName }}</div>
            <div class="teacher-userinfo">{{ item.career }}</div>
            <div class="teacher-item-text">
              {{ item.intro }}
            </div>
          </div>
        </div>
      </div>
      <div style="text-align: center; margin-top: 20px;">
        <el-pagination
          :current-page.sync="current"
          :page-size="size"
          layout="total, prev, pager, next"
          :total="total"
          @current-change="getTeacherList"
          v-if="total > 0"
        >
        </el-pagination>
      </div>
    </div>
  </div>
</template>

<script>
import None from "@/components/common/no-databox.vue";
import axios from 'axios';
export default {
  data() {
    return {
      searchValue: "",
      current: 1,
      size: 8,
      total: 0,
      teacherList: [],
    };
  },
  created() {
    this.getTeacherList();
  },
  components: {
    None,
  },
  methods: {
    getTeacherList() {
      // 使用axios调用后端API获取教师列表
      axios.get(`/api/teachers`, {
        params: {
          page: this.current,
          pageSize: this.size,
          name: this.searchValue
        }
      })
      .then((response) => {
        if (response.data.code === 20000) {
          this.teacherList = response.data.data.rows;
          this.total = response.data.data.total;
        } else {
          this.$message.error(response.data.message || '获取教师列表失败');
        }
      })
      .catch((error) => {
        console.error('获取教师列表失败:', error);
        this.$message.error('获取教师列表失败，请稍后重试');
      });
    },
  },
};
</script>

<style lang="scss" >
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

.teacher .el-dialog,
.teacher .el-pager li {
  background: transparent !important;
  -webkit-box-sizing: border-box;
}
.teacher .el-pagination .btn-next,
.teacher .el-pagination .btn-prev {
  background: transparent !important;
}

.teacher {
  width: 100%;
  overflow: hidden;
  padding-bottom: 80px;
}

.header-search {
  display: flex;
  justify-content: space-between;
  -webkit-box-align: center;
  margin-bottom: 20px;
  -webkit-box-pack: justify;
  align-items: center;
  .search_left {
    flex: 1;
    height: 44px;
    line-height: 44px;
    overflow: hidden;
  }
}

.teacher-list {
  .teacher-list-item {
    cursor: pointer;
    width: 280px;
    height: 280px;
    border: 1px solid #eef3f7;
    background: #f5f7fa;
    -webkit-box-shadow: 4px 4px 33px 0 #dae8f0;
    box-shadow: 4px 4px 33px 0 #dae8f0;
    -webkit-transition: all 0.3s ease;
    transition: all 0.3s ease;
    padding: 30px 20px;
    -webkit-box-sizing: border-box;
    letter-spacing: 1px;
    box-sizing: border-box;
    float: left;
    margin: 0 30px 40px 0;
    text-align: center;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    border-radius: 8px;
    
    &.no-avatar {
      .teacher-username {
        font-size: 20px;
        font-weight: 600;
        margin-bottom: 12px;
        color: #333;
      }
      
      .teacher-userinfo {
        font-size: 14px;
        color: #666;
        margin-bottom: 25px;
        padding: 6px 16px;
        background: #eef3f7;
        border-radius: 18px;
        display: inline-block;
      }
      
      .teacher-item-text {
        width: 100%;
        padding-top: 20px;
        margin-top: 20px;
        border-top: 1px solid #e3e3e3;
        line-height: 1.6;
        font-size: 14px;
        color: #555;
        text-align: left;
        overflow: hidden;
        display: -webkit-box;
        -webkit-box-orient: vertical;
        -webkit-line-clamp: 4;
        max-height: 90px;
      }
    }
    
    &:hover {
      background: #fff;
      transform: translateY(-5px);
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.12);
    }
    
    &:nth-child(4n) {
      margin-right: 0;
    }
  }
}

@media (max-width: 1200px) {
  .teacher-list .teacher-list-item {
    width: calc(50% - 15px);
    
    &:nth-child(2n) {
      margin-right: 0;
    }
  }
}

@media (max-width: 768px) {
  .teacher-list .teacher-list-item {
    width: 100%;
    margin-right: 0;
  }
  
  .header-search {
    flex-direction: column;
    align-items: flex-start;
    
    .search_right {
      width: 100%;
      margin-top: 15px;
    }
  }
}
</style>