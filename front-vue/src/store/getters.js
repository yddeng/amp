const getters = {
  isMobile: state => state.app.isMobile,
  lang: state => state.app.lang,
  theme: state => state.app.theme,
  color: state => state.app.color,
  token: state => state.user.token,
  nickname: state => state.user.name,
  multiTab: state => state.app.multiTab
}

export default getters
