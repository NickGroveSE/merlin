<script>
import CompassFilters from './components/CompassFilters.vue';
import HeroCompass from './components/HeroCompass.vue';
import MonitoringPopUp from './components/MonitoringPopUp.vue';
import { OverwatchService, CaptureService } from "../bindings/merlin-ow"
import { Events } from '@wailsio/runtime';

export default {
  name: 'App',
  components: {
    HeroCompass,
    CompassFilters,
    MonitoringPopUp
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
      loading: false,
      monitoring: false,
      status: {
        icon: "../public/assets/idle.svg",
        statusText: "Idle",
        gameData: {
          input: 'PC',
          queue: '',
          rank: 'All Ranks',
          role: '',
          map: 'All Maps',
          region: 'Americas'
        },
        message: "Messages with more in-depth status updates..."
      }
    }
  },
  mounted() {
    // Call scrape with default filters when component loads
    this.handleFiltersApplied(this.defaultFilters)

    Events.On('status-update', (data) => {
      this.status.icon = data.data[0].statusIcon
      this.status.statusText = data.data[0].statusText
      this.status.message = data.data[0].message

      // console.log("Made It Here")
      // console.log(data.data[0].statusText)
    });

    Events.On('message', (data) => {
      this.status.message = data.data[0]

      // console.log(message)
    })

    Events.On('queue-update', (data) => {
      this.status.gameData.queue = data.data[0]

      console.log(data)
    });

    Events.On('role-update', (data) => {
      this.status.gameData.role = data.data[0]
    });

    Events.On('map-update', (data) => {
      this.status.gameData.map = data.data[0]
    });



    // Events.On('test-emit', () => {
    //   console.log("Here")
    // })
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
    async handleStartMonitoring() {
      console.log('Begin Capturing Screen');
      
      this.loading = true;
      this.monitoring = true;
      try {
        const result = await CaptureService.StartMonitoring();
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
        this.monitoring = false;
        this.loading = false;
      }
    },
    async handleStopMonitoring() {
      console.log('Stop Capturing Screen');

      const result = await CaptureService.StopMonitoring();
    }
  }
}
</script>

<template>
  <div class="container" :class="{'disable-click': loading}">
    <HeroCompass v-if="heroData" class="component" id="compass" :heroData="heroData" @start-monitoring="handleStartMonitoring"/>
    <CompassFilters class="component" :queryParams="defaultFilters" @filters-applied="handleFiltersApplied"/>
    <div v-if="monitoring" id="monitor-wrapper">
      <MonitoringPopUp :statusIcon="status.icon" :statusText="status.statusText" :gameData="status.gameData" :message="status.message" @stop-monitoring="handleStopMonitoring"/>
    </div>
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

  .disable-click {
    position: relative;
  }

  .disable-click::before {
    pointer-events: none;
    content: ''; /* Required for pseudo-elements */
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.2); /* Translucent black overlay */
    z-index: 1; /* Place the overlay behind the content if needed */
  }

  #monitor-wrapper {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 2;
  }
</style>
