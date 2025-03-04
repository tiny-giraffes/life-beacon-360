<!--
 * Life Beacon 360
 * Copyright (C) 2025 Tim Yashin/tiny-giraffes
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 -->

<template>
  <v-container class="py-4">
    <h1 class="text-center mb-4">Location Tracking</h1>

    <v-row>
      <v-col cols="12" class="mb-4">
        <v-btn color="primary" block @click="goBack">Back to Settings</v-btn>
      </v-col>

      <v-col cols="12" class="text-center mb-4">
        <v-card class="pa-4">
          <div v-if="position">
            <h3 class="mb-2">Current Position</h3>
            <p>Latitude: {{ position.latitude.toFixed(6) }}</p>
            <p>Longitude: {{ position.longitude.toFixed(6) }}</p>
            <p v-if="lastSentTime">Last sent: {{ lastSentTime }}</p>
          </div>
          <div v-else>
            <p>Waiting for location data...</p>
          </div>
        </v-card>
      </v-col>

      <v-col cols="12" class="text-center">
        <v-btn
          v-if="!isTracking"
          color="success"
          size="large"
          block
          @click="startTracking"
        >
          Start Tracking
        </v-btn>
        <v-btn v-else color="error" size="large" block @click="stopTracking">
          Stop Tracking
        </v-btn>
      </v-col>

      <v-col cols="12" v-if="isTracking && sendLocation">
        <v-card class="pa-4 mt-4">
          <h3 class="mb-2">Server Status</h3>
          <p>Sending to: {{ serverAddress }}</p>
          <p>Interval: {{ locationInterval }} seconds</p>
          <p v-if="serverStatus">{{ serverStatus }}</p>
        </v-card>
      </v-col>

      <v-col cols="12" v-if="error">
        <v-alert type="error" class="mt-4">
          {{ error }}
        </v-alert>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, onBeforeUnmount } from "vue";
import { useRouter } from "vue-router";
import { registerPlugin } from "@capacitor/core";
import { Preferences } from "@capacitor/preferences";
import type { BackgroundGeolocationPlugin } from "@capacitor-community/background-geolocation";

const BackgroundGeolocation = registerPlugin<BackgroundGeolocationPlugin>(
  "BackgroundGeolocation"
);

export default defineComponent({
  setup() {
    const router = useRouter();
    const position = ref<{ latitude: number; longitude: number } | null>(null);
    const serverAddress = ref("");
    const serverToken = ref("");
    const sendLocation = ref(false);
    const locationInterval = ref(0);
    const isTracking = ref(false);
    const watcherId = ref<string | null>(null);
    const intervalId = ref<number | null>(null);
    const serverStatus = ref("");
    const error = ref("");
    const lastSentTime = ref("");

    const goBack = () => {
      if (isTracking.value) {
        if (confirm("Tracking is active. Stop tracking and go back?")) {
          stopTracking();
          router.push("/");
        }
      } else {
        router.push("/");
      }
    };

    const loadSettings = async () => {
      try {
        const { value } = await Preferences.get({ key: "appSettings" });
        if (value) {
          const settings = JSON.parse(value);
          serverAddress.value = settings.serverAddress || "";
          serverToken.value = settings.serverToken || "";
          sendLocation.value = settings.sendLocation || false;
          locationInterval.value = settings.locationInterval || 0;
        }
      } catch (err) {
        error.value = "Failed to load settings";
        console.error("Failed to load settings:", err);
      }
    };

    const formatTime = () => {
      const now = new Date();
      return now.toLocaleTimeString();
    };

    const sendToServer = async (location: {
      latitude: number;
      longitude: number;
    }) => {
      if (!sendLocation.value || !serverAddress.value) return;

      try {
        serverStatus.value = "Sending location...";
        console.log("Sending to address:", serverAddress.value);
        console.log("Using token:", serverToken.value);
        console.log(
          "Sending data:",
          JSON.stringify({
            latitude: location.latitude,
            longitude: location.longitude,
          })
        );

        const response = await fetch(serverAddress.value, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            // Try without Bearer prefix since your curl might be working that way
            Authorization: serverToken.value,
          },
          body: JSON.stringify({
            latitude: location.latitude,
            longitude: location.longitude,
          }),
        });

        const responseText = await response.text();
        console.log("Response status:", response.status);
        console.log("Response body:", responseText);

        if (response.ok) {
          serverStatus.value = "Location sent successfully";
          lastSentTime.value = formatTime();
        } else {
          serverStatus.value = `Error: ${response.status} - ${responseText}`;
        }
      } catch (err) {
        console.error("Send error details:", err);
      }
    };

    const startTracking = async () => {
      try {
        await loadSettings();

        // Request location permissions
        watcherId.value = await BackgroundGeolocation.addWatcher(
          {
            backgroundMessage: "Life Beacon 360 is tracking your location.",
            backgroundTitle: "Location Tracking Active",
            requestPermissions: true,
            stale: false,
            distanceFilter: 10, // Meters
          },
          (location, err) => {
            if (err) {
              error.value = `Location error: ${err.message}`;
              if (err.code === "NOT_AUTHORIZED") {
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

        // Set up interval for sending location if enabled and interval > 0
        if (sendLocation.value && locationInterval.value > 0) {
          intervalId.value = window.setInterval(() => {
            if (position.value) {
              sendToServer(position.value);
            }
          }, locationInterval.value * 1000);
        }

        isTracking.value = true;
        error.value = "";
      } catch (err) {
        error.value = "Failed to start tracking";
        console.error("Failed to start tracking:", err);
      }
    };

    const stopTracking = async () => {
      try {
        // Stop location watcher
        if (watcherId.value) {
          await BackgroundGeolocation.removeWatcher({
            id: watcherId.value,
          });
          watcherId.value = null;
        }

        // Clear sending interval
        if (intervalId.value !== null) {
          clearInterval(intervalId.value);
          intervalId.value = null;
        }

        isTracking.value = false;
        serverStatus.value = "";
      } catch (err) {
        error.value = "Failed to stop tracking";
        console.error("Failed to stop tracking:", err);
      }
    };

    // Clean up on unmount
    onBeforeUnmount(() => {
      if (isTracking.value) {
        stopTracking();
      }
    });

    // Load settings when component mounts
    onMounted(() => {
      loadSettings();
    });

    return {
      position,
      isTracking,
      serverAddress,
      locationInterval,
      sendLocation,
      serverStatus,
      error,
      lastSentTime,
      goBack,
      startTracking,
      stopTracking,
    };
  },
});
</script>

<style scoped>
.v-container {
  max-width: 600px;
  margin: auto;
}
</style>
