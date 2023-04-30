import {defineConfig} from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Memox Server",
  description: "Memox Open source Server",


  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      {text: 'Home', link: '/'},
      {text: 'API', link: '/API'}
    ],

    sidebar: [
      {
        text: 'guide',
        items: [
          {text: 'Get Started', link: 'https://github.com/MemoxApp/server#readme'},
          {text: 'API Document', link: '/API'}
        ]
      }
    ],

    socialLinks: [
      {icon: 'github', link: 'https://github.com/MemoxApp/server'}
    ]
  },


  
})
