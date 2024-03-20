<template>
   <div>
      <form @submit.prevent="searchInIndexer(searchQuery,searchType,MaxResults,Field)">
         <div
            class="grid grid-cols-[1fr,3fr,1fr] md:grid-cols-[9fr,21fr,9fr] grid-rows-2 text-fourth items-center text-center my-16 mx-4">
            <i></i>
            <i>
               <input type="text" placeholder="your query string here"
                  class="input input-bordered w-full max-w-xs md:max-w-md" v-model="searchQuery" required />
            </i>
            <i></i>
            <i></i>
            <i>
               <select class="select select-bordered w-full max-w-xs md:max-w-md text-fourth" required>
                  <option disabled selected value="">Search Type</option>
                  <option @click="recordSearchType(item.key)" v-for="item in searchTypes" :key="item.key">{{ item.name
                     }}
                  </option>
               </select>
            </i>
            <i></i>
         </div>
      </form>



      <div class="grid grid-cols-[1fr,42fr,1fr] md:grid-cols-[5fr,21fr,5fr] h-[36em]">
         <i></i>
         <i>
            <div class="grid grid-rows-[1fr,9fr] w-full border-opacity-50 bg-base-300 rounded-lg py-4 items-center">
               <i>
                  <div class="grid grid-cols-[3fr,3fr,6fr] items-center text-start px-4 ">
                     <i>
                        <h3>From</h3>
                     </i>
                     <i>
                        <h3>To</h3>
                     </i>
                     <i>
                        <h3>Subject</h3>
                     </i>
                  </div>
               </i>
               <i>
                  <div
                     class="grid grid-rows-1 card bg-base-300 place-items-start px-2 py-2 overflow-scroll h-[22rem] md:h-[32rem]">
                     <i>
                        <table class="table table-xs bg-base-200 rounded-none">
                           <tbody>
                              <tr v-for="item in ReturnedData" @click.prevent="ReturnRow(item)"
                                 class="cursor-pointer hover:bg-neutral">
                                 <td class="w-3/12 h-8 truncate border-t-[1px] border-neutral">
                                    {{ item._source.from }}
                                 </td>
                                 <td class="w-3/12 h-8 truncate border-t-[1px] border-neutral">
                                    {{ item._source.to }}
                                 </td>
                                 <td class="w-6/12 h-8 truncate border-t-[1px] border-neutral">
                                    <strong>{{ item._source.subject }}</strong>
                                 </td>
                              </tr>
                           </tbody>
                        </table>
                     </i>
                  </div>
               </i>
            </div>
         </i>
         <i></i>
      </div>

   </div>

</template>

<script>

import axios from 'axios'

export default {
   props: {
      MaxResults: {
         type: Number,
         default: 20
      },
      Field: {
         type: String,
         default: ""
      }
   },
   mounted() {
      this.getHealthCheck()
   },
   emits: ['returnedMail'],
   data() {
      return {
         ReturnedData: null,
         searchQuery: "",
         searchType: "",
         searchTypes: [
            { name: "Match", key: "match", selected: false },
            { name: "Match Phrase", key: "matchphrase", selected: false },
            { name: "Term", key: "term", selected: false },
            { name: "Prefix", key: "prefix", selected: false },
            { name: "Wildcard", key: "wildcard", selected: false },
            { name: "Fuzzy", key: "fuzzy", selected: false }
         ]
      }
   },
   methods: {
      recordSearchType(key) {
         this.searchType = key
      },
      fetchData() {
         fetch("/response.json")
            .then((res) => {
               if (!res.ok) {
                  throw new Error
                     (`HTTP error! Status: ${res.status}`);
               }
               return res.json();
            })
            .then((data) =>
               this.ReturnedData = data.hits.hits)
            .catch((error) =>
               console.error("Unable to fetch data:", error));
      },
      IsReturnedDataNull() {
         return (this.ReturnedData === null)
      },
      ReturnRow(item) {
         this.$emit("returnedMail", item._source)
      },
      getHealthCheck() {
         axios.get('http://localhost:4040/api/healthcheck').then((response) => { console.log(response) }).catch((error) => { console.log(error) })
      },
      searchInIndexer(searchQuery, searchType,MaxResults, Field) {
         axios.post('http://localhost:4040/api/search', {
            "type": searchType,
            "search_term": searchQuery,
            "field":Field,
            "sort_fields": ["-_score"],
            "from": 0,
            "max_results": MaxResults,
            "_source": []
         }).then((response) => {this.ReturnedData = response.data.hits.hits }).catch((error) => { console.log(error) })
      }
   }
}


</script>

<style scoped>
table {
   table-layout: fixed;
}

td {
   white-space: nowrap;
   overflow: hidden;
   /* <- this does seem to be required */
   text-overflow: ellipsis;
}
</style>