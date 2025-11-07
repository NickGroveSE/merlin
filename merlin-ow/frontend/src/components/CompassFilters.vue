<template>
  <div class="filters-container">
    <div class="filter-group">
      <label class="filter-label">Role</label>
      <select v-model="filters.role" class="filter-select">
        <option value="Tank">Tank</option>
        <option value="Damage">Damage</option>
        <option value="Support">Support</option>
        <!-- <option value="Custom">Custom</option> -->
      </select>
    </div>

    <div class="filter-group">
      <label class="filter-label">Input</label>
      <select v-model="filters.input" class="filter-select">
        <option value="PC">PC</option>
        <option value="Console">Console</option>
      </select>
    </div>

    <div class="filter-group">
      <label class="filter-label">Game Mode</label>
      <select v-model="filters.gameMode" class="filter-select" @change="onGameModeChange">
        <option value="0">Quick Play - Role Queue</option>
        <option value="2">Competitive - Role Queue</option>
      </select>
    </div>

    <!-- Conditional Rank Tier Dropdown -->
    <div v-if="filters.gameMode !== '0'" class="filter-group">
      <label class="filter-label">Rank Tier</label>
      <select v-model="filters.rankTier" class="filter-select">
        <option value="All">All Ranks</option>
        <option value="Bronze">Bronze</option>
        <option value="Silver">Silver</option>
        <option value="Gold">Gold</option>
        <option value="Platinum">Platinum</option>
        <option value="Diamond">Diamond</option>
        <option value="Master">Master</option>
        <option value="Grandmaster">Grandmaster & Champion</option>
      </select>
    </div>

    <div class="filter-group">
      <label class="filter-label">Map</label>
      <select v-model="filters.map" class="filter-select">
        <option value="all-maps">All Maps</option>
        <optgroup v-if="filters.gameMode === '0'" label="Clash">
          <option value="hanoka">Hanoka</option>
          <option value="throne-of-anubis">Throne of Anubis</option>
        </optgroup>
        <optgroup label="Control">
          <option value="antarctic-peninsula">Antarctic Peninsula</option>
          <option value="busan">Busan</option>
          <option value="ilios">Ilios</option>
          <option value="lijiang-tower">Lijiang Tower</option>
          <option value="nepal">Nepal</option>
          <option value="oasis">Oasis</option>
          <option value="samoa">Samoa</option>
        </optgroup>
        <optgroup label="Escort">
          <option value="circuit-royal">Circuit Royal</option>
          <option value="dorado">Dorado</option>
          <option value="havana">Havana</option>
          <option value="junkertown">Junkertown</option>
          <option value="rialto">Rialto</option>
          <option value="route-66">Route 66</option>
          <option value="shambali-monastery">Shambali Monastery</option>
          <option value="watchpoint-gibraltar">Watchpoint: Gilbraltar</option>
        </optgroup>
        <optgroup label="Flashpoint">
          <option value="aatlis">Aatlis</option>
          <option value="new-junk-city">New Junk City</option>
          <option value="suravasa">Suravasa</option>
        </optgroup>
        <optgroup label="Hybrid">
          <option value="blizzard-world">Blizzard World</option>
          <option value="eichenwalde">Eichenwalde</option>
          <option value="hollywood">Hollywood</option>
          <option value="kings-row">King's Row</option>
          <option value="midtown">Midtown</option>
          <option value="numbani">Numbani</option>
          <option value="paraiso">Paraíso</option>
        </optgroup>
        <optgroup label="Push">
          <option value="colosseo">Colosseo</option>
          <option value="esperanca">Esperança</option>
          <option value="new-queen-street">New Queen Street</option>
          <option value="runasapi">Runasapi</option>
        </optgroup>
      </select>
    </div>

    <div class="filter-group">
      <label class="filter-label">Region</label>
      <select v-model="filters.region" class="filter-select">
        <option value="Americas">Americas</option>
        <option value="Europe">Europe</option>
        <option value="Asia">Asia</option>
      </select>
    </div>

    <button @click="applyFilters" class="apply-filters-btn">
      Apply Filters
    </button>
  </div>
</template>

<script>

export default {
  name: 'OverwatchFilters',
  props: {
    queryParams: Object
  },
  data() {
    return {
      filters: {
        role: this.queryParams.role,
        input: this.queryParams.input,
        gameMode: this.queryParams.gameMode,
        rankTier: this.queryParams.rankTier, // Default rank tier
        map: this.queryParams.map,
        region: this.queryParams.region
      }
    }
  },
  watch: {
    queryParams: {
      handler(newParams) {
        // Update filters when queryParams prop changes
        this.filters = {
          role: newParams.role,
          input: newParams.input,
          gameMode: newParams.gameMode,
          rankTier: newParams.rankTier,
          map: newParams.map,
          region: newParams.region
        };
      },
      deep: true  // Watch nested properties
    }
  },
  methods: {
    onGameModeChange() {
      // Reset rank tier when switching to Quick Play
      if (this.filters.gameMode === '0') {
        this.filters.rankTier = 'All';
      }
    },
    applyFilters() {
        const filtersToSend = { ...this.filters };
        if (this.filters.gameMode === '0') {
            delete filtersToSend.rankTier;
        }

        // Emit the filters to parent component
        this.$emit('filters-applied', filtersToSend);
    }
  }
}
</script>

<style scoped>

.filters-container {
  background: #111827;
  border-radius: 12px;
  padding: 24px 24px 24px 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(31,41,55, 0.7);
  width: 300px;
}

/* Filter Group Styling */
.filter-group {
  margin-bottom: 15px;
}

.filter-label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: #e0e0e0;
  margin-bottom: 5px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* Dropdown Styling */
.filter-select {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 2px solid rgba(249, 168, 38, 0.3);
  border-radius: 8px;
  color: #ffffff;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg width='12' height='8' viewBox='0 0 12 8' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M1 1L6 6L11 1' stroke='%23f9a826' stroke-width='2' stroke-linecap='round'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 16px center;
  padding-right: 44px;
}

.filter-select:hover {
  border-color: rgba(249, 168, 38, 0.6);
  box-shadow: 0 0 15px rgba(249, 168, 38, 0.2);
}

.filter-select:focus {
  outline: none;
  border-color: #f9a826;
  box-shadow: 0 0 20px rgba(249, 168, 38, 0.3);
}

.filter-select optgroup {
  background: #1a1a2e;
  color: #06b6d4;
  padding: 10px;
}

.filter-select option {
  background: #1a1a2e;
  color: #ffffff;
  padding: 10px;
}


/* Apply Button */
.apply-filters-btn {
  width: 100%;
  height: min-content;
  padding: 10px 24px;
  margin: 20px auto 0 auto;
  background: linear-gradient(135deg, #f9a826 0%, #06b6d4 100%);
  border: none;
  border-radius: 8px;
  color: #ffffff;
  font-size: 16px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(249, 168, 38, 0.4);
}

.apply-filters-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 25px rgba(249, 168, 38, 0.6);
}

.apply-filters-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 10px rgba(249, 168, 38, 0.4);
}
</style>