import { createRouter, createWebHashHistory } from 'vue-router' // 改用 Vue Router 4 的方法

// 页面组件导入部分保持不变
const Index = () => import(/* webpackChunkName: "about" */ '../views/Index.vue')
const HomeCenter = () => import(/* webpackChunkName: "about" */ '../views/pages/HomeCenter.vue')
const CourseCenter = () => import(/* webpackChunkName: "about" */ '../views/pages/CourseCenter.vue')
const ForumCenter = () => import(/* webpackChunkName: "about" */ '../views/pages/ForumCenter.vue')
const LearnCenter = () => import(/* webpackChunkName: "about" */ '../views/pages/LearnCenter.vue')
const TeacherCenter = () => import(/* webpackChunkName: "about" */ '../views/pages/TeacherCenter.vue')
const know = () => import(/* webpackChunkName: "about" */ '../views/pages/know.vue')
const TeacherDetail = () => import(/* webpackChunkName: "about" */ '../views/details/TeacherDetail.vue')
const VideoDetail = () => import(/* webpackChunkName: "about" */ '../views/details/VideoDetail.vue')
const webrtc = () => import(/* webpackChunkName: "about" */ '../views/pages/webrtc.vue')
const CourseDetail = () => import(/* webpackChunkName: "about" */ '../views/details/CourseDetail.vue')
const PersonalCenter = () => import(/* webpackChunkName: "about" */ '../views/pages/PersonalCenter.vue')
const NoticeCenter = () => import(/* webpackChunkName: "about" */ '../views/pages/NoticeCenter.vue')
const ForumPage = () => import(/* webpackChunkName: "about" */ '../views/details/ForumPage.vue')
const ForumDetail = () => import(/* webpackChunkName: "about" */ '../views/details/ForumDetail.vue')
const AskRelease = () => import(/* webpackChunkName: "about" */ '../views/details/AskRelease.vue')
const Information = () => import(/* webpackChunkName: "personal" */ '../views/personal/Information.vue')
const MyClasses = () => import(/* webpackChunkName: "personal" */ '../views/personal/MyClasses.vue')
const Answer = () => import(/* webpackChunkName: "personal" */ '../views/personal/Answer.vue')
const Msg = () => import(/* webpackChunkName: "notice" */ '../views/notice/Msg.vue')
const Exam = () => import(/* webpackChunkName: "notice" */ '../views/notice/Exam.vue')
const Busywork = () => import(/* webpackChunkName: "notice" */ '../views/notice/Busywork.vue')
const NewsCenter = () => import(/* webpackChunkName: "news" */ '../views/pages/NewsCenter.vue')
const NewsList = () => import(/* webpackChunkName: "news" */ '../views/details/NewsList.vue')
const NewsDetail = () => import(/* webpackChunkName: "news" */ '../views/details/NewsDetail.vue')
const Home = () => import('../views/Home.vue')

// 路由配置表保持不变
export const routes = [
  {
    path: '/',
    name: 'Index',
    component: Index,
    redirect: "/homeCenter",
    children: [
      { path: '/homeCenter', name: 'HomeCenter', component: HomeCenter },
      { path: '/know', name: 'know', component: know },
      { path: '/courseCenter', name: 'CourseCenter', component: CourseCenter },
      {
        path: '/forumCenter',
        name: 'ForumCenter',
        component: ForumCenter,
        redirect: "/forumCenter/page",
        children: [
          { path: '/forumCenter/page', name: 'ForumPage', component: ForumPage },
          { path: '/forumCenter/detail/:id', name: 'ForumDetail', component: ForumDetail },
        ]
      },
      { path: '/learnCenter', name: 'LearnCenter', component: LearnCenter },
      { path: '/teacherCenter', name: 'TeacherCenter', component: TeacherCenter },
      { path: '/teacherDetail', name: 'TeacherDetail', component: TeacherDetail },
      { path: '/courseDetail', name: 'CourseDetail', component: CourseDetail },
      { path: '/webrtc', name: 'Webrtc', component: webrtc },
      { path: '/player/:vid', name: 'VideoDetail', component: VideoDetail },
      {
        path: '/personalCenter',
        name: 'PersonalCenter',
        redirect: "/personalCenter/information",
        component: PersonalCenter,
        children: [
          { path: '/personalCenter/information', name: 'Information', component: Information },
          { path: '/personalCenter/myClasses', name: 'MyClasses', component: MyClasses },
          { path: '/personalCenter/answer', name: 'Answer', component: Answer }
        ]
      },
      {
        path: '/notice',
        name: 'Notice',
        redirect: "/notice/msg",
        component: NoticeCenter,
        children: [
          { path: '/notice/msg', name: 'Msg', component: Msg },
          { path: '/notice/exam', name: 'exam', component: Exam },
          { path: '/notice/Busywork', name: 'Busywork', component: Busywork }
        ]
      },
      { path: '/home', name: 'Home', component: Home },
      { path: '/askrelease', name: 'AskRelease', component: AskRelease },
      {
        path: '/news',
        name: 'NewsCenter',
        component: NewsCenter,
        redirect: "/news/list",
        children: [
          { path: '/news/list', name: 'NewsList', component: NewsList },
          { path: '/news/detail/:id', name: 'NewsDetail', component: NewsDetail },
        ]
      }
    ]
  },
  {
    path: '/about',
    name: 'About',
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  },
  {
    path: '/do',
    name: 'ExamPaperDo',
    component: () => import('@/views/exam/paper/do')
  },
  {
    path: '/read',
    name: 'ExamPaperRead',
    component: () => import('@/views/exam/paper/read')
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHashHistory(), // 使用 hash 模式，对应原来的默认模式
  routes
})

// 导航守卫
router.beforeEach((to, from, next) => {
  window.scrollTo(0, 0) // 简化的滚动到顶部写法
  next()
})

export default router