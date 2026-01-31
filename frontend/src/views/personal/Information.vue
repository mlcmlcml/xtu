<template>
  <div>
    <div class="tab-title-box">
      <span class="tab-item">个人资料</span>
    </div>
    <div class="tab-content-box">
      <el-form
        class="personal-el-form"
        ref="form"
        label-width="80px"
        size="mini"
        :model="formUser"
      >
        <h2 class="personal-base-header">基本信息</h2>
        
        <el-form-item label="头像">
          <el-avatar :size="100" :src="formUser.userHead" fit="cover" />
        </el-form-item>

        <el-form-item label="昵称">
          <el-input
            v-model="formUser.nickName"
            :readonly="true"
          ></el-input>
        </el-form-item>

        <el-form-item label="邮箱">
          <el-input v-model="formUser.userEmail" :readonly="true" />
        </el-form-item>
        
        <el-form-item label="用户身份">
          <el-input v-model="formUser.userName" :readonly="true" />
        </el-form-item>
        
        <el-form-item label="学号">
          <el-input
            v-model="formUser.stuId"
            :disabled="true"
          />
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
import { mapGetters } from "vuex";
export default {
  data() {
    return {
      formUser: {
        stuId: "",
        userName: "",
        nickName: "",
        userHead: "",
        userEmail: "",
      },
      defaultAvatar: 'https://tse3-mm.cn.bing.net/th/id/OIP-C.qidgOqAsPEdzAg5inmSK3AAAAA?rs=1&pid=ImgDetMain'
    };
  },
  computed: {
    ...mapGetters('user', ['userInfo']),
    
    normalizedUserInfo() {
      console.log('userInfo in normalizedUserInfo:', this.userInfo);
      return {
        stuId: this.userInfo.stuId || "",
        userName: this.userInfo.userName || "",
        nickName: this.userInfo.nickName || "",
        userHead: this.userInfo.userHead || this.defaultAvatar,
        userEmail: this.userInfo.userEmail || ""
      };
    }
  },
  watch: {
    normalizedUserInfo: {
      immediate: true,
      deep: true,
      handler(newVal) {
        console.log('normalizedUserInfo changed:', newVal);
        this.formUser = {
          ...this.formUser,
          ...newVal
        };
      }
    }
  },
};
</script>

<style lang="scss" scoped>
.tab-title-box {
  display: flex;
  height: 51px;
  line-height: 51px;
  border-bottom: 1px solid #e2e2e2;
  background: #ebeef5;

  .tab-item {
    margin-left: 40px;
    color: $theme-color-font;
    line-height: 51px;
    height: 52px;
    min-width: 65px;
    padding: 0 15px;
    text-align: center;
    cursor: pointer;
    font-size: 18px;
    float: left;
    box-sizing: border-box;
    border-top: 2px solid $theme-color-font;
    background: #fff;
  }
}

.tab-content-box {
  background: #fff;
  padding: 20px 30px;
  
  .personal-base-header {
    font-size: 16px;
    font-weight: 600;
    padding-bottom: 15px;
    margin-bottom: 20px;
    border-bottom: 1px solid #eee;
  }
  
  .personal-el-form {
    max-width: 600px;
    
    .el-form-item {
      margin-bottom: 25px;
    }
  }
}
</style>