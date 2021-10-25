import Vue from 'vue'
// import ElementUI from 'element-ui'
import {
    Tabs,
    TabPane,
    Form,
    FormItem,
    Button,
    Input,
    Table,
    TableColumn,
    MessageBox,
    Message,
    Popover,
    Upload,
    Link
} from 'element-ui'

import 'element-ui/lib/theme-chalk/index.css'
import App from './App.vue'

// Vue.use(ElementUI)
Vue.use(Tabs)
Vue.use(TabPane)
Vue.use(Form)
Vue.use(FormItem)
Vue.use(Button)
Vue.use(Input)
Vue.use(Table)
Vue.use(TableColumn)
// Vue.use(MessageBox)
// Vue.use(Message)
Vue.use(Popover)
Vue.use(Upload)
Vue.use(Link)

Vue.prototype.$msgbox = MessageBox;
Vue.prototype.$confirm = MessageBox.confirm
Vue.prototype.$message = Message

new Vue({
  el: '#app',
  render: h => h(App)
})
