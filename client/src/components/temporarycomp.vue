<script>
export default {
   methods: {
      toggleDrowdown() {
         this.isButtonClicked = !this.isButtonClicked
         this.buttonClassNames = (!this.isButtonClicked)
            ? ["hidden", "absolute right-0", "mt-2", "rounded-md", "shadow-lg", "bg-white", "ring-1", "ring-black", "ring-opacity-5", "p-1", "space-y-1", "transition-all", "duration-500", "opacity-0"]
            : ["absolute right-0", "mt-2", "rounded-md", "shadow-lg", "bg-white", "ring-1", "ring-black", "ring-opacity-5", "p-1", "space-y-1", "transition-all", "duration-1000", "opacity-95"]
      },
      selectSearchOption(key, callback) {
         this.selectedSearchType = key
         let typesArr = this.searchTypes

         for (const item of typesArr) {
            this.IsTypeSelected(item)
         }
         callback()
      },
      IsTypeSelected(typesArrItem) {
         typesArrItem.selected = (typesArrItem.key === this.selectedSearchType)
            ? true
            : false
      },
      activateAdvancedOpctions(item) {
         item.isActivated = !item.isActivated
      },
      getSearchTypesClassName(item) {
         return ["text-xs", "block", "px-4", "py-2", "text-secondary", "hover:bg-third", "cursor-pointer", "rounded-md", item.selected ? "bg-third" : ""]
      },
      toggleActivatedOptions(item) {
         return ["transition", "duration-700", "rounded-lg", "px-4", item.isActivated ? "text-secondary bg-fourth" : ""]
      }
   },
   data() {
      return {
         isButtonClicked: false,
         selectedSearchType: "",
         buttonClassNames: ["hidden", "absolute right-0", "mt-2", "rounded-md", "shadow-lg", "bg-white", "ring-1", "ring-black", "ring-opacity-5", "p-1", "space-y-1"],
         searchTypes: [
            { name: "Match", key: "match", selected: false },
            { name: "Match Phrase", key: "matchphrase", selected: false },
            { name: "Term", key: "term", selected: false },
            { name: "Prefix", key: "prefix", selected: false },
            { name: "Wildcard", key: "wildcard", selected: false },
            { name: "Fuzzy", key: "fuzzy", selected: false },
         ],
         advancedOptions: {
            isActivated: false,
         },
         isTagGreaterThanZero: false
      }
   }
}
</script>

<template>
   <div class="grid grid-cols-1  grid-rows-[1fr,2fr,2fr] items-center h-96 text-center  bg-primary">
      <i>
         <div class="grid grid-rows-1 grid-cols-3 items-center text-gray-400 m-auto">
            <i></i>
            <i> Logo here</i>
            <i></i>
         </div>
      </i>
      <i>
         <div class="grid grid-cols-[1fr,3fr,1fr] md:grid-cols-[9fr,21fr,9fr]  grid-rows-2 items-center text-third">
            <i></i>
            <i>
               <div
                  class="grid grid-rows-1 grid-cols-[5fr,21fr,5fr] md:grid-cols-[3fr,21fr,3fr]  items-center  bg-secondary rounded-3xl h-12 w-72 md:w-[44rem] p-1 m-auto">
                  <i>

                     <div class="min-h-screen flex items-center justify-start">
                        <div class="relative group">
                           <button ref="dropdownButton" @click="toggleDrowdown($event)"
                              class="inline-flex justify-end w-full px-4 py-2 text-sm font-medium text-gray-700 rounded-xl shadow-sm">
                              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                                 stroke="currentColor" class="w-6 h-6 text-third">
                                 <path stroke-linecap="round" stroke-linejoin="round"
                                    d="M6 13.5V3.75m0 9.75a1.5 1.5 0 0 1 0 3m0-3a1.5 1.5 0 0 0 0 3m0 3.75V16.5m12-3V3.75m0 9.75a1.5 1.5 0 0 1 0 3m0-3a1.5 1.5 0 0 0 0 3m0 3.75V16.5m-6-9V3.75m0 3.75a1.5 1.5 0 0 1 0 3m0-3a1.5 1.5 0 0 0 0 3m0 9.75V10.5" />
                              </svg>
                           </button>

                           <div ref="dropdownMenu" :class="buttonClassNames">
                              <label for="search-types" class="text-secondary capitalize">
                                 search types
                              </label>
                              <!-- Search input -->

                              <!-- Dropdown content goes here -->
                              <a @click="selectSearchOption(item.key, toggleDrowdown)" v-for="item in searchTypes"
                                 :key="item.key" :class="getSearchTypesClassName(item)">{{
                              item.name }}</a>
                           </div>
                        </div>
                     </div>
                  </i>
                  <i>
                     <input ref="searchInput"
                        class="block w-full py-2 text-fourth rounded-lg focus:outline-none bg-secondary " type="text"
                        placeholder="Search" autocomplete="off">
                  </i>
                  <i>
                     <div class="grid grid-cols-1 grid-rows-1 place-items-end px-4 py-2">
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                           stroke="currentColor" class="w-6 h-6 text-third">
                           <path stroke-linecap="round" stroke-linejoin="round"
                              d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" />
                        </svg>
                     </div>
                  </i>
               </div>
            </i>
            <i></i>
            <i></i>
            <i>
               <div
                  class="grid grid-cols-[1fr,1fr,3fr] md:grid-cols-3 grid-rows-1 text-gray-400 items-center text-end px-2  md:px-8">
                  <i></i>
                  <i></i>
                  <i>
                     <div class="grid grid-cols-1 items-end">
                        <i>
                           <button @click="activateAdvancedOpctions(advancedOptions)"
                              :class="toggleActivatedOptions(advancedOptions)">
                              Advanced options
                           </button>
                        </i>
                     </div>
                  </i>
               </div>
            </i>
            <i></i>
         </div>
      </i>
      <i>
         <transition>
            <div v-show="advancedOptions.isActivated"
               class="grid grid-cols-1 grid-rows-1 place-items-center text-gray-500 opacity-95">
               <i>
                  <div class="max-w-lg m-6">
                     <div class="relative">
                        <input
                           class="appearance-none block w-full bg-third text-gray-700 border border-gray-200 rounded py-2 px-4 leading-tight"
                           placeholder="Enter some tags">
                        <div class="hidden">
                           <div class="absolute z-40 left-0 mt-2 w-full">
                              <div class="py-1 text-sm bg-white rounded shadow-lg border border-gray-300">
                                 <a class="block py-1 px-5 cursor-pointer hover:bg-indigo-600 hover:text-white">Add
                                    tag"<span class="font-semibold" x-text="textInput"></span>"</a>
                              </div>
                           </div>
                        </div>
                        <!-- selections -->
                        <div v-show="isTagGreaterThanZero"
                           class="bg-secondary inline-flex items-center text-sm rounded mt-2 mr-1 overflow-hidden">
                           <span class="ml-2 mr-1 leading-relaxed truncate max-w-xs px-1" x-text="tag"></span>
                           <button
                              class="w-6 h-8 inline-block align-middle text-gray-500 bg-secondary focus:outline-none">
                              <svg class="w-6 h-6 fill-current mx-auto" xmlns="http://www.w3.org/2000/svg"
                                 viewBox="0 0 24 24">
                                 <path fill-rule="evenodd"
                                    d="M15.78 14.36a1 1 0 0 1-1.42 1.42l-2.82-2.83-2.83 2.83a1 1 0 1 1-1.42-1.42l2.83-2.82L7.3 8.7a1 1 0 0 1 1.42-1.42l2.83 2.83 2.82-2.83a1 1 0 0 1 1.42 1.42l-2.83 2.83 2.83 2.82z" />
                              </svg>
                           </button>
                        </div>
                     </div>
                  </div>
               </i>
            </div>
         </transition>
      </i>
   </div>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
   transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
   opacity: 0;
}
</style>
