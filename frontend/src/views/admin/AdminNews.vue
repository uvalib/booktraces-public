<template>
   <div class="admin">
      <h2>System Admin Panel <span class="login"><b>Logged in as:</b>{{loginName}}</span></h2>
      <div>
         <h3>News <span @click="addNewsClick" class="add pure-button pure-button-primary">Add News Item</span></h3>
         <div class="error">{{error}}</div>
         <div class="news-item" v-for="(item,idx) in news" :key="item.id">
            <div class="controls">
               <template v-if="!editingItem(item.id)">
                  <i :data-id="item.id" title="Delete" class="action fas fa-trash-alt" @click="deleteClicked"></i>
                  <i :data-id="item.id" title="Edit" class="action fa fa-edit" @click="editClicked"></i>
               </template>
               <template v-else>
                  <i class="action cancel fas fa-times-circle" @click="cancelClicked"></i>
                  <i class="action save fas fa-check-circle" @click="saveClicked"></i>
               </template>
            </div>
            <template v-if="editingItem(item.id)">
               <div><label>Title:</label><input class="title" type="text" v-model="editDetails.title"></div>
               <vue-editor v-model="editDetails.content"></vue-editor>
            </template>
            <template v-else>
               <div class="title">{{item.title}}</div>
               <div class="date">{{formatDate(item.createdAt)}}</div>
               <div @click="togglePublication(idx)" title="Toggle publication status" class="published">{{publishedText(item)}}</div>
               <div class="text" v-html="item.content"></div>
            </template>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState } from 'vuex'
import { mapGetters } from 'vuex'
import { VueEditor } from 'vue2-editor'
export default {
   name: "admin-news",
   components: {
      "vue-editor": VueEditor
   },
   data: function () {
      return {
         edit: false,
         editDetails: null,
         addingNew: false,
      }
   },
   computed: {
      ...mapState({
         error: state => state.error,
         loading: state => state.loading,
         news: state => state.news.list,
      }),
      ...mapGetters({
         loginName: 'admin/loginName',
      })
   },
   methods: {
      togglePublication(idx) {
         this.$store.dispatch("news/togglePublication", idx)
      },
      publishedText( item ) {
         if (item.published) return "Published"
         return "Not Published"
      },
      formatDate(date) {
         return date.split("T")[0]
      },
      addNewsClick() {
         this.$store.commit("news/addPlaceholder")
         this.editDetails = Object.assign({}, this.news[0])
         this.edit = true
         this.addingNew = true
      },
      cancelClicked() {
         this.edit = false
         this.editDetails = null
         if (this.addingNew) {
            this.addingNew = false
            this.$store.commit("news/cancelAdd")
         }
      },
      saveClicked() {
         if (this.addingNew) {
            this.$store.dispatch('news/addNews',this.editDetails).then((/*response*/) => {
               this.edit=false
               this.editDetails = null
               this.addingNew = false
            })
         } else {
            this.$store.dispatch('news/updateNews',this.editDetails).then((/*response*/) => {
               this.edit=false
               this.editDetails = null
               this.addingNew = false
            })
         }
      },
      editingItem(id) {
         if (this.edit == false) return false 
         return this.editDetails.id == id
      },
      deleteClicked(event) {
         let tgt = event.currentTarget
         let newsID = tgt.dataset.id
         let resp = confirm("Delete this news item? All data will be permanently lost. Continue?")
         if (resp) {
             this.$store.dispatch('news/deleteNews', newsID)
         }
      },
      editClicked(event) {
         let tgt = event.currentTarget
         var newsID = tgt.dataset.id
         var tgtIdx = -1
         this.news.some( function(e,idx) {
          if (e.id == newsID) {
            tgtIdx = idx
            return true
          }
          return false
        })
        if (tgtIdx > -1) {
           this.editDetails = Object.assign({}, this.news[tgtIdx])
           this.edit = true
        }
      },
   },
   created() {
      this.$store.dispatch('news/getAll')
   },
};
</script>

<style scoped>
i.action.cancel, .error {
   color: firebrick;
}
i.action.save {
   color: green;
}
i.action {
   color: #666;
   opacity: 0.6;
   margin: 0 5px;
}
i.action:hover {
   cursor: pointer;
   opacity: 1;
}
span.login {
   font-family: sans-serif;
   font-size: 0.5em;
   float: right;
   font-weight: 100;
}
span.login b {
   margin-right: 5px;
}
h2 {
  font-size: 1.5em;
  font-weight: bold;
  border-bottom: 1px dashed #666;
  font-family: 'Special Elite', cursive;
  padding-bottom: 5px;
  margin-bottom: 15px;
}
div.admin {
   padding: 15px 25px;
   min-height: 600px;
   background: white;
   color: #444;
}
h3 {
   position: relative;
}
h3 span.pure-button.add {
   font-size: 0.6em;
   font-weight: 100;
   position: absolute;
   right: 0;
   bottom:0;
}
.news-item {
   margin: 10px 0 25px 0;
   padding: 15px;
   border:1px solid #ccc;
   position: relative;
}
div.text {
   margin-top: 15px;
}
div.published {
   font-weight: bold;
   color: cornflowerblue;
}
div.published:hover {
   text-decoration: underline;
   cursor: pointer;
}
div.controls {
   position: absolute;
   right: 15px;
   top: 15px;
}
div.text >>> p {
   margin: 0px 0;
}
div.title {
   color: #444;
   font-weight: bold;
}
div.date {
   font-weight: 100;
   font-size: 0.9em;
}
label {
   display: inline-block;
   font-weight: bold;
   color: #555; 
   text-align: right;
   margin-right: 15px;
}
input.title {
   width: 800px;
   margin-bottom: 10px;
}
</style>