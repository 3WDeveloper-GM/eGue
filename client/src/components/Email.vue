<template>
   <div class="grid grid-cols-[1fr,42fr,1fr] md:grid-cols-[5fr,21fr,5fr] items-center my-10  md:mx-10">
      <i></i>
      <i>
         <div class="grid grid-cols-1 border-opacity-50">
            <div
               class="grid grid-rows-[1fr,19fr] card bg-base-300 rounded-box place-items-start h-[36rem] md:h-[48rem] w-full overflow-x-clip px-5 py-10 text-balance break-all">
               <i>
                  <div @click.prevent="changePage()">
                     <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                        stroke="currentColor"
                        class="w-6 h-6 hover:bg-neutral transition-all duration-500 m-3 p-1 rounded-lg">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M10.5 19.5 3 12m0 0 7.5-7.5M3 12h18" />
                     </svg>
                  </div>
               </i>
               <i>
                  <div
                     class="grid grid-rows-1 card bg-base-200 place-items-start px-2 py-2 overflow-y-scroll h-[22rem] md:h-[40rem]">
                     <i>
                        <div v-if="IsCurrentMailNull(currentEmail)">
                           <h3>Nothing to show so far</h3>
                        </div>
                        <div v-else>
                           <h4 v-for="field in emailFields" class="text-gray-300 text-pretty w-full md:w-11/12">
                              <span class="capitalize"> <strong>{{ field }}: </strong> &nbsp;&nbsp; </span>
                              <span> {{ currentEmail[field] }}</span>
                           </h4>
                           <br>

                           <br>
                           <p class="text-gray-400 text-pretty w-full md:w-11/12"> {{ currentEmail.body }}</p>
                        </div>
                     </i>
                  </div>

               </i>
            </div>
         </div>
      </i>
      <i></i>
   </div>
</template>

<script>
export default {
   data() {
      return {
         emailFields: [
            "subject",
            "date",
            "to",
            "from",
            "cc",
            "bcc",
         ]
      }
   },
   props: {
      currentEmail: {
         type: Object,
         default: null
      }
   },
   methods: {
      printCurrentMail(currentEmail) {
         return (currentEmail === null) ? "There are no emails to display right now" : JSON.stringify(currentEmail)
      },
      IsCurrentMailNull(currentEmail) {
         return (currentEmail === null)
      },
      changePage() {
         this.$emit('page-change')
      }
   },
}
</script>

<style lang="scss" scoped></style>