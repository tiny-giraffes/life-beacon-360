# Life Beacon 360

A location tracking application with a mobile client and backend server.

## Project Overview

Life Beacon 360 is a prototype location tracking system that consists of:

1. **Mobile Client**: Android app that tracks user location and sends it to a server
2. **Server**: Go backend that receives and stores location data

This project allows for real-time or interval-based location tracking with authorization.

## Mobile Client

### Setup Instructions

1. Install dependencies:

   ```bash
   cd mobile_client
   npm install
   ```

2. Build the app:

   ```bash
   npm run build
   ```

3. Add Android platform:

   ```bash
   npx cap add android
   ```

4. Copy the web assets to Android:

   ```bash
   npx cap sync android
   ```

5. Open in Android Studio:

   ```bash
   npx cap open android
   ```

6. Build and install the app on your device from Android Studio

### Usage Instructions

1. **Configure Settings**:

   - Enter the Server Address (e.g., `http://your-server:8080/api/locations`)
   - Enter the Server Token (same as configured in your server's `.env` file)
   - Set the Location Interval (in seconds, 0 = immediate)
   - Enable "Send Location to Server" if you want to send data
   - Press "Save Settings"

2. **Start Tracking**:

   - Navigate to the Tracking page
   - Press "Start Tracking" to begin location tracking
   - The app will request location permissions
   - Position data will appear on screen
   - If sending is enabled, location will be sent at the specified interval
   - Press "Stop Tracking" to stop location updates

3. **Background Operation**:
   - The app will continue tracking location even when minimized or the screen is off
   - A notification will appear showing tracking is active

## Server

### Setup Instructions

1. Create a `.env` file in the `server/config` directory:

   ```bash
   cp server/config/.env.example server/config/.env
   ```

2. Edit the `.env` file with your database credentials and API token:

   ```
   POSTGRES_USER=your_db_user
   POSTGRES_PASSWORD=your_db_password
   POSTGRES_DB=life_beacon_db
   POSTGRES_HOST=localhost
   POSTGRES_PORT=6000

   API_TOKEN=your_api_token
   ```

3. Start the PostgreSQL database:

   ```bash
   cp docker-compose-example.yml docker-compose.yml
   # Edit docker-compose.yml to match your credentials
   docker-compose up -d
   ```

4. Build and run the server:
   ```bash
   cd server
   go build -o life-beacon-server ./cmd
   ./life-beacon-server
   ```

### API Endpoints

#### Save Location

- **URL**: `/api/locations`
- **Method**: `POST`
- **Auth Required**: Yes (Token in Authorization header)
- **Headers**:
  - `Content-Type: application/json`
  - `Authorization: Bearer YOUR_API_TOKEN`
- **Body**:
  ```json
  {
    "latitude": 37.7749,
    "longitude": -122.4194
  }
  ```
- **Success Response**:
  - Code: 201
  - Content: `{ "message": "Location saved successfully" }`

## Troubleshooting

### Common Issues

1. **Location Not Updating**:

   - Ensure location permissions are granted
   - Check that GPS is enabled on your device
   - On Android 10+, make sure background location permission is granted

2. **Failed to Send Location**:

   - Verify server address is correct and includes the full path
   - Confirm the API token matches between client and server
   - Check that the server is running and accessible
   - Ensure your device has internet connectivity

3. **Database Connection Issues**:
   - Verify PostgreSQL is running (`docker ps`)
   - Check the database credentials in `.env`
   - Ensure the port is accessible

### Debugging

- Mobile client logs can be viewed in Android Studio's Logcat
- Server logs appear in the terminal where the server is running

## Security Notes

This is a prototype and has basic security features:

- Token-based API authentication
- Data sent over HTTPS is recommended for production

For a production environment, consider adding:

- User authentication
- Data encryption
- More robust error handling
- Rate limiting
- TLS/HTTPS
