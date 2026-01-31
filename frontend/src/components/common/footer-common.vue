<template>
  <div class="footer">
    <keep-alive>
      <div class="footer-ct">
        <div class="footer-list">
          <h4 class="footer-title">
            <span>友情链接</span>
          </h4>
          <ul class="footer-ul clearfix">
            <li v-for="link in friendLinks" :key="link.id">
              <a :href="link.url" :title="link.title" target="_blank">{{ link.name }}</a>
            </li>
          </ul>
        </div>
        <div class="b-foot">
          <section class="fl b-foot-link">
            <section class="b-f-link">
              <a href="#" title="关于我们" target="_blank">关于我们</a>|
              <a href="#" title="联系我们" target="_blank">联系我们</a>|
              <a href="#" title="帮助中心" target="_blank">帮助中心</a>
              <a href="#" title="资源下载" target="_blank">资源下载</a>| <br />
              <span>服务热线：0731-58292892(湘潭)</span><br />
              <span>Email：jsjxy@xtu.edu.cn</span>
            </section>
            <section class="b-f-link" style="margin-top: 5px">
              <span>©网站版权均归个人所有</span>
            </section>
          </section>
          <aside class="fr">
            <br />
            <span>湘潭大学 </span>
          </aside>
        </div>
      </div>
    </keep-alive>
  </div>
</template>

<script>
export default {
  name: "footer-common",
  data() {
    return {
      friendLinks: [
        { id: 1, name: '中央网信办', url: 'https://www.cac.gov.cn/', title: '国家互联网信息办公室' },
        { id: 2, name: '教育部', url: 'https://www.moe.gov.cn/', title: '中华人民共和国教育部' },
        { id: 3, name: '公安部', url: 'https://www.mps.gov.cn/', title: '中华人民共和国公安部' },
        { id: 4, name: '湘潭大学', url: 'https://www.xtu.edu.cn/', title: '湘潭大学官网' },
        { id: 5, name: '中国网安协会', url: 'https://www.nissei.org/', title: '中国网络安全协会' }
      ]
    };
  }
};
</script>

<style lang="scss" scoped>
.footer {
  background-color: #3d4553;
  border-top-width: 5px;
  border-top-style: solid;
  color: rgb(56, 95, 211);
  width: 100%;
  overflow: hidden;
  padding: 30px 0 10px;
}
.footer-ct {
  color: #999;
  width: 1200px;
  margin: 0 auto;
  position: relative;
  letter-spacing: 2px;
}
.footer-list {
  .footer-title {
    height: 30px;
    line-height: 30px;
    font-size: 18px;
  }
  .footer-ul {
    li {
      float: left;
      margin: 15px 30px 0 0;
      color: #666;
      font-size: 12px;
    }
  }
}
.b-foot {
  font-size: 14px;
  border-top: 1px solid #393939;
  padding: 20px 0 40px 0;
  margin-top: 30px;
  line-height: 25px;
  .b-foot-link {
    .b-f-link {
      color: #666;
      a {
        color: #666;
        display: block;
        margin: 0 10px 5px;
        display: inline-block;
        vertical-align: middle;
        &:first-child {
          margin-left: 0px;
        }
      }
    }
  }
  section {
    max-width: 1000px;
  }
  aside {
    img {
      width: 100px;
    }
  }
}
</style>
