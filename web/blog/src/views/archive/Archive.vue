<template>
  <div>
    <!-- banner -->
    <div class="archive-banner banner">
      <h1 class="banner-title">归档</h1>
    </div>
    <!-- 归档列表 -->
    <v-card class="blog-container">
      <timeline>
        <timeline-title> 目前共计{{ total }}篇文章，继续加油 </timeline-title>
        <timeline-item v-for="item of artList" :key="item.ID">
          <!-- 日期 -->
          <span class="time">{{ item.CreatedAt | date }}</span>
          <!-- 文章标题 -->
          <router-link
            :to="'/articles/' + item.ID"
            style="color: #666; text-decoration: none"
          >
            {{ item.title }}
          </router-link>
        </timeline-item>
      </timeline>
      <!-- 分页按钮 -->
      <v-pagination
        color="#00C4B6"
        total-visible="7"
        v-model="queryParam.pagenum"
        :length="Math.ceil(total / queryParam.pagesize)"
        @input="getArtList()"
      ></v-pagination>
    </v-card>
  </div>
</template>

<script>
import { Timeline, TimelineItem, TimelineTitle } from "vue-cute-timeline";
export default {
  mounted() {
    this.getArtList();
  },
  components: {
    Timeline,
    TimelineItem,
    TimelineTitle,
  },
  data: function () {
    return {
      artList: [],
      queryParam: {
        pagesize: 10,
        pagenum: 1,
      },
      total: 0,
    };
  },
  methods: {
    // 获取文章列表
    async getArtList() {
      const { data: res } = await this.$http.get("article", {
        params: {
          pagesize: this.queryParam.pagesize,
          pagenum: this.queryParam.pagenum,
        },
      });
      // console.log(res)
      this.artList = res.data;
      this.total = res.total;
    },
  },
};
</script>

<style scoped>
.archive-banner {
  background: url(../../assets/img/archives.jpg)
    center center / cover no-repeat;
  background-color: #49b1f5;
}
.time {
  font-size: 0.75rem;
  color: #555;
  margin-right: 1rem;
}
</style>
