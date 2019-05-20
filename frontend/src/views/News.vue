<template>
   <div class="content events">
      <h2>News</h2>
      <div class="news-item" v-for="item in news" :key="item.id">
         <h3>{{item.title}}<span class="date">{{formatDate(item.createdAt)}}</span></h3>
         <div class="text" v-html="item.content"></div>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
export default {
   name: "news",
   computed: {
      ...mapState({
         news: state => state.public.news,
      })
   },
   created: function () {
      this.$store.dispatch('public/getNews')
   },
   methods: {
      formatDate(date) {
         return date.split("T")[0]
      },
   }
};
</script>

<style scoped>
.news-item {
   margin: 10px 0 25px 0;
   padding: 0px;
}
h3 {
   margin: 0;
   position: relative;
   padding: 4px 10px;
   background: #568;
   color: white;
}
div.text {
   padding: 10px;
}
div.text >>> p {
   margin: 0 0;
}
span.date {
   position: absolute;
   right: 10px;
   top: 6px;
   font-weight: 100;
   font-size: 0.8em; 
}
</style>