<template>
   <div class="content submit">
      <h2>Submit a Book</h2>
      <div class="info">
         <p>We are currently collecting pre-1923 copies of library books. Share your finds with us!</p>
         <p>Submit images and information about your library book here:</p>
      </div>
      <form>
         <div class="row">
            <label for="title">Title of Book<span class="required">*</span></label>
            <InputText id="title" v-model="upload.title" />
         </div>
         <div class="row">
            <label for="author">Author (Last name, first name)<span class="required">*</span></label>
            <InputText id="author" v-model="upload.author" />
         </div>
         <div class="row">
            <label for="publication">Publication place/date (e.g., London, 1888)</label>
            <InputText id="publication" v-model="upload.publication" />
         </div>
         <div class="row">
            <label for="institution">Institution where found<span class="required">*</span></label>
            <Select v-model="upload.institution" :options="system.institutions"
               optionLabel="name" :editable="true"
               placeholder="Select or create an institution" />
         </div>
         <div class="row">
            <label for="callnumber">Call Number</label>
            <InputText id="callnumber" v-model="upload.callNumber" />
         </div>
         <div class="row">
            <label for="description">Description</label>
            <Textarea id="description" v-model="upload.description" rows="4" />
         </div>
         <div class="row">
            <label for="images">Images of unique features (e.g., marginalia, inserts)<span class="required">*</span></label>
            <FileUpload id="images" mode="advanced" url="/api/upload" name="uploadimage"
               :showUploadButton="false" :showCancelButton="false"
               accept="image/jpeg,image/png" :multiple="false"
               @upload="imageUploaded" @before-upload="beforeImageUpload" @removeUploadedFile="imageRmoved"
               :auto="true" chooseLabel="Upload Images"
            >
               <template #empty>
                  <div class="upload title">Drag and drop to upload content</div>
                  <div class="upload subtitle">or click 'Upload Images' to select a file from your computer</div>
               </template>
            </FileUpload>
         </div>
         <div class="row">
            <label for="submitter">Submitted By<span class="required">*</span></label>
            <InputText id="submitter" v-model="upload.submitter" />
            <p class="note">The name you provide will appear on the Booktraces site along with your submission.</p>
         </div>
         <div class="row">
            <label for="email">Email Address<span class="required">*</span></label>
            <InputText id="email" v-model="upload.email" />
            <p class="note">We will keep your email address private and will only use if if we need to contact you about your submission.</p>
         </div>
         <div class="row">
            <label>Post Tags<span class="required">*</span></label>
            <div class="tags">
               <div v-for="tag in system.tags" :key="tag.name" class="tag">
                  <Checkbox v-model="upload.tags" :inputId="`tag-${tag.id}`" name="tag" :value="tag.id" />
                  <label :for="`tag-${tag.id}`">{{ tag.name }}</label>
               </div>
            </div>
            <p class="note">At least one tag is required.</p>
         </div>
         <div class="error" v-if="upload.error">{{upload.error}}</div>
      </form>
      <Button @click="submitClicked" label="SUBMIT" :disabled="upload.canSubmit == false || upload.working"/>
      <p class="error" v-if="!upload.canSubmit">Fill in all required fields to enable the submit button.</p>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useSystemStore } from "@/stores/system"
import { useUploadStore } from "@/stores/uploads"
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import FileUpload from 'primevue/fileupload'
import Checkbox from 'primevue/checkbox'

const system = useSystemStore()
const upload = useUploadStore()

onMounted(() => {
   upload.getUploadID()
   system.getTags()
   system.getInstitutions()
})

const beforeImageUpload = (( evt ) => {
   evt.formData.append('uploadID', upload.uploadID);
})
const imageUploaded = ((evt) => {
   upload.addImage( evt.files[0] )
})
const imageRmoved = ((evt) => {
   upload.removeUploadedImage( evt.file.name )
})

const submitClicked = (() => {
   upload.submit()
})

</script>

<style lang="scss" scoped>
.error {
   margin: 10px 0;
   color: firebrick;
   font-style: italic;
   padding: 10px;
   background: #fff0f0;
   border-radius: 5px;
   border: 1px solid firebrick;
}

.info {
   margin-bottom: 25px;
}

form {
   display: flex;
   flex-direction: column;
   margin-bottom: 25px;
   gap: 15px;
   .row {
      display: flex;
      flex-direction: column;
      gap: 5px;
      .tags {
         display: flex;
         flex-flow: row wrap;
         gap: 15px 20px;
         margin: 5px 0;
         border: 1px solid #dedede;
         border-radius: 5px;
         padding: 20px;
         .tag {
            display: flex;
            flex-flow: row nwrap;
            gap: 10px;
            align-items: center;
            justify-content: flex-start;
         }
      }
   }
   .note {
      font-style: italic;
      color: #999;
      font-size: .9em;
      margin: 0;
   }
   .required {
      font-weight: bold;
      padding-left: 5px;
      color: firebrick;
   }
   .upload.title {
      font-weight: 700;
      font-size: 1.2em;
   }
}

</style>
