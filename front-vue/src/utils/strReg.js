export function getCmdContextVar (str) {
  // 以字母下划线开头，后接数字下划线和字母
  var patt = new RegExp(/\{\{(_*[a-zA-Z]+[_a-zA-Z0-9]*)\}\}/g)
  var mm = {}
  var rets = str.match(patt)
  if (rets) {
    for (let i = 0; i < rets.length; i++) {
      const name = rets[i].split('{{')[1].split('}}')[0]
      if (!mm.hasOwnProperty(name)) {
        mm[name] = ''
      }
    }
  }
  return mm
}
