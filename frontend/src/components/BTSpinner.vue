<template>
   <div v-if="overlay" class="bt-spinner-overlay">
      <div class="bt-spinner">
         <h3 v-if="message">{{ message }}</h3>
         <div class="spinner-animation">
            <div class="bounce bounce1"></div>
            <div class="bounce bounce2"></div>
            <div class="bounce bounce3"></div>
         </div>
      </div>
   </div>
   <div v-else class="bt-spinner embed">
      <h3 v-if="message">{{ message }}</h3>
      <div class="spinner-animation">
         <div class="bounce bounce1"></div>
         <div class="bounce bounce2"></div>
         <div class="bounce bounce3"></div>
      </div>
   </div>
</template>

<script setup>

const props = defineProps({
   message: {
      type: String,
      default: "",
   },
   overlay: {
      type: Boolean,
      default: false,
   },
})
</script>

<style lang="scss" scoped>
div.bt-spinner-overlay {
   position: fixed;
   left: 0;
   top: 0;
   width: 100%;
   height: 100%;
   z-index: 1000;
   background: rgba(0, 0, 0, 0.1);
   div.bt-spinner {
      text-align: center;
      background: white;
      margin: 12vw auto;

      padding: 10px 150px 25px 150px;
      display: inline-block;
      font-weight: bold;
      color: #444;
      box-shadow: 0 3px 6px rgba(0, 0, 0, 0.16), 0 3px 6px rgba(0, 0, 0, 0.23);
      h3 {
         text-align: center;
      }
   }
}

div.bt-spinner.embed {
   box-shadow: none;
   padding: 0;
   margin: 0;
   background: transparent;
   margin-top: 6%;
   text-align: center;
   h3 {
      text-align: center;
   }
}

@media only screen and (min-width: 768px) {
   div.bt-spinner {
      padding: 40px 90px;
   }
}
@media only screen and (max-width: 768px) {
   div.bt-spinner {
      width: 95%;
      padding: 40px 0;
      margin-top: 30%;
   }
}

.spinner-animation {
   margin: 0 auto;
   width: 80px;
   text-align: center;
   .bounce {
      background-color: #888;
   }
}
.spinner-animation > div {
   width: 18px;
   height: 18px;
   border-radius: 100%;
   display: inline-block;
   -webkit-animation: sk-bouncedelay 1.4s infinite ease-in-out both;
   animation: sk-bouncedelay 1.4s infinite ease-in-out both;
   margin: 0 2px;
}

.spinner-animation .bounce1 {
   -webkit-animation-delay: -0.32s;
   animation-delay: -0.32s;
}

.spinner-animation .bounce2 {
   -webkit-animation-delay: -0.16s;
   animation-delay: -0.16s;
}
@-webkit-keyframes sk-bouncedelay {
   0%,
   80%,
   100% {
      -webkit-transform: scale(0);
   }
   40% {
      -webkit-transform: scale(1);
   }
}
@keyframes sk-bouncedelay {
   0%,
   80%,
   100% {
      -webkit-transform: scale(0);
      transform: scale(0);
   }
   40% {
      -webkit-transform: scale(1);
      transform: scale(1);
   }
}
</style>
