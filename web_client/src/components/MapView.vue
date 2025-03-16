<template>
  <v-container fluid>
    <v-card class="mb-4">
      <v-card-title class="text-h4">Location Map</v-card-title>
      <v-card-text>
        <div id="map" style="height: 600px; width: 100%"></div>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          prepend-icon="mdi-map-marker"
          :loading="loading"
          @click="fetchLocations"
        >
          Show Locations
        </v-btn>
      </v-card-actions>
    </v-card>
    <v-snackbar v-model="snackbar" :color="snackbarColor">
      {{ snackbarText }}
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import axios from "axios";
import "leaflet/dist/leaflet.css";
import * as L from "leaflet";

// State
const map = ref<L.Map | null>(null);
const markers = ref<L.Marker[]>([]);
const loading = ref(false);
const snackbar = ref(false);
const snackbarText = ref("");
const snackbarColor = ref("success");

// Server URL and token - these should come from environment variables in a real app
const API_URL = "http://localhost:8080/api/locations";
const API_TOKEN = import.meta.env.VITE_API_TOKEN || "";

// Calculate distance between two points in meters using the Haversine formula
function calculateDistance(
  lat1: number,
  lon1: number,
  lat2: number,
  lon2: number
): number {
  const R = 6371e3; // Earth radius in meters
  const φ1 = (lat1 * Math.PI) / 180;
  const φ2 = (lat2 * Math.PI) / 180;
  const Δφ = ((lat2 - lat1) * Math.PI) / 180;
  const Δλ = ((lon2 - lon1) * Math.PI) / 180;

  const a =
    Math.sin(Δφ / 2) * Math.sin(Δφ / 2) +
    Math.cos(φ1) * Math.cos(φ2) * Math.sin(Δλ / 2) * Math.sin(Δλ / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
}

// Fetch locations from the API
async function fetchLocations() {
  loading.value = true;

  try {
    const response = await axios.get(API_URL, {
      headers: {
        Authorization: `Bearer ${API_TOKEN}`,
      },
    });

    if (
      response.data &&
      Array.isArray(response.data) &&
      response.data.length > 0
    ) {
      displayLocationsOnMap(response.data);
      snackbarColor.value = "success";
      snackbarText.value = `Loaded ${response.data.length} location(s)`;
    } else {
      snackbarColor.value = "warning";
      snackbarText.value = "No locations found";
    }
  } catch (error) {
    console.error("Error fetching locations:", error);
    snackbarColor.value = "error";
    snackbarText.value = "Failed to fetch locations";
  } finally {
    loading.value = false;
    snackbar.value = true;
  }
}

// Display locations on the map
function displayLocationsOnMap(locations: any[]) {
  // Clear existing markers
  markers.value.forEach((marker) => {
    marker.remove();
  });
  markers.value = [];

  // Guard clause
  if (!map.value || locations.length === 0) return;

  // Create a strong reference to the map that TypeScript can understand
  const mapInstance = map.value as L.Map;

  // Get the latest location (first in the array since they're ordered by created_at DESC)
  const latestLocation = locations[0];

  // Check if all locations are within 10 meters of the latest location
  const allNearby = locations.every((loc) => {
    if (loc.id === latestLocation.id) return true;
    return (
      calculateDistance(
        latestLocation.latitude,
        latestLocation.longitude,
        loc.latitude,
        loc.longitude
      ) <= 10
    );
  });

  if (allNearby) {
    // If all locations are within 10 meters, only show the latest
    const latLng: L.LatLngExpression = [
      latestLocation.latitude,
      latestLocation.longitude,
    ];
    const marker = L.marker(latLng, {
      icon: createMarkerIcon("red"),
    });

    marker.addTo(mapInstance);
    marker.bindPopup(
      `Latest Location (${new Date(latestLocation.createdAt).toLocaleString()})`
    );
    markers.value.push(marker);

    // Set view to the latest location
    mapInstance.setView(latLng, 16);
  } else {
    // Add markers for all locations
    const bounds: L.LatLngBoundsExpression = [];

    locations.forEach((loc, index) => {
      // Red for latest location, blue for others
      const isLatest = index === 0;
      const latLng: L.LatLngExpression = [loc.latitude, loc.longitude];
      bounds.push(latLng);

      const marker = L.marker(latLng, {
        icon: createMarkerIcon(isLatest ? "red" : "blue"),
      });

      marker.addTo(mapInstance);
      marker.bindPopup(
        `${isLatest ? "Latest Location" : "Location"} (${new Date(
          loc.createdAt
        ).toLocaleString()})`
      );
      markers.value.push(marker);
    });

    // Fit the map to show all markers
    if (bounds.length > 0) {
      mapInstance.fitBounds(bounds as L.LatLngBoundsLiteral, {
        padding: [50, 50],
      });
    }
  }
}

// Create custom marker icon
function createMarkerIcon(color: "red" | "blue"): L.Icon {
  const iconUrl =
    color === "red"
      ? "https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png"
      : "https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-blue.png";

  return L.icon({
    iconUrl,
    shadowUrl:
      "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41],
  });
}

// Initialize map when component mounts
onMounted(() => {
  // Create map with a div element reference
  const mapElement = document.getElementById("map");
  if (!mapElement) return;

  // Initialize the map
  const mapInstance = L.map(mapElement).setView([0, 0], 2);
  map.value = mapInstance;

  // Add the OpenStreetMap tiles
  L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution: "© OpenStreetMap contributors",
  }).addTo(mapInstance);
});

// Clean up when component unmounts
onUnmounted(() => {
  if (map.value) {
    map.value.remove();
    map.value = null;
  }
});
</script>
