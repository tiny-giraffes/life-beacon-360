<template>
  <div>
    <h1>Tracking Page</h1>
    <!-- Add navigation button at the top -->
    <div class="mb-4">
      <v-btn color="primary" @click="goBack">Back to Settings</v-btn>
    </div>

    <p v-if="position">
      Current Position: Latitude: {{ position.latitude }}, Longitude:
      {{ position.longitude }}
    </p>
    <p v-else>Waiting for location...</p>
  </div>
</template>

<script lang="ts">
import { ref, onMounted, onBeforeUnmount } from "vue";
import { useRouter } from "vue-router";
import { registerPlugin } from "@capacitor/core";
import { Preferences } from "@capacitor/preferences";
import type { BackgroundGeolocationPlugin } from "@capacitor-community/background-geolocation";

const BackgroundGeolocation = registerPlugin<BackgroundGeolocationPlugin>(
  "BackgroundGeolocation"
);

export default {
  setup() {
    const router = useRouter();
    const position = ref<{ latitude: number; longitude: number } | null>(null);
    const serverAddress = ref("");
    const serverToken = ref("");
    const sendLocation = ref(false);
    const locationInterval = ref(0);
    let watcherId: string | null = null;
    let intervalId: number | null = null;

    const goBack = () => {
      router.push("/");
    };

    const loadSettings = async () => {
      const { value } = await Preferences.get({ key: "appSettings" });
      if (value) {
        const settings = JSON.parse(value);
        serverAddress.value = settings.serverAddress;
        serverToken.value = settings.serverToken;
        sendLocation.value = settings.sendLocation;
        locationInterval.value = settings.locationInterval;
      }
    };

    const sendToServer = async (location: {
      latitude: number;
      longitude: number;
    }) => {
      try {
        await fetch(serverAddress.value, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${serverToken.value}`,
          },
          body: JSON.stringify({
            latitude: location.latitude,
            longitude: location.longitude,
            timestamp: new Date().toISOString(),
          }),
        });
      } catch (error) {
        console.error("Failed to send location:", error);
      }
    };

    const startTracking = async () => {
      await loadSettings();

      watcherId = await BackgroundGeolocation.addWatcher(
        {
          backgroundMessage: "Cancel to prevent battery drain.",
          backgroundTitle: "Tracking You.",
          requestPermissions: true,
          stale: false,
          distanceFilter: 50,
        },
        (location, error) => {
          if (error) {
            if (error.code === "NOT_AUTHORIZED") {
              if (confirm("Location permission required. Open settings?")) {
                BackgroundGeolocation.openSettings();
              }
            }
            return;
          }

          if (location) {
            position.value = {
              latitude: location.latitude,
              longitude: location.longitude,
            };

            // Send immediately if interval is 0
            if (sendLocation.value && locationInterval.value === 0) {
              sendToServer(position.value);
            }
          }
        }
      );

      // Setup interval if enabled
      if (sendLocation.value && locationInterval.value > 0) {
        intervalId = window.setInterval(() => {
          if (position.value) sendToServer(position.value);
        }, locationInterval.value * 1000);
      }
    };

    const stopTracking = () => {
      if (watcherId) BackgroundGeolocation.removeWatcher({ id: watcherId });
      if (intervalId) clearInterval(intervalId);
    };

    onMounted(startTracking);
    onBeforeUnmount(stopTracking);

    return { position, goBack };
  },
};
</script>
