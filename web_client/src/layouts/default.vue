<template>
  <v-app>
    <v-app-bar color="primary">
      <v-app-bar-title>Life Beacon 360</v-app-bar-title>

      <v-spacer></v-spacer>

      <v-btn to="/" variant="text">Home</v-btn>
      <v-btn to="/map" variant="text">Map</v-btn>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>

    <v-footer app color="primary" class="text-center">
      <span class="white--text"
        >&copy; {{ new Date().getFullYear() }} Life Beacon 360 - AGPL-3.0
        Licensed</span
      >
    </v-footer>
  </v-app>
</template>

<script setup lang="ts">
// No external component imports
</script>
