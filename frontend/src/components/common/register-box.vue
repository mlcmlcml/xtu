<template>
  <div>
    <el-dialog
      title="注册"
      :visible.sync="dialogFormRegister"
      @close="closeDialog"
      width="25%"
      center
    >
      <el-form
        :model="RegisterForm"
        status-icon
        ref="RegisterForm"
        :rules="RegisterRules"
        label-width="100px"
      >
        <el-form-item label="学号" prop="stuId">
          <el-input
            placeholder="请输入学号"
            v-model="RegisterForm.stuId"
          ></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input
            placeholder="请输入邮箱"
            v-model="RegisterForm.email"
          ></el-input>
        </el-form-item>

        <el-form-item label="设置密码" prop="password">
          <el-input
            placeholder="6~12个字符"
            type="password"
            v-model="RegisterForm.password"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="昵称" prop="nickName">
          <el-input
            placeholder="请输入昵称"
            v-model="RegisterForm.nickName"
          ></el-input>
        </el-form-item>
        <el-form-item label="头像" prop="userHead">
          <!-- 使用上传组件 -->
          <el-upload
            class="avatar-uploader"
            action="#"
            :show-file-list="false"
            :on-change="handleAvatarChange"
            :auto-upload="false"
            accept="image/*"
          >
            <img v-if="RegisterForm.userHead" :src="RegisterForm.userHead" class="avatar" />
            <i v-else class="el-icon-plus avatar-uploader-icon"></i>
          </el-upload>
          <input type="hidden" v-model="RegisterForm.userHead" />
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button
          style="width: 100%"
          type="primary"
          round
          @click.native="registerClick"
          :loading="loading"
        >确 定</el-button>
        <div class="button-wjmm">
          <a @click="gologinClick()">已激活，去登录</a>
        </div>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import axios from 'axios';
import { mapActions } from 'vuex';

export default {
  data() {
    return {
      dialogFormRegister: false,
      RegisterForm: {
        stuId: "",
        email: "",
        password: "",
        nickName: "",
        userHead: "https://tse3-mm.cn.bing.net/th/id/OIP-C.qidgOqAsPEdzAg5inmSK3AAAAA?rs=1&pid=ImgDetMain", // 默认头像
        avatarFile: null // 添加文件对象字段
      },
      loading: false,
      RegisterRules: {
        stuId: [
          { required: true, trigger: "blur", message: "请输入学号" },
          { min: 6, max: 20, message: "学号长度在6-20个字符", trigger: "blur" }
        ],
        email: [
          { required: true, trigger: "blur", message: "请输入邮箱" },
          { type: 'email', message: '请输入正确的邮箱地址', trigger: ['blur', 'change'] }
        ],
        password: [
          { required: true, trigger: "blur", message: "请输入密码" },
          { min: 6, max: 12, message: "密码长度在6-12个字符", trigger: "blur" }
        ],
        nickName: [
          { required: true, trigger: "blur", message: "请输入昵称" },
          { min: 2, max: 10, message: "昵称长度在2-10个字符", trigger: "blur" }
        ],
        userHead: [
          { required: true, trigger: "blur", message: "请上传头像" }
        ]
      },
    };
  },
  computed: {
    dialogregister() {
      return this.$store.state.dialog.register;
    },
  },
  watch: {
    dialogregister(val) {
      this.dialogFormRegister = val;
      if (!val) {
        this.$refs.RegisterForm.resetFields();
        // 重置为默认头像
        this.RegisterForm.userHead = "https://tse3-mm.cn.bing.net/th/id/OIP-C.qidgOqAsPEdzAg5inmSK3AAAAA?rs=1&pid=ImgDetMain";
        this.RegisterForm.avatarFile = null;
      }
    },
  },
  methods: {
    ...mapActions('user', ['login']),
    
    gologinClick() {
      this.$store.dispatch("dialog/setregister", false);
      this.$store.dispatch("dialog/setlogin", true);
    },
    
    closeDialog() {
      this.$store.dispatch("dialog/setregister", false);
    },
    
    // 处理头像选择
    handleAvatarChange(file) {
      if (file) {
        // 预览图片
        const reader = new FileReader();
        reader.onload = (e) => {
          this.RegisterForm.userHead = e.target.result;
        };
        reader.readAsDataURL(file.raw);
        // 保存文件对象
        this.RegisterForm.avatarFile = file.raw;
      }
      return false; // 阻止自动上传
    },
    
    registerClick() {
      this.$refs.RegisterForm.validate((valid) => {
        if (valid) {
          this.loading = true;
          
          // 使用FormData处理文件上传
          const formData = new FormData();
          formData.append('stuId', this.RegisterForm.stuId);
          formData.append('email', this.RegisterForm.email);
          formData.append('password', this.RegisterForm.password);
          formData.append('nickName', this.RegisterForm.nickName);
          
          // 如果有上传文件则添加
          if (this.RegisterForm.avatarFile) {
            formData.append('avatar', this.RegisterForm.avatarFile);
          }
          
          axios.post('/api/register', formData)
          .then((res) => {
          this.loading = false;
          // 修正提示信息为注册成功
          this.$message.success('注册成功');
        })
        .catch((error) => {
          this.loading = false;
          this.$message.error('注册失败：' + (error.response?.data?.message || error.message));
        })
          .finally(() => {
            this.loading = false;
          });
        } else {
          console.log("表单验证失败");
          return false;
        }
      });
    },
  },
};
</script>

<style scoped>
.button-wjmm {
  margin-top: 15px;
  text-align: center;
}

.button-wjmm a {
  color: #409EFF;
  cursor: pointer;
  text-decoration: none;
}

.button-wjmm a:hover {
  text-decoration: underline;
}

/* 头像上传样式 */
.avatar-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  width: 178px;
  height: 178px;
}

.avatar-uploader:hover {
  border-color: #409EFF;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  line-height: 178px;
  text-align: center;
}

.avatar {
  width: 178px;
  height: 178px;
  display: block;
}
</style>