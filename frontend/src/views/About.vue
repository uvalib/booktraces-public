<template>
   <div class="about content">
      <h2>About</h2>
      <p>
         “…some hand, now forever quiet, had made at certain passages strong pen and ink marks,
         long since browned by time.” — George Eliot, The Mill on the Floss (1860)
      </p>
      <p>
         <router-link to="/">Book Traces</router-link> is a crowd-sourced web project aimed at identifying unique copies of
         <b>nineteenth- and early twentieth-century books on library shelves.</b> Our focus is on customizations made by original owners
         in personal copies, primarily in the form of marginalia, annotation, inscriptions, and insertions.
      </p>
      <p>
        Sponsored by <a href="http://nines.org" target="_blank">NINES</a> 
        at the <a href="http://virginia.edu" target="_blank">University of Virginia</a> and led by 
        <a href="http://www.engl.virginia.edu/people/ams4k" target="_blank">Andrew Stauffer</a>, Book Traces is 
        meant to engage <b>the question of the future of the print record</b> in the wake of wide-scale digitization. 
        The issue is particularly urgent for the materials from the long nineteenth century. In most cases, 
        pre-1800 books have been moved to special collections, and post-1923 materials remain in copyright 
        and thus on the shelves for circulation. But academic libraries are now reconfiguring access to public-domain 
        texts via repositories such as Hathi Trust and Google Books. The nineteenth-century books themselves aren’t 
        getting much use. As a result, we are anticipating the withdrawal of large portions of the print collections in favor of 
        digital surrogates and shared access to books across libraries.
      </p>
      <p>
        We believe that, as libraries begin to manage down their print holdings, special attention 
        should be paid to these historical materials, which in many cases came to university libraries 
        from alumni donors and <b>bear marks of use by their original nineteenth-century owners.</b> 
        These books thus constitute a <b>massive, distributed archive of the history of reading and book-use,</b> 
        hidden in plain sight in the circulating collections. Marginalia, inscriptions, photos, original manuscripts, 
        letters, drawings, and many other unique pieces of historical data can be found in individual copies, 
        many of them associated with the history of the institution that collected the books in the first place. These 
        unique attributes cannot be located by any electronic catalog. <span class="big">Each book has to be opened and examined.</span>
      </p>
      <p>
        <span class="big">That’s where you come in.</span> Thousands of such books — marked or otherwise customized by nineteenth-century 
        owners – are on library shelves. The Book Traces project hopes to inspire the development of processes 
        for discovering them, cataloguing them more fully, preserving them, and making better-informed decisions 
        about print collections management. Further, we hope to model this process as part of a national effort 
        for coordinating the management of nineteenth-century printed materials in academic libraries.
      </p>
      <p>
        If you have any questions or remarks, please contact <a href="mailto:ams4k@virginia.edu">Andrew Stauffer (ams4k@virginia.edu)</a> 
        or <a href="mailto:khj5c@virginia.edu">Kristin Jensen (khj5c@virginia.edu)</a>. 
        And please, <b>share the site 
        widely!</b> We want as many seekers out there in the library stacks as we can muster. The more examples 
        we can locate, the more we will know about the history and future of the nineteenth century book.
      </p>
      <p><router-link to="/press">Want more information on Book Traces?</router-link></p>
   </div>
</template>

<script>
export default {
   name: "about",
   methods: {
   }
};
</script>

<style scoped>
span.big {
   font-weight: bold;
   font-size: 1.1em;
}
</style>