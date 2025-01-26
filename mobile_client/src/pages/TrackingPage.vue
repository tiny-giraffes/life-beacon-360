<template>
  <div>
    <h1>Tracking Page</h1>
    <p v-if="position">
      Current Position: Latitude: {{ position.latitude }}, Longitude:
      {{ position.longitude }}
    </p>
    <p v-else>Waiting for location...</p>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import { registerPlugin } from "@capacitor/core";
import type { BackgroundGeolocationPlugin } from "@capacitor-community/background-geolocation";

// Register the plugin with the correct type
const BackgroundGeolocation = registerPlugin<BackgroundGeolocationPlugin>(
  "BackgroundGeolocation"
);

export default {
  setup() {
    const position = ref<{ latitude: number; longitude: number } | null>(null);
    let watcherId: string | null = null;

    const startTracking = async () => {
      try {
        // Start location tracking with a watcher
        watcherId = await BackgroundGeolocation.addWatcher(
          {
            backgroundMessage: "Cancel to prevent battery drain.",
            backgroundTitle: "Tracking You.",
            requestPermissions: true,
            stale: false,
            distanceFilter: 50, // Minimum distance in meters for updates
          },
          (location, error) => {
            if (error) {
              // Handle permission-related errors
              if (error.code === "NOT_AUTHORIZED") {
                if (
                  window.confirm(
                    "This app needs your location but does not have permission.\n\nOpen settings now?"
                  )
                ) {
                  BackgroundGeolocation.openSettings();
                }
              }
              console.error("Background geolocation error:", error);
              return;
            }

            // Update position state with the new location
            if (location) {
              position.value = {
                latitude: location.latitude,
                longitude: location.longitude,
              };
              console.log("Location updated:", position.value);
            }
          }
        );
      } catch (err) {
        console.error("Failed to start tracking:", err);
      }
    };

    const stopTracking = () => {
      if (watcherId) {
        BackgroundGeolocation.removeWatcher({ id: watcherId }).catch((err) => {
          console.error("Failed to remove watcher:", err);
        });
      }
    };

    onMounted(startTracking);
    onBeforeUnmount(stopTracking);

    return { position };
  },
};
</script>
