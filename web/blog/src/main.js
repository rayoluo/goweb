import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import vuetify from "./plugins/vuetify";
import animated from "animate.css";
import "./assets/css/index.css";
import "./assets/css/iconfont.css";
import "./assets/css/markdown.css";
import "./assets/css/iconfont_csdn.css";
// import "./assets/css/markdown_new.css";
import config from "./assets/js/config";
import Share from "vue-social-share";
import "vue-social-share/dist/client.css";
import { vueBaberrage } from "vue-baberrage";
import axios from "axios";
import VueAxios from "vue-axios";
import moment from "moment";
import InfiniteLoading from "vue-infinite-loading";
import "highlight.js/styles/atom-one-dark.css";
import VueImageSwipe from "vue-image-swipe";
import "vue-image-swipe/dist/vue-image-swipe.css";
import Toast from "./components/toast/index";
import NProgress from "nprogress";
import "nprogress/nprogress.css";

Vue.prototype.config = config;
Vue.config.productionTip = false;
Vue.use(animated);
Vue.use(Share);
Vue.use(vueBaberrage);
Vue.use(InfiniteLoading);
Vue.use(VueAxios, axios);
Vue.use(VueImageSwipe);
Vue.use(Toast);

Vue.filter("date", function(value) {
  return moment(value).format("YYYY-MM-DD");
});

Vue.filter("hour", function(value) {
  return moment(value).format("HH:mm:ss");
});

Vue.filter("num", function(value) {
  if (value >= 1000) {
    return (value / 1000).toFixed(1) + "k";
  }
  return value;
});

router.beforeEach((to, from, next) => {
  NProgress.start();
  if (to.meta.title) {
    document.title = to.meta.title;
  }
  next();
});

router.afterEach(() => {
  window.scrollTo({
    top: 0,
    behavior: "instant"
  });
  NProgress.done();
});

axios.interceptors.response.use(
  function(response) {
    switch (response.data.code) {
      case 50000:
        Vue.prototype.$toast({ type: "error", message: "系统异常" });
    }
    return response;
  },
  function(error) {
    return Promise.reject(error);
  }
);

let Url = "https://rayoluo.top/api/v1/";
// let Url = 'http://localhost:3000/api/v1/'

axios.defaults.baseURL = Url;

// 设置withcredentials
axios.defaults.withCredentials = true;

new Vue({
  router,
  store,
  vuetify,
  render: h => h(App)
}).$mount("#app");
