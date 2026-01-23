import type { MenuMixedOptions } from './interface'

export const routeRecordRaw: MenuMixedOptions[] = [
  {
    path: 'articles',
    name: 'articleManagement',
    icon: 'iconify ph--article',
    label: '文章管理',
    redirect: 'articles/list',
    children: [
      {
        path: 'list',
        name: 'articleList',
        label: '文章列表',
        icon: 'iconify ph--list-bullets',
        meta: {
          componentName: 'ArticleList',
          showTab: true,
        },
        component: 'articles/index',
      },
      {
        path: 'edit/new',
        name: 'articleCreate',
        label: '新建文章',
        show: false,
        meta: {
          componentName: 'ArticleEdit',
          showTab: true,
          enableMultiTab: true,
          renderTabTitle() {
            return '新建文章'
          },
        },
        component: 'articles/edit',
      },
      {
        path: 'edit/:id',
        name: 'articleEdit',
        label: '编辑文章',
        show: false,
        meta: {
          componentName: 'ArticleEdit',
          showTab: true,
          enableMultiTab: true,
          renderTabTitle({ id }) {
            return `编辑文章${id ? `-${id}` : ''}`
          },
        },
        component: 'articles/edit',
      },
    ],
  },
  {
    path: 'notes',
    name: 'noteManagement',
    icon: 'iconify ph--note',
    label: '手记管理',
    redirect: 'notes/list',
    children: [
      {
        path: 'list',
        name: 'noteList',
        label: '手记列表',
        icon: 'iconify ph--note',
        meta: {
          componentName: 'NoteList',
          showTab: true,
        },
        component: 'notes/index',
      },
    ],
  },
  {
    path: 'pages',
    name: 'pageManagement',
    icon: 'iconify ph--layout',
    label: '页面管理',
    redirect: 'pages/list',
    children: [
      {
        path: 'list',
        name: 'pageList',
        label: '页面列表',
        icon: 'iconify ph--file-text',
        meta: {
          componentName: 'PageList',
          showTab: true,
        },
        component: 'pages/index',
      },
    ],
  },
  {
    path: 'albums',
    name: 'albumManagement',
    icon: 'iconify ph--image',
    label: '相册管理',
    redirect: 'albums/list',
    children: [
      {
        path: 'list',
        name: 'albumList',
        label: '相册列表',
        icon: 'iconify ph--image',
        meta: {
          componentName: 'AlbumList',
          showTab: true,
        },
        component: 'albums/index',
      },
    ],
  },
  {
    path: 'interactions',
    name: 'commentInteraction',
    icon: 'iconify ph--chat-circle-text',
    label: '评论与互动',
    meta: {
      componentName: 'CommentInteraction',
      showTab: true,
    },
    component: 'interactions/index',
  },
  {
    path: 'friend-links',
    name: 'friendLinkManagement',
    icon: 'iconify ph--link',
    label: '友链',
    redirect: 'friend-links/list',
    children: [
      {
        path: 'list',
        name: 'friendLinkList',
        label: '友链列表',
        icon: 'iconify ph--link',
        meta: {
          componentName: 'FriendLinkList',
          showTab: true,
        },
        component: 'friend-links/index',
      },
    ],
  },
  {
    path: 'union',
    name: 'unionManagement',
    icon: 'iconify ph--circles-three',
    label: '联合',
    redirect: 'union/list',
    children: [
      {
        path: 'list',
        name: 'unionList',
        label: '联合管理',
        icon: 'iconify ph--circles-three',
        meta: {
          componentName: 'UnionList',
          showTab: true,
        },
        component: 'union/index',
      },
    ],
  },
  {
    path: 'files',
    name: 'fileManagement',
    icon: 'icon-[fluent--cloud-arrow-up-24-regular]',
    label: '文件管理',
    redirect: 'files/list',
    children: [
      {
        path: 'list',
        name: 'fileList',
        label: '文件列表',
        icon: 'icon-[fluent--cloud-arrow-up-24-regular]',
        meta: {
          componentName: 'FileList',
          showTab: true,
        },
        component: 'uploads/index',
      },
    ],
  },
  {
    path: 'plugins',
    name: 'pluginManagement',
    icon: 'iconify ph--puzzle-piece',
    label: '插件与云函数',
    redirect: 'plugins/list',
    children: [
      {
        path: 'list',
        name: 'pluginList',
        label: '插件与云函数',
        icon: 'iconify ph--puzzle-piece',
        meta: {
          componentName: 'PluginList',
          showTab: true,
        },
        component: 'plugins/index',
      },
    ],
  },
  {
    path: 'advanced',
    name: 'advancedInfo',
    icon: 'iconify ph--info',
    label: '高级信息',
    redirect: 'advanced/overview',
    children: [
      {
        path: 'overview',
        name: 'advancedOverview',
        label: '高级信息',
        icon: 'iconify ph--info',
        meta: {
          componentName: 'AdvancedInfo',
          showTab: true,
        },
        component: 'advanced/index',
      },
    ],
  },
  {
    path: 'monitoring',
    name: 'systemMonitor',
    icon: 'iconify ph--activity',
    label: '系统监控',
    redirect: 'monitoring/overview',
    children: [
      {
        path: 'overview',
        name: 'systemMonitorOverview',
        label: '系统监控',
        icon: 'iconify ph--activity',
        meta: {
          componentName: 'SystemMonitor',
          showTab: true,
        },
        component: 'monitoring/index',
      },
    ],
  },
  {
    path: '/about',
    key: 'about',
    name: 'about',
    icon: 'iconify ph--info',
    label: '关于',
    component: 'about/index',
    meta: {
      showTab: true,
    },
  },
]
