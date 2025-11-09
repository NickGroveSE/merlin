<template>
  <div class="chart-container">
    <div class="chart-wrapper">
      <div class="axis-label label-top">Strong</div>
      <div class="axis-label label-bottom">Weak</div>
      <div class="axis-label label-left">Niche</div>
      <div class="axis-label label-right">Popular</div>
      
      <svg viewBox="0 0 560 500" class="chart-svg" preserveAspectRatio="xMidYMid meet">
        <!-- Grid lines -->
        <line
          v-for="i in 7"
          :key="`h-${i}`"
          :x1="margin.left"
          :y1="margin.top + (height / 6) * (i - 1)"
          :x2="margin.left + width"
          :y2="margin.top + (height / 6) * (i - 1)"
          class="grid-line"
        />
        <line
          v-for="i in 7"
          :key="`v-${i}`"
          :x1="margin.left + (width / 6) * (i - 1)"
          :y1="margin.top"
          :x2="margin.left + (width / 6) * (i - 1)"
          :y2="margin.top + height"
          class="grid-line"
        />

        <!-- Center axes -->
        <line
          :x1="margin.left"
          :y1="centerY"
          :x2="margin.left + width"
          :y2="centerY"
          class="axis-line"
        />
        <line
          :x1="centerX"
          :y1="margin.top"
          :x2="centerX"
          :y2="margin.top + height"
          class="axis-line"
        />

        <!-- Data points -->
        <g v-for="(hero, i) in data" :key="hero.Name">
          <text
            v-if="tooltip.pinnedHero !== hero.Name && (tooltip.pinned || hoveredHero !== hero.Name)"
            :x="getLabelPosition(hero, i).x"
            :y="getLabelPosition(hero, i).y"
            text-anchor="middle"
            fill="#fff"
            font-size="12"
            font-weight="600"
            class="hero-dot-text"
          >
            {{ hero.Name }}
          </text>
          <circle
            :cx="scaleX(hero.PickRate)"
            :cy="scaleY(hero.WinRate)"
            r="10"
            :fill="hero.Color"
            opacity="0.75"
            :class="{ 'hero-dot': true, 'pinned': tooltip.pinned && tooltip.pinnedHero === hero.Name }"
            @mouseenter="handleMouseEnter(hero, $event)"
            @mousemove="handleMouseMove($event)"
            @mouseleave="handleMouseLeave"
            @click="handleDotClick(hero, $event)"
          />
          <circle
            v-if="tooltip.pinned && tooltip.pinnedHero === hero.Name || hoveredHero === hero.Name"
            :cx="scaleX(hero.PickRate)"
            :cy="scaleY(hero.WinRate)"
            r="10"
            fill="none"
            :stroke="hero.Color"
            opacity="1"
            stroke-width="3"
            class="pulse-ring"
          />

        </g>
      </svg>

      <div
        v-if="tooltip.show"
        class="tooltip"
        :style="{
          left: `${tooltip.x + 5}px`,
          top: `${tooltip.y - 5}px`
        }"
      >
        <button 
          v-if="tooltip.pinned"
          class="tooltip-close"
          @click="closeTooltip"
        >
          ×
        </button>
        <img 
          class="tooltip-image"
          :src="tooltip.data.Image"
        />
        <div class="tooltip-name">{{ tooltip.data.Name }}</div>
        <div class="tooltip-stat">Pick Rate: {{ tooltip.data.PickRate }}%</div>
        <div class="tooltip-stat">Win Rate: {{ tooltip.data.WinRate }}%</div>
      </div>
    </div>
    <div id="capture-btn" @click="startMonitoring">Monitor Screen</div>
  </div>
</template>

<script>

export default {
  name: 'OverwatchScatterPlot',
  props: {
    heroData: Object,
    error: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      data: this.heroData,
      margin: { top: 40, right: 60, bottom: 30, left: 60 },
      hoveredHero: null, 
      tooltip: {
        show: false,
        pinned: false,
        pinnedHero: null,
        x: 0,
        y: 0,
        data: null
      },
      labelPositions: []
    };
  },
  watch: {
    heroData: {
      handler(newVal) {
        this.data = newVal;
      },
      deep: true,
      immediate: true
    },
    data: {
      handler() {
        // Reset tooltip whenever data changes
        this.tooltip = {
          show: false,
          pinned: false,
          pinnedHero: null,
          x: 0,
          y: 0,
          data: null
        };
        this.hoveredHero = null;
        if (this.data && this.data.length > 0) {
          this.$nextTick(() => {
            this.calculateLabelPositions();
          });
        }
      },
      deep: true
    }
  },
  mounted() {
    this.$el.addEventListener('hero-data-updated', (event) => {
      // This updates internal data, not the prop
      this.data = event.detail;
    });

    // Wait for data to be fully set before calculating positions
    this.$nextTick(() => {
      if (this.data && this.data.length > 0) {
        this.calculateLabelPositions();
      }
    });
  },
  computed: {
    width() {
      return 560 - this.margin.left - this.margin.right;
    },
    height() {
      return 500 - this.margin.top - this.margin.bottom;
    },
    pickRates() {
      return this.data.map(d => d.PickRate);
    },
    winRates() {
      return this.data.map(d => d.WinRate);
    },
    minPickRate() {
      return Math.min(...this.pickRates);
    },
    maxPickRate() {
      return Math.max(...this.pickRates);
    },
    minWinRate() {
      return Math.min(...this.winRates);
    },
    maxWinRate() {
      return Math.max(...this.winRates);
    },
    centerPickRate() {
      return (this.minPickRate + this.maxPickRate) / 2;
    },
    centerWinRate() {
      return (this.minWinRate + this.maxWinRate) / 2;
    },
    pickRateRange() {
      return this.maxPickRate - this.minPickRate;
    },
    winRateRange() {
      return this.maxWinRate - this.minWinRate;
    },
    paddedMinPickRate() {
      return this.minPickRate - this.pickRateRange * 0.1;
    },
    paddedMaxPickRate() {
      return this.maxPickRate + this.pickRateRange * 0.1;
    },
    paddedMinWinRate() {
      return this.minWinRate - this.winRateRange * 0.1;
    },
    paddedMaxWinRate() {
      return this.maxWinRate + this.winRateRange * 0.1;
    },
    centerX() {
      return this.scaleX(this.centerPickRate);
    },
    centerY() {
      return this.scaleY(this.centerWinRate);
    }
  },
  methods: {
    scaleX(value) {
      return this.margin.left + (value - this.paddedMinPickRate) / 
             (this.paddedMaxPickRate - this.paddedMinPickRate) * this.width;
    },
    scaleY(value) {
      return this.margin.top + this.height - (value - this.paddedMinWinRate) / 
             (this.paddedMaxWinRate - this.paddedMinWinRate) * this.height;
    },
    calculateLabelPositions() {
      this.labelPositions = [];

      const positions = [];
      const labelHeight = 16; // Approximate height of label with padding
      
      this.data.forEach((hero, i) => {
        const x = this.scaleX(hero.PickRate);
        const y = this.scaleY(hero.WinRate);
        
        // Estimate label width based on character count (roughly 7px per character)
        const nameLength = hero.Name.length;
        const estimatedWidth = nameLength * 7;
        const halfWidth = estimatedWidth / 2;
        
        // Possible offsets to try (top, bottom, left, right, diagonal)
        // Adjust horizontal offsets based on name length
        const offsets = [
          { x: 0, y: -12 },                    // top (default)
          { x: 0, y: 20 },                     // bottom
          { x: -halfWidth - 10, y: 4 },        // left (adjusted)
          { x: halfWidth + 10, y: 4 },         // right (adjusted)
          { x: -halfWidth - 5, y: -10 },       // top-left (adjusted)
          { x: halfWidth + 5, y: -10 },        // top-right (adjusted)
          { x: -halfWidth - 5, y: 20 },        // bottom-left (adjusted)
          { x: halfWidth + 5, y: 20 }          // bottom-right (adjusted)
        ];
        
        let bestOffset = offsets[0];
        let bestScore = -Infinity;
        
        // Try each offset and score based on distance from other labels AND dots
        for (const offset of offsets) {
          const testX = x + offset.x;
          const testY = y + offset.y;
          
          let minLabelDist = Infinity;
          let minDotDist = Infinity;
          
          // Check distance to all existing labels (rectangular collision)
          for (let j = 0; j < positions.length; j++) {
            const pos = positions[j];
            const otherHero = this.data[j];
            const otherWidth = otherHero.Name.length * 7;
            
            // Calculate rectangular overlap
            const xDist = Math.abs(testX - pos.x);
            const yDist = Math.abs(testY - pos.y);
            const widthSum = (estimatedWidth + otherWidth) / 2 + 10;
            
            // Penalize if labels would overlap
            if (xDist < widthSum && yDist < labelHeight) {
              minLabelDist = Math.min(minLabelDist, 0);
            } else {
              const dist = Math.sqrt(Math.pow(xDist, 2) + Math.pow(yDist, 2));
              minLabelDist = Math.min(minLabelDist, dist);
            }
          }
          
          // Check distance to all dots (avoid covering other dots)
          for (let j = 0; j < this.data.length; j++) {
            if (j === i) continue; // Skip own dot
            
            const otherHero = this.data[j];
            const dotX = this.scaleX(otherHero.PickRate);
            const dotY = this.scaleY(otherHero.WinRate);
            
            // Check if label would cover the dot
            const xDist = Math.abs(testX - dotX);
            const yDist = Math.abs(testY - dotY);
            
            if (xDist < halfWidth && yDist < labelHeight / 2) {
              minDotDist = Math.min(minDotDist, 0); // Label covers dot
            } else {
              const dist = Math.sqrt(Math.pow(xDist, 2) + Math.pow(yDist, 2));
              minDotDist = Math.min(minDotDist, dist);
            }
          }
          
          // Score this position (prioritize avoiding overlaps, then maximize distance)
          const score = minLabelDist * 2 + minDotDist;
          
          // If this position is better, use it
          if (score > bestScore) {
            bestScore = score;
            bestOffset = offset;
          }
        }
        
        positions.push({
          x: x + bestOffset.x,
          y: y + bestOffset.y,
          width: estimatedWidth
        });
      });
      
      this.labelPositions = positions;
    },
    getLabelPosition(hero, index) {
      if (this.labelPositions[index]) {
        return this.labelPositions[index];
      }
      // Fallback to default position
      return {
        x: this.scaleX(hero.PickRate),
        y: this.scaleY(hero.WinRate) - 12
      };
    },
    handleMouseEnter(hero, event) {
      this.hoveredHero = hero.Name;
      // Don't show tooltip on hover if it's pinned to a different hero
      if (this.tooltip.pinned && this.tooltip.pinnedHero !== hero.Name) {
        return;
      }
      
      this.tooltip = {
        show: true,
        pinned: this.tooltip.pinned,
        pinnedHero: this.tooltip.pinnedHero,
        x: event.pageX,
        y: event.pageY,
        data: hero
      };
    },
    handleMouseMove(event) {
      // Don't update position if tooltip is pinned
      if (this.tooltip.show && !this.tooltip.pinned) {
        this.tooltip.x = event.pageX;
        this.tooltip.y = event.pageY;
      }
    },
    handleMouseLeave() {
      this.hoveredHero = null;
      // Don't hide tooltip if it's pinned
      if (!this.tooltip.pinned) {
        this.tooltip = {
          show: false,
          pinned: false,
          pinnedHero: null,
          x: 0,
          y: 0,
          data: null
        };
      }
    },
    handleDotClick(hero, event) {
      event.stopPropagation();
      
      // If clicking the same hero that's pinned, unpin it
      if (this.tooltip.pinned && this.tooltip.pinnedHero === hero.Name) {
        this.tooltip = {
          show: false,
          pinned: false,
          pinnedHero: null,
          x: 0,
          y: 0,
          data: null
        };
      } else {
        // Pin the tooltip to this hero
        this.tooltip = {
          show: true,
          pinned: true,
          pinnedHero: hero.Name,
          x: event.pageX,
          y: event.pageY,
          data: hero
        };
      }
    },
    closeTooltip() {
      this.tooltip = {
        show: false,
        pinned: false,
        pinnedHero: null,
        x: 0,
        y: 0,
        data: null
      };
      this.hoveredHero = null;
    },
    startMonitoring() {
        // Emit to parent component
        this.$emit('start-monitoring');
    }
  }
};
</script>

<style scoped>
.chart-container {
  width: 100%;
  max-width: 560px;
  margin: 0 auto;
  overflow: visible;
}

.chart-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 560 / 500;
  border-radius: 10px;
  overflow: visible;
}

.axis-label {
  position: absolute;
  font-size: 20px;
  font-weight: 600;
  color: #fbbf24;
  padding: 4px 8px;
  border-radius: 4px;
  white-space: nowrap;
}

.label-top {
  top: 0;
  left: 50%;
  transform: translateX(-50%);
}

.label-bottom {
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
}

.label-left {
  top: 50%;
  left: -10px;
  transform: translateY(-50%);
}

.label-right {
  top: 50%;
  right: -30px;
  transform: translateY(-50%);
}

.chart-svg {
  width: 100%;
  height: 100%;
}

.grid-line {
  stroke: rgba(255, 255, 255, 0.1);
  stroke-width: 1;
}

.axis-line {
  stroke: rgba(255, 255, 255, 0.5);
  stroke-width: 2;
}

.hero-dot-text {
  text-shadow:
    -1px -1px 0 #111827,
    1px -1px 0 #111827,
    -1px 1px 0 #111827,
    1px 1px 0 #111827;
  transition: all 0.3s ease;
}

.hero-dot {
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.hero-dot:hover {
  opacity: 1 !important;
}

.hero-dot.pinned {
  opacity: 1 !important;
}

.tooltip {
  position: fixed;
  background: rgba(0, 0, 0, 0.9);
  color: #fff;
  padding: 10px 15px;
  border-radius: 8px;
  font-size: 13px;
  pointer-events: auto;
  z-index: 100;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.tooltip-close {
  position: absolute;
  top: 5px;
  right: 5px;
  background: none;
  border: none;
  color: #fff;
  font-size: 24px;
  line-height: 1;
  cursor: pointer;
  padding: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.2s ease;
}

.tooltip-close:hover {
  color: #fbbf24;
}

.tooltip-image {
  display: block;
  width: 50px;
  height: 50px;
  object-fit: contain;
  margin-bottom: 5px;
}

.tooltip-name {
  font-weight: bold;
  font-size: 15px;
  margin-bottom: 5px;
  color: #fbbf24;
}

.tooltip-stat {
  margin: 3px 0;
}

.pulse-ring {
  animation: pulse-ring 2s ease-out infinite;
}

@keyframes pulse-ring {
  0% {
    r: 10;
    opacity: 1;
  }
  100% {
    r: 18;
    opacity: 0;
  }
}

#capture-btn {
  width: 180px;
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

#capture-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 25px rgba(249, 168, 38, 0.6);
}

#capture-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 10px rgba(249, 168, 38, 0.4);
}
</style>