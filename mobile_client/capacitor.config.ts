import type { CapacitorConfig } from "@capacitor/cli";

const config: CapacitorConfig = {
  appId: "com.github.tiny_giraffes.life_beacon_360",
  appName: "Life Beacon 360",
  webDir: "dist",
  plugins: {
    BackgroundGeolocation: {
      // Android-specific configuration
      notification: {
        title: "Location Tracking Active",
        text: "Life Beacon 360 is tracking your location in the background",
        channelName: "Location Tracking",
        smallIcon: "ic_launcher_foreground",
      },
      // General configuration
      distanceFilter: 10, // Distance in meters
      stoppedElapsedTimeInSeconds: 0, // Start tracking immediately
      requestPermissions: true, // Request permissions on startup
    },
  },
  // Additional configuration for Android
  android: {
    backgroundColor: "#FFFFFF",
    contentInset: "automatic",
    allowMixedContent: true, // Required for some API requests to work in WebView
  },
  // Configuration for all platforms
  server: {
    cleartext: true, // Allow HTTP (non-HTTPS) connections for development/testing
  },
};

export default config;
