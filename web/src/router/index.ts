import { createRouter, createWebHistory } from 'vue-router'
import taskView from '../views/task/task.vue'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/task',
      children: [
        {
          path: '/',
          name: 'task-index',
          component: taskView
        },
        {
          path: 'cron',
          name: 'task-cron',
          component: () => import('../views/task/cron.vue')
        },
        {
          path: 'memory',
          name: 'task-memory',
          component: () => import('../views/task/memory.vue')
        },
      ]
    },
    {
      path: '/rpa',
      children: [
        {
          path: 'index',
          name: 'rpa-index',
          component: () => import('../views/rpa/rpa/index.vue')
        },
        {
          path: 'var',
          name: 'rpa-var',
          component: () => import('../views/rpa/var/index.vue')
        },
        {
          path: 'group',
          name: 'rpa-group',
          component: () => import('../views/rpa/group/index.vue')
        },
      ]
    },
    { 
      path:"/batch",
      name:'batch',
      component: () => import('../views/batch/index.vue')
    },
    { 
      path:"/log",
      name:'log',
      component: () => import('../views/log/index.vue')
    },
    { 
      path:"/wecom",
      name:'wecom',
      component: () => import('../views/wecom/index.vue')
    },
    { 
      path:"/file",
      name:'file',
      component: () => import('../views/file/index.vue')
    },
  ]
})

export default router
