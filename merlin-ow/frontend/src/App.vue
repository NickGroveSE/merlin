<script>
import CompassFilters from './components/CompassFilters.vue';
import HeroCompass from './components/HeroCompass.vue';
import { OverwatchService, CaptureService } from "../bindings/changeme"

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
    },
    async handleCaptureTrigger() {
      console.log('Begin Capturing Screen');
      
      this.loading = true;
      try {
        const result = await CaptureService.Capture();
        console.log('Service Result:', result);
        this.heroData = result[0];
        this.defaultFilters = {
          role: result[1].role,
          input: result[1].input,
          gameMode: result[1].gameMode,
          rankTier: result[1].rankTier,
          map: result[1].map,
          region: result[1].region
        };
        console.log(this.defaultFilters)
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
    <HeroCompass v-if="heroData" class="component" id="compass" :heroData="heroData" @capture-triggered="handleCaptureTrigger"/>
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
    margin: 0 3rem 3rem 0;
  }
</style>
