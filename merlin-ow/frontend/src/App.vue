<script>
import CompassFilters from './components/CompassFilters.vue';
import HeroCompass from './components/HeroCompass.vue';
import { OverwatchService } from "../bindings/changeme"

export default {
  name: 'App',
  components: {
    HeroCompass,
    CompassFilters,
  },
  data() {
    return {
      defaultFilters: {
        role: 'Support',
        input: 'PC',
        gameMode: '2',
        rankTier: 'All',
        map: 'all-maps',
        region: 'Americas'
      },
      heroData: null,
      loading: false
    }
  },
  mounted() {
    // Call scrape with default filters when component loads
    this.handleFiltersApplied(this.defaultFilters)
  },
  methods: {
    async handleFiltersApplied(filters) {
      console.log('Filters received in App.vue:', filters);
      
      this.loading = true;
      try {
        const result = await OverwatchService.Scrape(filters);
        // console.log('Scrape Result:', result);
        this.heroData = result;
      } catch (err) {
        console.error('Error Scraping:', err);
      } finally {
        this.loading = false;
      }
    }
  }
}
</script>

<template>
  <div class="container">
    <HeroCompass v-if="heroData" class="component" id="compass" :heroData="heroData"/>
    <CompassFilters class="component" :queryParams="defaultFilters" @filters-applied="handleFiltersApplied"/>
  </div>
</template>

<style scoped>
  .container {
    padding: 2rem;
    width: calc(100% - 4rem);
    height: calc(100% - 4rem);
    margin: 0 auto;
    color: white;
    display: inline-block;
    vertical-align: top;
  }

  .component {
    display: inline-block;
    vertical-align: top;
  }

  #compass {
    margin: 0 3rem 0 0;
  }
</style>
