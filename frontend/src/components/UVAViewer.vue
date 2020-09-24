<template>
   <div class="uva-viewer">
      <div class="image-wrap" @mousedown="mouseDown()" :style="{height: `${height}px`}" >
         <img class="thumb" :src="url" :style="{transform: rotation}"/>
      </div>
      <div class="ctrl">
         <i @click="rotLeft()" class="btn fas fa-undo-alt ccw"></i>
         <span>
            <i @click="zoomOut()" class="btn zoom fas fa-search-minus"></i>
            <span @click="reset()" style="font-size:0.95em" class="btn">Reset</span>
            <i @click="zoomIn()" class="btn zoom fas fa-search-plus"></i>
         </span>
         <i @click="rotRight()" class="btn fas fa-undo-alt cw"></i>
      </div>
   </div>
</template>

<script>

export default {
   props: {
      url: {
         type: String,
         required: true
      },
      height: {
         type: Number,
         default: 400
      }
   },
   data: function () {
      return {
         rotate: 0,
         zoom: 1.0,
         translateX: 0,
         translateY: 0,
         currX: 0,
         currY: 0,
         dragging: false
      }
   },
   computed: {
      rotation() {
         return `rotate(${this.rotate}deg) translateX(${this.translateX}px)  translateY(${this.translateY}px) scale(${this.zoom})`
      },
   },
   methods: {
      reset() {
         this.rotate = 0
         this.zoom = 1.0
         this.translateX = 0
         this.translateY = 0
         this.currX = 0
         this.currY = 0
      },
      zoomIn() {
         this.zoom += 0.5
      },
      zoomOut() {
         if (this.zoom > 0.5) {
            this.zoom -= 0.5
         }
      },
      mouseDown() {
         event.preventDefault()
         event.stopPropagation()
         this.dragging = true
         this.currX = event.screenX
         this.currY = event.screenY
         //console.log(`MOUSE: ${event.screenX}, ${event.screenY}`)
      },
      mouseUp() {
         this.dragging = false
      },
      mouseMove() {
         if ( this.dragging ) {
            event.preventDefault()
            event.stopPropagation()
            let delX = event.screenX - this.currX
            let delY = event.screenY - this.currY
            //console.log(`CURR: ${this.rotate} deg - ${this.currX}, ${this.currY}`)
            //console.log(`MOUSE: ${event.screenX}, ${event.screenY}`)

            if ( this.rotate == 90 || this.rotate == -270) {
               delY = this.currX - event.screenX
               delX = event.screenY - this.currY
            } else if ( this.rotate == 270 || this.rotate == -90 ) {
               delY = event.screenX - this.currX
               delX = this.currY - event.screenY
            } else if ( this.rotate == 180 || this.rotate == -180 ) {
               delX = this.currX - event.screenX
               delY = this.currY - event.screenY
            }

            this.translateX += delX
            this.translateY += delY
            this.currX = event.screenX
            this.currY = event.screenY
            //console.log(`DEL: ${delX}, ${delY}`)

         }
      },
      rotRight() {
         if ( this.rotate == 270) {
            this.rotate = 0
         } else {
            this.rotate += 90
         }
      },
      rotLeft() {
         if ( this.rotate == -270) {
            this.rotate = 0
         } else {
            this.rotate -= 90
         }
      },
   },
   mounted() {
      window.addEventListener("mouseup", this.mouseUp)
      window.addEventListener("mousemove", this.mouseMove)
   },
   destroyed() {
      window.removeEventListener("mouseup", this.mouseUp)
      window.removeEventListener("mousemove", this.mouseMove)
   }
}
</script>

<style scoped>
.uva-viewer {
   position: relative;
   background: #222;
   border: 1px solid #999;
}
.uva-viewer img {
   display: inline-block;
   width: 100%;
}
.image-wrap {
   overflow: hidden;
}
.uva-viewer .ctrl {
  display: flex;
  flex-flow: row nowrap;
  justify-content: space-between;
  background: #444;
  padding: 5px;
  border-top: 1px solid black;
  border-bottom: 1px solid black;
}
.uva-viewer .ctrl .btn {
   color: white;
   background: black;
   padding: 10px;
   border-radius: 5px;
   opacity: 0.45;
}
.btn.cw {
      transform: scaleX(-1);
}
.zoom {
   margin: 0 5px;
}
</style>