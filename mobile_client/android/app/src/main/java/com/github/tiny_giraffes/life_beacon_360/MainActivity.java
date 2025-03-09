package com.github.tiny_giraffes.life_beacon_360;

import android.content.Intent;
import android.os.Bundle;
import com.getcapacitor.BridgeActivity;
import com.getcapacitor.Plugin;
import com.getcapacitor.PluginCall;
import com.getcapacitor.PluginMethod;
import com.getcapacitor.annotation.CapacitorPlugin;
import com.getcapacitor.JSObject;

public class MainActivity extends BridgeActivity {
    @Override
    public void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        
        // Register our plugin
        registerPlugin(LocationServicePlugin.class);
    }
    
    @CapacitorPlugin(name = "LocationService")
    public static class LocationServicePlugin extends Plugin {
        @PluginMethod
        public void startService(PluginCall call) {
            Intent serviceIntent = new Intent(getContext(), LocationService.class);
            if (android.os.Build.VERSION.SDK_INT >= android.os.Build.VERSION_CODES.O) {
                getContext().startForegroundService(serviceIntent);
            } else {
                getContext().startService(serviceIntent);
            }
            call.resolve();
        }
        
        @PluginMethod
        public void stopService(PluginCall call) {
            Intent serviceIntent = new Intent(getContext(), LocationService.class);
            getContext().stopService(serviceIntent);
            call.resolve();
        }
    }
}