<template>
   <div class="about content">
      <h2>About</h2>
      <dl>
         <dt>Director</dt>
         <dd>
            <a href="https://english.as.virginia.edu/people/profile/ams4k" target="_blank">Andrew Stauffer</a>,&nbsp;
            Department of English, University of Virginia
         </dd>
         <dt>Project Manager</dt>
         <dd>
            <a href="https://www.library.virginia.edu/staff/khj5c/" target="_blank">Kristin Jensen</a>,&nbsp;
            University of Virginia Library
         </dd>
      </dl>
      <p>
         <router-link to="/">Book Traces</router-link> is a collaborative web project aimed at discovering, cataloging, and preserving
         unique copies of
         <b>nineteenth- and early twentieth-century books on library shelves</b>.
         Our focus is on
         <b>circulating books</b> (not special collections) that were customized by their original owners, who left in them
         marginalia, annotations, inscriptions, and insertions. At least until the 1950s, academic libraries in particular built their
         collections from donations, collecting secondhand books from alumni, local families, faculty, and others.
         These books thus constitute a
         <b>massive, distributed archive of the history of reading and book use,</b>
         hidden in plain sight in the stacks of libraries. We want to find these books and make them available to researchers.
      </p>
      <p>
         Book Traces is meant to engage
         <b>the question of the future of the print record</b>
         in the wake of wide-scale digitization. The issue is particularly urgent for the materials from the long nineteenth century.
         In most cases, pre-1800 books have been moved to special collections, and post-1924 materials remain in copyright and thus on the
         shelves for circulation. But academic libraries are now reconfiguring access to public-domain texts via digital repositories such as
         Hathi Trust and Google Books. As a result, we are facing the
         <a href="https://www.tandfonline.com/doi/full/10.1080/01930826.2015.1105668" target="_blank">withdrawal of many books</a>
         in favor of digital surrogates and
         <a href="https://www.oclc.org/research/publications/library/2013/2013-09r.html" target="_blank">
         shared access to books across libraries</a>.
      </p>
      <p>
         We believe that special attention should be paid to these historical materials, which <b>bear marks of use by their original owners.</b>
         &nbsp;<router-link to="/submissions/1">Marginalia</router-link>,
         &nbsp;<router-link to="/submissions/361">inscriptions</router-link>,
         &nbsp;<router-link to="/submissions/69">photos</router-link>,
         &nbsp;<router-link to="/submissions/1233">original manuscripts</router-link>,
         &nbsp;<router-link to="/submissions/508">letters</router-link>,
         &nbsp;<router-link to="/submissions/578">drawings</router-link>,
         &nbsp;<router-link to="/submissions/397">locks of hair</router-link>,
         &nbsp;<router-link to="/submissions/1196">flowers</router-link>, and many other unique pieces
         of historical data can be found in individual copies, many of them associated with the history of the institution that
         collected the books in the first place. These unique attributes cannot be located by any electronic catalog.
         <span class="big">Each book has to be opened and examined.</span>
      </p>
      <p>
        <span class="big">That’s where you come in.</span> Thousands of such books — marked or otherwise customized by nineteenth-century
        owners – are on library shelves. The Book Traces project hopes to inspire the development of processes
        for discovering them, cataloguing them more fully, preserving them, and making better-informed decisions
        about print collections management. Further, we hope to model this process as part of a national effort
        for coordinating the management of nineteenth-century printed materials in academic libraries.
      </p>
      <p>
         Book Traces is sponsored by <a href="http://nines.org" target="_blank">NINES</a> at
         the <a href="http://virginia.edu" target="_blank">University of Virginia</a>, with
         <a href="https://news.library.virginia.edu/2017/06/01/hidden-in-plain-sight-book-traces-uva-a-clir-hidden-collections-project/"
            targte="_blank">funding from CLIR</a> and the
         <a href="https://www.mellon.org/grants/grants-database/grants/university-of-virginia/1901-06518/" target="_blank">Andrew W. Mellon Foundation</a>.
      </p>
      <p>
        If you have any questions or remarks, please contact <a href="mailto:ams4k@virginia.edu">Andrew Stauffer (ams4k@virginia.edu)</a>
        or <a href="mailto:khj5c@virginia.edu">Kristin Jensen (khj5c@virginia.edu)</a>.
        And please, <b>share the site widely!</b> We want as many seekers out there in the library stacks as we can muster.
        The more examples we can locate, the more we will know about the history and future of the nineteenth-century book. We can also help co-ordinate
        class exercises based on stack searches in your library.
      </p>
      <h4>Book Traces Scholarly Advisory Board</h4>
      <ul>
         <li>Thomas Augst (NYU)</li>
         <li>Michael Cohen (UCLA)</li>
         <li>Patricia Crain (NYU)</li>
         <li>Steve Cushman (UVA)</li>
         <li>Lindsey Eckert (Florida State)</li>
         <li>Sarah Gardner (Mercer</li>
         <li>Caroline Janney (UVA)</li>
         <li>Mary Louise Kete (Vermont)</li>
         <li>Amanda Licastro (Stevenson)</li>
         <li>Devoney Looser (Arizona State)</li>
         <li>Christina Lupton (Copenhagen)</li>
         <li>Deborah Lutz (Louisville)</li>
         <li>Deidre Lynch (Harvard)</li>
         <li>Jerome McGann (UVA)</li>
         <li>Elizabeth McHenry (NYU)</li>
         <li>Jonathan Mulrooney (Holy Cross)</li>
         <li>Andrew Piper (McGill)</li>
         <li>Leah Price (Rutgers)</li>
         <li>Jonathan Sachs (Concordia)</li>
         <li>Steve Stowe (Indiana)</li>
         <li>Mark Towsey (Liverpool)</li>
         <li>Jon Wells (Michigan)</li>
      </ul>
      <h4>Partner Libraries</h4>
      <ul>
         <li>Arizona State University</li>
         <li>Brandeis University</li>
         <li>Bryn Mawr College</li>
         <li>Louisiana State University</li>
         <li>University of Connecticut</li>
         <li>University of Miami</li>
         <li>University of Virginia</li>
      </ul>
   </div>
</template>

<script>
export default {
   name: "about",
   methods: {}
};
</script>

<style scoped>
dl {
   padding-bottom: 10px;
   margin-bottom: 10px;
   border-bottom: 1px dashed;
}
dt, dd {
   margin-bottom: 10px;
}
dt {
   font-weight: bold;
}
span.big {
   font-weight: bold;
   font-size: 1.1em;
}
</style>