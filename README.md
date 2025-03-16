# Life Beacon 360 (an early prototype)

A location tracking application with mobile and web clients and a backend server.

## Project Overview

Life Beacon 360 is a prototype location tracking system that consists of:

1. **Mobile Client**: Android app that tracks user location and sends it to a server
2. **Web Client**: Vue.js frontend for managing and viewing location data
3. **Server**: Go backend that receives and stores location data

This project allows for real-time or interval-based location tracking with authorization.

## Web Client

### Setup Instructions

1. Navigate to the web client directory:

   ```bash
   cd web_client
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Run the development server:

   ```bash
   npm run dev
   ```

4. Build for production:

   ```bash
   npm run build
   ```

5. Preview production build:
   ```bash
   npm run preview
   ```

### Project Structure

- Uses Vue 3 with Vuetify for UI components
- Typescript for type safety
- Vite as the build tool
- Pinia for state management
- Vue Router for navigation

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
- Web client can be debugged using browser developer tools
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

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0) - see the [LICENSE](LICENSE) file for details.

The AGPL-3.0 requires that if you modify this software and provide it as a service over a network (e.g., a web application), you must make your modified source code available to the users of that service.
