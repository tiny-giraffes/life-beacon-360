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
    <h1 class="text-center mb-4">Life Beacon 360</h1>
    <h2 class="text-center mb-4">Settings</h2>

    <v-card class="pa-4 mb-4">
      <v-form @submit.prevent="saveSettings">
        <v-row>
          <v-col cols="12">
            <v-text-field
              v-model="serverAddress"
              label="Server Address"
              hint="Enter the complete server URL (e.g., http://your-server:8080/api/locations)"
              persistent-hint
              outlined
              required
            />
          </v-col>

          <v-col cols="12">
            <v-text-field
              v-model="serverToken"
              label="Server Token"
              hint="Token for server authentication"
              persistent-hint
              outlined
              required
            />
          </v-col>

          <v-col cols="12">
            <v-checkbox
              v-model="sendLocation"
              label="Send Location to Server"
              hint="Enable to send location data to the server"
              persistent-hint
            />
          </v-col>

          <v-col cols="12" v-if="sendLocation">
            <v-text-field
              v-model.number="locationInterval"
              label="Location Interval (seconds)"
              hint="How often to send location updates (0 = immediate/real-time)"
              type="number"
              min="0"
              outlined
              required
            />
          </v-col>

          <v-col cols="12" class="text-center">
            <v-btn color="primary" size="large" type="submit" :loading="saving">
              Save Settings
            </v-btn>
          </v-col>
        </v-row>
      </v-form>
    </v-card>

    <v-card class="pa-4">
      <v-col cols="12" class="text-center">
        <p class="mb-2">Start location tracking:</p>
        <v-btn
          color="success"
          size="large"
          block
          @click="goToTrackingPage"
          :disabled="!areSettingsComplete"
        >
          Go to Tracking Page
        </v-btn>
        <p v-if="!areSettingsComplete" class="mt-2 text-error">
          Please complete and save settings first
        </p>
      </v-col>
    </v-card>

    <v-snackbar v-model="showSnackbar" :color="snackbarColor">
      {{ snackbarText }}
    </v-snackbar>
  </v-container>
</template>

<script lang="ts">
import { defineComponent, ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { Preferences } from "@capacitor/preferences";

export default defineComponent({
  name: "SettingsPage",
  setup() {
    // Router for navigation
    const router = useRouter();

    // Reactive states for settings
    const serverAddress = ref<string>("");
    const locationInterval = ref<number>(30);
    const serverToken = ref<string>("");
    const sendLocation = ref<boolean>(false);

    // UI states
    const saving = ref<boolean>(false);
    const showSnackbar = ref<boolean>(false);
    const snackbarText = ref<string>("");
    const snackbarColor = ref<string>("success");

    // Check if settings are complete enough to proceed
    const areSettingsComplete = computed(() => {
      if (!sendLocation.value) {
        return true; // If not sending location, no need for other settings
      }
      return !!serverAddress.value && !!serverToken.value;
    });

    // Save settings to local storage
    const saveSettings = async () => {
      if (sendLocation.value && (!serverAddress.value || !serverToken.value)) {
        showSnackbarMessage("Please fill in all required fields", "error");
        return;
      }

      saving.value = true;
      try {
        const settings = {
          serverAddress: serverAddress.value,
          locationInterval: locationInterval.value,
          serverToken: serverToken.value,
          sendLocation: sendLocation.value,
        };

        await Preferences.set({
          key: "appSettings",
          value: JSON.stringify(settings),
        });

        showSnackbarMessage("Settings saved successfully!", "success");
      } catch (error) {
        console.error("Failed to save settings:", error);
        showSnackbarMessage("Failed to save settings", "error");
      } finally {
        saving.value = false;
      }
    };

    // Show snackbar message
    const showSnackbarMessage = (
      message: string,
      color: string = "success"
    ) => {
      snackbarText.value = message;
      snackbarColor.value = color;
      showSnackbar.value = true;
    };

    // Load settings from local storage
    const loadSettings = async () => {
      try {
        const { value } = await Preferences.get({ key: "appSettings" });
        if (value) {
          const settings = JSON.parse(value);
          serverAddress.value = settings.serverAddress || "";
          locationInterval.value = settings.locationInterval ?? 30;
          serverToken.value = settings.serverToken || "";
          sendLocation.value = settings.sendLocation || false;
        }
      } catch (error) {
        console.error("Failed to load settings:", error);
        showSnackbarMessage("Failed to load settings", "error");
      }
    };

    // Navigate to the TrackingPage
    const goToTrackingPage = () => {
      if (!areSettingsComplete.value) {
        showSnackbarMessage("Please complete and save settings first", "error");
        return;
      }
      router.push("/TrackingPage");
    };

    // Load settings on app initialization
    onMounted(() => {
      loadSettings();
    });

    return {
      serverAddress,
      locationInterval,
      serverToken,
      sendLocation,
      saving,
      showSnackbar,
      snackbarText,
      snackbarColor,
      areSettingsComplete,
      saveSettings,
      goToTrackingPage,
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
