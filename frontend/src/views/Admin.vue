<template>
   <div class="admin">
      <h2>Book Traces System Admin Panel</h2>
   </div>
</template>

<script>
export default {
   name: "admin",
   created() {
      this.$store.dispatch("admin/getSubmissions")
   }
};
</script>

<style scoped>
div.admin {
   padding: 15px 25px;
   min-height: 600px;
   background: white;
   color: #444;
}
</style>