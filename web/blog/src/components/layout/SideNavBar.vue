<template>
  <v-navigation-drawer
    app
    v-model="drawer"
    width="250"
    disable-resize-watcher
    right
    overlay-opacity="0.8"
  >
    <!-- 博主介绍 -->
    <div class="blogger-info">
      <v-avatar size="110" style="margin-bottom: 0.5rem">
        <img :src="blogInfo.avatar" />
      </v-avatar>
    </div>
    <!-- 博客信息 -->
    <div class="blog-info-wrapper">
      <div class="blog-info-data">
        <router-link to="/archives">
          <div style="font-size: 0.875rem">文章</div>
          <div style="font-size: 1.125rem">
            {{ blogInfo.articleCount }}
          </div>
        </router-link>
      </div>
      <div class="blog-info-data">
        <router-link to="/categories">
          <div style="font-size: 0.875rem">分类</div>
          <div style="font-size: 1.125rem">
            {{ blogInfo.categoryCount }}
          </div>
        </router-link>
      </div>
      <div class="blog-info-data">
        <router-link to="/tags">
          <div style="font-size: 0.875rem">标签</div>
          <div style="font-size: 1.125rem">
            {{ blogInfo.tagCount }}
          </div>
        </router-link>
      </div>
    </div>
    <hr />
    <!-- 页面导航 -->
    <div class="menu-container">
      <div class="menus-item">
        <router-link to="/">
          <i class="iconfont iconzhuye" /> 首页
        </router-link>
      </div>
      <div class="menus-item">
        <router-link to="/archives">
          <i class="iconfont iconguidang" /> 归档
        </router-link>
      </div>
      <div class="menus-item">
        <router-link to="/categories">
          <i class="iconfont iconfenlei" /> 分类
        </router-link>
      </div>
      <div class="menus-item">
        <router-link to="/tags">
          <i class="iconfont iconbiaoqian" /> 标签
        </router-link>
      </div>
      <!-- <div class="menus-item">
        <router-link to="/links">
          <i class="iconfont iconlianjie" /> 友链
        </router-link>
      </div> -->
      <div class="menus-item">
        <router-link to="/about">
          <i class="iconfont iconzhifeiji" /> 关于
        </router-link>
      </div>
      <!-- <div class="menus-item">
        <router-link to="/message">
          <i class="iconfont iconpinglunzu" /> 留言
        </router-link>
      </div> -->
      <div class="menus-item" v-if="!this.$store.state.avatar">
        <a @click="openLogin"><i class="iconfont icondenglu" /> 登录 </a>
      </div>
      <div v-else>
        <div class="menus-item">
          <router-link to="/user">
            <i class="iconfont icongerenzhongxin" /> 个人中心
          </router-link>
        </div>
        <div class="menus-item">
          <a @click="logout"><i class="iconfont icontuichu" /> 退出</a>
        </div>
      </div>
    </div>
  </v-navigation-drawer>
</template>

<style scoped>
.blogger-info {
  padding: 26px 30px 0;
  text-align: center;
}
.blog-info-wrapper {
  display: flex;
  align-items: center;
  padding: 12px 10px 0;
}
.blog-info-data {
  flex: 1;
  line-height: 2;
  text-align: center;
}
hr {
  border: 2px dashed #d2ebfd;
  margin: 20px 0;
}
.menu-container {
  padding: 0 10px 40px;
  animation: 0.8s ease 0s 1 normal none running sidebarItem;
}
.menus-item a {
  padding: 6px 30px;
  display: block;
  line-height: 2;
}
.menus-item i {
  margin-right: 2rem;
}
@keyframes sidebarItem {
  0% {
    transform: translateX(200px);
  }
  100% {
    transform: translateX(0);
  }
}
</style>

<script>
export default {
  created() {
    this.getBlogInfo();
  },
  computed: {
    drawer: {
      set(value) {
        this.$store.state.drawer = value;
      },
      get() {
        return this.$store.state.drawer;
      },
    },
    isLogin() {
      return this.$store.state.userId;
    },
  },
  data: function () {
    return {
      id: 1,
      blogInfo: {
        avatar: "",
        intro: "",
        articleCount: 0,
        categoryCount: 0,
        tagCount: 0,
      },
      queryParams: {
        pageSize: 5,
        pageNum: 1,
      },
    };
  },
  methods: {
    openLogin() {
      this.$store.state.loginFlag = true;
    },
    logout() {
      //如果在个人中心则跳回上一页
      if (this.$route.path == "/user") {
        this.$router.go(-1);
      }
      this.axios.get("logout").then(({ data }) => {
        if (data.status === 200) {
          this.$store.commit("logout");
          this.$toast({ type: "success", message: data.message });
        } else {
          this.$toast({ type: "error", message: data.message });
        }
      });
    },
    getBlogInfo() {
      this.axios.get(`profile/${this.id}`).then(({ data }) => {
        // console.log(data.data)
        this.blogInfo.avatar = data.data.avatar;
        this.blogInfo.nickname = data.data.name;
        this.blogInfo.intro = data.data.desc;
        this.$store.commit("checkBlogInfo", this.blogInfo);
      });
      this.axios.get(`tag`).then(({ data }) => {
        this.blogInfo.tagCount = data.total;
      });
      this.axios.get(`category`).then(({ data }) => {
        this.blogInfo.categoryCount = data.total;
      });
      this.axios
        .get("article", {
          params: {
            pagesize: this.queryParams.pageSize,
            pagenum: this.queryParams.pageNum,
          },
        })
        .then(({ data }) => {
          // console.log(data.data);
          this.blogInfo.articleCount = data.total;
        });
    },
  },
};
</script>
