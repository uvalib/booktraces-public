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
          <label for="author">Publication place/date (e.g., London, 1888)</label>
          <input id="author" class="pure-u-1-1" type="text">
        </div>
        <div class="pure-u-1-1 gap">
          <label for="library">Library where found<span class="required">*</span></label>
          <input id="library" class="pure-u-1-1" type="text">
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
          <input type="hidden" id="submitted-files" name="submitted-files" :value="uploadedFiles">
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
          <label for="tags">Post Tags</label>
          <div class="choices">
            <span v-for="tag in tags" :key="tag.id">
               <label class="pure-checkbox inline">
                  <input type="checkbox" name="tag" :value="tag.id">
                  {{ tag.name }}
               </label>
            </span>
         </div>
        </div>
        <button class="submit pure-button pure-button-primary">SUBMIT</button>
      </fieldset>
    </form>
  </div>
</template>

<script>
import vue2Dropzone from 'vue2-dropzone'
import 'vue2-dropzone/dist/vue2Dropzone.min.css'

export default {
  name: 'submit',
  components: {
    vueDropzone: vue2Dropzone
  },
  data: function () {
    return {
      dropzoneOptions: {
        url: '/api/upload',
        createImageThumbnails: true,
        maxFilesize: null,
        chunking: true,
        chunkSize: 1000000, // bytes = 1Mb,
        addRemoveLinks: true  
      }
    }
  },
  computed: {
    uploadedFiles: function() {
      return this.$store.getters.uploadedFiles
    },
    tags: function() {
      return this.$store.getters.tags
    }
  },
  created: function () {
    this.$store.dispatch('getTags')
    this.$store.dispatch('getUploadID')
  },
  methods: {
    fileAddedEvent (file) {
      // just adds filename to store list 
      this.$store.commit("addUploadedFile",file.name)
    },
    fileRemovedEvent (file) {
      // makes an ajax call to the service to remove the file
      this.$store.dispatch("removeUploadedFile",file.name)
    },
    sendingEvent (file, xhr, formData) {
      formData.append('identifier', this.$store.getters.uploadID);
    }
  }
}
</script>

<style scoped>
/* NOTE: Styles for DropZone must go in an un-scoped style section. App.vue has them */   
.pure-button.submit {
   background: #24890d;
   font-weight: bold;
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
