import {defineConfig} from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
    title: "Memox Server",
    description: "Memox Open source Server",


    themeConfig: {
        // https://vitepress.dev/reference/default-theme-config
        nav: [
            {text: '首页', link: '/'},
            {text: '快速开始', link: '/GET_STARTED'},
            {text: 'API', link: '/API'},
        ],

        sidebar: [
            {
                text: 'guide',
                items: [
                    {text: '快速开始', link: '/GET_STARTED'},
                    {text: '配置说明', link: '/CONFIG'},
                    {text: 'API 文档', link: '/API'},
                    {text: '已知问题', link: '/ISSUE'},
                    {text: '常见问题', link: '/TROUBLE_SHOOTING'}
                ]
            }
        ],

        socialLinks: [
            {icon: 'github', link: 'https://github.com/MemoxApp/server'}
        ]
    },


})
