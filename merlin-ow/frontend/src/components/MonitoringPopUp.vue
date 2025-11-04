<template>
  <div class="status-container">
    <!-- Header Section -->
    <div class="status-header">
      <img class="status-icon" :src="statusIcon">
      <h2 class="status-text">Status : {{ statusText }}</h2>
      <button class="close-button" @click="stopMonitoring">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <!-- Content Section -->
    <div class="content-wrapper">
      <!-- Left Side - Game Data -->
      <div class="game-data">
        <div class="data-row">
          <span class="data-label">Input</span>
          <span class="data-separator">:</span>
          <span class="data-value">{{ gameData.input }}</span>
        </div>
        <div class="data-row">
          <span class="data-label">Queue</span>
          <span class="data-separator">:</span>
          <span class="data-value">{{ gameData.queue }}</span>
        </div>
        <div class="data-row">
          <span class="data-label">Rank</span>
          <span class="data-separator">:</span>
          <span class="data-value">{{ gameData.rank }}</span>
        </div>
        <div class="data-row">
          <span class="data-label">Role</span>
          <span class="data-separator">:</span>
          <span class="data-value">{{ gameData.role }}</span>
        </div>
        <div class="data-row">
          <span class="data-label">Map</span>
          <span class="data-separator">:</span>
          <span class="data-value">{{ gameData.map }}</span>
        </div>
        <div class="data-row">
          <span class="data-label">Region</span>
          <span class="data-separator">:</span>
          <span class="data-value">{{ gameData.region }}</span>
        </div>
      </div>

      <!-- Right Side - Message Log -->
      <div class="message-box">
        <p class="message-description">
          {{message}}
        </p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MonitoringComponent',
  props: {
    statusText: {
      type: String,
      default: 'Idle'
    },
    statusIcon: {
      type: String,
      default: ''
    },
    gameData: {
      type: Object,
      default: () => ({
        input: 'PC',
        queue: '',
        rank: 'All Ranks',
        role: '',
        map: 'All Maps',
        region: 'Americas'
      })
    },
    message: {
      type: String,
      default: "Messages with more in-depth status updates..."
    }
  },
  methods: {
    stopMonitoring() {
        this.$emit('stop-monitoring')
    }
  }
}
</script>

<style scoped>
.status-container {
  width: 600px;
  margin: 150px auto 0 auto;
  border: 2px solid rgba(249, 168, 38, 0.3);
  border-radius: 8px;
  background: #111827;
  transition: all 0.3s ease;
  /* font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; */
}

.status-header {
  display: flex;
  align-items: center;
  padding: 20px 24px;
  /* border-bottom: 2px solid rgba(249, 168, 38, 1); */
  gap: 16px;
}

.status-icon {
  width: auto;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transform: translatey(2px);
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
	0% {
		transform: translatey(2px);
	}
	50% {
		transform: translatey(-2px);
	}
	100% {
		transform: translatey(2px);
	}
}

.status-text {
  flex: 1;
  font-size: 32px;
  text-align: left;
  margin: 0 0 0 20px;
}

.close-button {
  width: 32px;
  height: 32px;
  border: 2px solid #333;
  background: #B20000;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.close-button:hover {
  background: #E50000;
}

.content-wrapper {
  display: grid;
  grid-template-columns: 1fr 1fr;
  border-top: 2px solid rgba(249, 168, 38, 0.3);
}

.game-data {
  padding: 24px;
  border-right: 2px solid rgba(249, 168, 38, 0.3);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.data-row {
  display: grid;
  grid-template-columns: 80px 20px 1fr;
  gap: 8px;
  font-size: 16px;
}

.data-label {
  font-weight: 600;
  text-align: left;
}

.data-separator {
  text-align: center;
}

.data-value {
  text-align: left;
}

.message-box {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message-description {
  font-size: 14px;
  line-height: 1.5;
  margin: 0;
}

.message-log {
  flex: 1;
  background: #f9f9f9;
  border: 1px solid #ddd;
  border-radius: 4px;
  padding: 12px;
  max-height: 200px;
  overflow-y: auto;
  font-family: monospace;
  font-size: 12px;
}

.log-entry {
  display: block;
  margin-bottom: 8px;
  line-height: 1.4;
}

.log-entry:last-child {
  margin-bottom: 0;
}

.log-timestamp {
  color: #666;
  margin-right: 8px;
}

.log-text {
  color: #333;
}
</style>