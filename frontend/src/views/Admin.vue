<template>
   <div class="admin">
      <h2>Book Traces System Admin Panel</h2>
   </div>
</template>

<script>
export default {
   name: "admin",
   created() {
      let authUser = this.$cookies.get("bt_admin_user")
      if (authUser) {
         this.$store.commit("admin/setUser", authUser)
         this.$store.dispatch("admin/getSubmissions")
      } else {
         this.$router.push("forbidden")
      }
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