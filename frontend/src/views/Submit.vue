<template>
  <div class="content submit">
    <h2>Submit a Book</h2>
    <div class="info">
      <p>We are currently  collecting pre-1923 copies of library books.  Share your finds with us!</p>
      <p>Submit images and information about your library book here:</p>
    </div>
    <form class="pure-form pure-form-stacked">
      <fieldset>
        <div class="pure-u-1-1 gap">
          <label for="title">Title of Book<span class="required">*</span></label>
          <input id="title" class="pure-u-1-1" type="text" required>
        </div>
        <div class="pure-u-1-1 gap">
          <label for="author">Author (Last name, first name)<span class="required">*</span></label>
          <input id="author" class="pure-u-1-1" type="text" required>
        </div>
        <div class="pure-u-1-1 gap">
          <label for="publication">Publication place/date (e.g., London, 1888)</label>
          <input id="publication" class="pure-u-1-1" type="text">
        </div>
        <div class="pure-u-1-1 gap">
          <label for="institution">Institution where found<span class="required">*</span></label>
          <multiselect v-model="selectedInstitution" class="folders"  
              placeholder="Select or create an institution"
              :showLabels="false" 
              :searchable="true"
              :taggable="true"
              track-by="id" label="name"
              tagPlaceholder="Press enter to create a new institution"
              @tag="addInstitution"
              :options="institutions">
          </multiselect>
        </div>
        <div class="pure-u-1-1 gap">
          <label for="call-number">Call Number</label>
          <input id="call-number" class="pure-u-1-1" type="text">
        </div>
        <div class="pure-u-1-1 gap">
          <label for="description">Description</label>
          <textarea rows="4" id="description" class="pure-u-1-1"></textarea>
        </div>
        <div class="pure-u-1-1 gap">
          <label>Images of unique features (e.g., marginalia, inserts)<span class="required">*</span></label>
          <vue-dropzone :useCustomSlot=true id="customdropzone" 
            :options="dropzoneOptions" 
            v-on:vdropzone-sending="sendingEvent"
            v-on:vdropzone-success="fileAddedEvent"
            v-on:vdropzone-removed-file="fileRemovedEvent">
            <div class="dropzone-custom">
              <div class="upload title">Drag and drop to upload content</div>
              <div class="upload subtitle">or click to select a file from your computer</div>
            </div>
          </vue-dropzone>
        </div>
        <div class="pure-u-1-1 gap">
          <label for="submitter">Submitted By<span class="required">*</span></label>
          <input id="submitter" class="pure-u-1-1" type="text">
          <p class="note">The name you provide will appear on the Booktraces site along with your submission.</p>
        </div>
        <div class="pure-u-1-1 gap">
          <label for="email">Email Address<span class="required">*</span></label>
          <input id="email" class="pure-u-1-1" type="text">
          <p class="note">We will keep your email address private and will only use if if we need to contact you about your submission.</p>
        </div>
        <div class="pure-u-1-1 gap">
          <label for="tags">Post Tags<span class="required">*</span></label>
          <div class="choices">
            <span v-for="tag in tags" :key="tag.id">
               <label class="pure-checkbox inline">
                  <input type="checkbox" name="tag" :value="tag.id" v-model="selectedTags">
                  {{ tag.name }}
               </label>
            </span>
         </div>
         <p class="note">At least one tag is required.</p>
        </div>
      </fieldset>
    </form>
    <div class="error">{{error}}</div>
    <span @click="submitClicked" class="submit pure-button pure-button-primary">SUBMIT</span>
  </div>
</template>

<script>
import vue2Dropzone from 'vue2-dropzone'
import 'vue2-dropzone/dist/vue2Dropzone.min.css'
import axios from 'axios'
import { mapState } from 'vuex'
import Multiselect from 'vue-multiselect'
export default {
  name: 'submit',
  components: {
    vueDropzone: vue2Dropzone, Multiselect
  },
  data: function () {
    return {
      submitted: false,
      selectedTags: [],
      selectedInstitution: null,
      dropzoneOptions: {
        url: '/api/upload',
        createImageThumbnails: true,
        maxFilesize: null,
        chunking: true,
        chunkSize: 1000000, // bytes = 1Mb,
        addRemoveLinks: true,
        acceptedFiles: "image/jpeg,image/png"
      }
    }
  },
  computed: {
    ...mapState({
        error: state => state.error,
        uploadedFiles: state => state.public.uploadedFiles,
        uploadID: state => state.public.uploadID,
        tags: state => state.tags,
        institutions: state => state.institutions,
    }),
  },
  created: function () {
    this.$store.dispatch('getTags')
    this.$store.dispatch('getInstitutions')
  },
  methods: {
    addInstitution(newInstitutionName) {
        this.$store.dispatch("addInstitution", newInstitutionName).then( () => {
          this.institutions.some( (i)=> {
              if (i.name == newInstitutionName) {
                this.selectedInstitution = i
                return true
              }
              return false
          })
        }).catch((error) => {
          // TODO something else maybe?
          alert(error)
        })
    },
    fileAddedEvent (file) {
      this.$store.commit("public/addUploadedFile",file.name)
    },
    fileRemovedEvent (file) {
      if ( this.submitted === false) {
        this.$store.dispatch("public/removeUploadedFile",file.name)
      }
    },
    sendingEvent (_file, _xhr, formData) {
      formData.append('uploadID', this.uploadID);
    },
    submitClicked(/*event*/) {
      if (this.selectedTags.length == 0) {
         this.$store.commit("setError", "At least one tag is required") 
         return
      }
      let form = {
        uploadID: this.uploadID,
        title: document.getElementById("title").value,
        author: document.getElementById("author").value,
        publication: document.getElementById("publication").value,
        institution_id: this.selectedInstitution.id,
        institution: this.selectedInstitution.name,
        callNumber: document.getElementById("call-number").value,
        description: document.getElementById("description").value,
        files: this.uploadedFiles,
        submitter: document.getElementById("submitter").value,
        email: document.getElementById("email").value,
        tags: this.selectedTags
      }
      
      axios.post("/api/submit", form).then((/*response*/)  =>  {
        this.$store.commit("public/clearUploadedFiles")
        this.submitted = true
        this.$router.push("/thanks")
      }).catch((error) => {
        this.$store.commit("setError",error.response.data) 
      })
    }
  }
}
</script>

<style src="vue-multiselect/dist/vue-multiselect.min.css"></style>
<style>
.pure-form input[type=text].multiselect__input {
  border: none;
}
</style>
<style scoped>
/* NOTE: Styles for DropZone must go in an un-scoped style section. App.vue has them */   
.pure-button.submit {
   background: #24890d;
   font-weight: bold;
}
.error {
  margin: 5px 0 10px 0;
  color: firebrick;
  font-style: italic;
}
label.pure-checkbox.inline {
  display: inline-block;
  margin: 5px 10px;
}
p.note {
  font-style: italic;
  color: #999;
  font-size: 0.8em;
  margin:0;
}
.gap {
  margin-bottom: 10px;
}
.required {
  font-weight: bold;
  padding-left: 5px;
  color: firebrick;
}
div.dropzone-custom {
  color: #666;
}
#customdropzone  {
  padding: 5px;
  height: 180px;
}
div.upload.title {
  font-weight: bold;
  font-size: 1.5em;
  margin-top: 15px;
}
div.upload.subtitle {
  font-weight: 500;
  margin-top: 5px;
}
</style>
