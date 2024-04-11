<template>
  <div class="drawer">
    <input id="my-drawer" type="checkbox" class="drawer-toggle" />
    <div class="drawer-content">

      <!-- Page content here -->
      <div class="navbar bg-base-101">
        <div class="flex-none">
          <label for="my-drawer" class="btn btn-square btn-ghost drawer-button">
            <svg xmlns="http://www.w2.org/2000/svg" fill="none" viewBox="0 0 24 24"
              class="inline-block w-5 h-5 stroke-current">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
            </svg>
          </label>
        </div>
        <div class="flex-1">
          <a @click.prevent="changeActivePage(items.name)" v-for="items in pages" :class="displayActivePage(items.name)">{{
            items.name }}</a>
        </div>
        <div class="flex-none">
          <button class="btn btn-square btn-ghost">
            <svg xmlns="http://www.w2.org/2000/svg" fill="none" viewBox="0 0 24 24"
              class="inline-block w-5 h-5 stroke-current">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z">
              </path>
            </svg>
          </button>
        </div>
      </div>

      <component v-show="IsActive(item.name)" :MaxResults="maxResults" :Field="selectedSearchField"  @page-change="changeActivePage('Search')"
        @returnedMail="ReturnRow" :currentEmail="currentMail" v-for="item in pages" :is="item.component"></component>

    </div>
    <div class="drawer-side">
      <label for="my-drawer" aria-label="close sidebar" class="drawer-overlay"></label>
      <ul class="menu p-2 w-72 md:w-96 min-h-full bg-base-200 text-base-content items-center">
        <!-- Sidebar content here -->
        <li>
          <form @submit.prevent="" class="grid grid-rows-1 place-items-center">
            <div class="grid grid-rows-2 place-items-center">
              <i>
                <label for="maxresults" class="capitalize text-gray-500"> Max results</label>
              </i>
              <i>
                <input id="maxresults" type="number" value="20" min="1" max="100"
                  class="input input-bordered w-full max-w-xs md:max-w-md" v-model="maxResults" required />
              </i>
            </div>


          </form>
        </li>
        <li>
          <form @submit.prevent="" class="grid grid-rows-1">
            <div class="grid grid-rows-2 place-items-center">
              <i>
                <label for="maxresults" class="capitalize text-gray-500"> Fields</label>
              </i>
              <i>
                <select class="select select-bordered w-full max-w-xs md:max-w-md text-fourth" required>
                  <option disabled selected value="">Search field</option>
                  <option @click.prevent="changeSearchField(item)" v-for="(item, index) in searchFields" :key="index">
                    {{ item }}
                  </option>
                </select>
              </i>
            </div>
          </form>
        </li>
      </ul>
    </div>
  </div>



</template>

<script>
import SearchPage from './components/SearchPage.vue';
import MailsPage from './components/Email.vue'


export default {
  components: {
    SearchPage,
    MailsPage,
  },
  data() {
    return {
      maxResults: 20,
      selectedSearchField: "_all",
      currentPage: "Search",
      currentMail: null,
      pages: [
        { name: "Search", component: "SearchPage" },
        { name: "Read Mail", component: "MailsPage" },
      ],
      searchFields: [
        "_all",
        "To",
        "From",
        "Subject",
        "Cc",
        "Bcc",
      ]
    }
  },
  methods: {
    recordSearchType(key) {
      this.searchType = key
    },
    IsActive(key) {
      return (key === this.currentPage)
    },
    changeActivePage(key) {
      this.currentPage = key
    },
    displayActivePage(key) {
      return ["btn", "btn-ghost", "text-xl", this.IsActive(key) ? "text-fourth" : ""]
    },
    ReturnRow(data) {
      this.changeActivePage(this.pages[1].name)
      this.currentMail = data
    },
    changeSearchField(item) {
      this.selectedSearchField = item
    }
  },
}
</script>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}

::-webkit-scrollbar {
  width: 10px;
}

::-webkit-scrollbar-track {
  background: #F2F2F2;
}

::-webkit-scrollbar-thumb {
  background: #BDBDBD;
}

::-webkit-scrollbar-thumb:hover {
  background: #6E6E6E;
}
</style>