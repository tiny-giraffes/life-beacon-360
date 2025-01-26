<template>
  <v-container class="py-4">
    <v-row>
      <v-col cols="12">
        <v-text-field
          v-model="serverAddress"
          label="Server Address"
          hint="Enter the server URL"
          outlined
          dense
        />
      </v-col>
      <v-col cols="12">
        <v-text-field
          v-model.number="locationInterval"
          label="Location Interval (seconds)"
          hint="Interval for sending location"
          type="number"
          outlined
          dense
        />
      </v-col>
      <v-col cols="12">
        <v-text-field
          v-model="serverToken"
          label="Server Token"
          hint="Token for server authentication"
          outlined
          dense
        />
      </v-col>
      <v-col cols="12">
        <v-checkbox v-model="sendLocation" label="Send Location to Server" />
      </v-col>
      <v-col cols="12" class="text-center">
        <v-btn color="primary" @click="saveSettings">Save Settings</v-btn>
      </v-col>
      <v-col cols="12" class="text-center mt-4">
        <v-btn color="secondary" @click="goToTrackingPage"
          >Go to Tracking Page</v-btn
        >
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { Preferences } from "@capacitor/preferences";

export default defineComponent({
  name: "App",
  setup() {
    // Router for navigation
    const router = useRouter();

    // Reactive states for settings
    const serverAddress = ref<string>("");
    const locationInterval = ref<number>(0);
    const serverToken = ref<string>("");
    const sendLocation = ref<boolean>(false);

    // Save settings to local storage
    const saveSettings = async () => {
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
      alert("Settings saved successfully!");
    };

    // Load settings from local storage
    const loadSettings = async () => {
      const { value } = await Preferences.get({ key: "appSettings" });
      if (value) {
        const settings = JSON.parse(value);
        serverAddress.value = settings.serverAddress || "";
        locationInterval.value = settings.locationInterval || 0;
        serverToken.value = settings.serverToken || "";
        sendLocation.value = settings.sendLocation || false;
      }
    };

    // Navigate to the TrackingPage
    const goToTrackingPage = () => {
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
