import { createApp } from 'vue' // Vue 3 新写法
import App from './App.vue'
import router from './router'
import store from './store'

// 1. Element Plus 替换 Element UI
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css' // 默认样式
import './assets/css/element-variables.scss' // 你的自定义主题

// 2. 静态资源
import './assets/fonts/iconfont/iconfont.css'
import './assets/fonts/iconfont/iconfont.js'
import '@/utils/date.js'

// 3. Markdown 编辑器 (Vue 3 适配版写法)
import VMdPreview from '@kangc/v-md-editor/lib/preview';
import VueMarkdownEditor from '@kangc/v-md-editor';
import '@kangc/v-md-editor/lib/style/base-editor.css';
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js';
import '@kangc/v-md-editor/lib/theme/style/github.css';
import hljs from 'highlight.js';

VueMarkdownEditor.use(githubTheme, {
  Hljs: hljs,
});
VMdPreview.use(githubTheme, {
  Hljs: hljs,
});

// 4. 创建并挂载应用
const app = createApp(App)

app.use(store)       // 挂载 Vuex 4
app.use(router)      // 挂载 Router
app.use(ElementPlus) // 挂载 Element Plus
app.use(VueMarkdownEditor)
app.use(VMdPreview)

app.mount('#app')