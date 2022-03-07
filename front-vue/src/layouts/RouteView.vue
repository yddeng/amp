<script>
export default {
  name: 'RouteView',
  props: {
    keepAlive: {
      type: Boolean,
      default: true
    }
  },
  data () {
    return {}
  },
  render () {
    const inKeep = (
      <keep-alive>
        <router-view />
      </keep-alive>
    )
    // const notKeep = (
    //   <router-view />
    // )
    return inKeep
  }
}
</script>
