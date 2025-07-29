# LB360 MVP Functional Requirements

## Project Overview

Life Beacon 360 (LB360) is a location tracking application designed for families and companies to track the location of group members with geofencing capabilities and alert system.

## Tech Stack

### Backend
- **Framework**: Go with Echo v4
- **Database**: PostgreSQL with GORM
- **Cache**: Redis for session and permission caching
- **Authentication**: Server-side sessions with simple tokens
- **Notifications**: Push notifications (FCM/APNS) for critical updates

### Web Client
- **Framework**: Vue 3 + Vuetify 3
- **Build Tool**: Vite
- **State Management**: Pinia

### Mobile Client
- **Frontend**: Vue 3 + Vuetify 3 (Capacitor)
- **Background Service**: Native Java service for location tracking

### Desktop Client
- **Technology**: TBD (Electron or similar)
- **Features**: System service for background tracking

## Core Architecture Principles

### Permission & Settings Philosophy
- **Default State**: All actions and settings are admin-only by default
- **Global Settings**: All settings are global by default
- **Delegation**: Admin can grant specific permissions to users or groups
- **Revocation**: Admin can revoke any permission or setting override
- **Precedence**: User-specific settings override group settings, which override global settings
- **Consistency**: ALL settings and permissions follow this same pattern - no exceptions

## 1. User Management & Registration

### Admin Bootstrap
- **Setup Wizard**: First-time app setup creates initial admin when database is empty
- One-time setup process accessible via **web, mobile, and desktop interfaces**
- Requires: email, username, full name, phone, password
- **System automatically creates initial group** "Main" during setup
- **Initial admin is assigned to the "Main" group**

### Group Management
- **Initial group**: System creates "Main" group during first setup
- **Group operations** (permission: `can_manage_groups`):
  - Create new groups
  - Edit group names/details
  - Delete groups (only if no users assigned)
  - **Cannot delete the "Main" group** (system protection)
- **Group structure**: Flat, no hierarchies
- **User assignment**: Each user belongs to exactly one group

### User Registration Flow
1. **User creation** (permission: `can_manage_users`):
   - Required: email, username, full name, phone, **group selection**
   - Optional: photo, job title
   - **Role selection**: Only admins can create other admins (users with permission can only create users)
2. **System sends invitation email** with temporary credentials
3. **User logs in with temporary password**
4. **System suggests password change** on first login (optional)
5. **User account activated**

### Role System
- **Two roles only**: `admin` and `user`
- **Admin protection**: 
  - Cannot delete or demote the last admin
  - **Admin permissions CANNOT be revoked** (admins always have full permissions)
- **Admin-to-Admin actions**:
  - One admin can demote another admin to user (except the last admin)
  - No granular permission management between admins
  - To limit an admin's access, they must be demoted to user first
- **Default permissions**:
  - Admin: All permissions always (cannot be modified)
  - User: No permissions by default
- **Permission delegation**: 
  - **Only admins** can grant/revoke permissions
  - Users can NEVER edit permissions
  - Permissions can only be granted/revoked for users, not admins

### User Profile Fields
- **Required**: email, username, full name, phone, password
- **Optional**: photo (stored as binary data), job title
- **System fields**: role, group assignment, active status, first login tracking

## 2. Authentication & Sessions

### Session Management Approach
- **Simple session tokens**: Random UUID/string tokens (not JWTs)
- **Server-side sessions**: All session data stored on server
- **Unified session store**: Same session store for web, mobile, and desktop clients
- **Instant revocation**: Admin can terminate any session immediately
- **Multiple concurrent sessions**: Users can be logged in from multiple devices

### Client-Specific Implementation
- **Web client**: Session cookies (httpOnly, secure, SameSite)
- **Mobile/Desktop client**: Session token in Authorization header
- **Token storage**: Secure storage (Keychain/Keystore)

### Session Configuration
- **Session lifetime**: Configured through settings hierarchy
  - Uses GlobalSettings/GroupSettings/UserSettings
  - Absolute maximum duration
  - Activity-based extension
- **Session tracking**:
  - Login time, last activity timestamp
  - Device/browser information
  - IP address
  - Client type (web/mobile/desktop)

### Authentication Flow
1. **Login**: User submits credentials
2. **Server validation**: Verify credentials against database
3. **Session creation**: Generate random token, store session data
4. **Token delivery**: Return token to client
5. **Request validation**: Every API request validates token
6. **Logout**: Remove session from server store

## 3. Location Tracking Rules

### Core Permission Philosophy
*Same as defined in Core Architecture Principles*

### Permission System
- **Permission Types**:
  - `can_control_own_tracking`: User can start/stop their own tracking
  - `can_view_location`: User can view location data (specify targets)
  - `can_control_tracking`: User can start/stop tracking (specify targets)
  - `can_export_data`: User can export location data
  - `can_view_reports`: User can view tracking reports
  - `can_manage_users`: User management operations
  - `can_manage_groups`: Group management operations
  - Additional permissions can be added as needed
  - **Critical exception**: Only admins can grant/revoke permissions or create other admins
- **Permission Targets**:
  - `self`: Permission applies to user's own data
  - `group`: Permission applies to user's group members
  - `specific_users`: Permission applies to listed users
  - `all`: Permission applies to all users
- **Permission Management**:
  - **Grant**: Only admins can grant permissions to users/groups
  - **Revoke**: Only admins can revoke permissions from users/groups
  - **Admin immunity**: Admin permissions cannot be modified - admins always have full access
  - **Reset**: Only admins can clear all permissions for a user/group
  - **Audit**: All permission changes are logged
  - **Users limitation**: Users can NEVER modify permissions, regardless of what permissions they have

### Settings Hierarchy
- **Global Settings** (default for all):
  - Tracking interval
  - Polling interval
  - Accuracy mode
  - Data retention period
  - Session lifetime settings
  - Battery thresholds
  - All other configurable parameters
- **Group Settings** (override global):
  - Any global setting can be overridden at group level
  - Admin can remove group overrides
- **User Settings** (highest precedence):
  - Any setting can be overridden at user level
  - Admin can remove user overrides

### Tracking Modes
- **Client Tracking**: Any client (mobile app, desktop app, or web client) can send location data
- **Multi-Client Support**: If user has multiple clients running, each sends location independently
- **Remote Tracking Control**: Users can start/stop tracking based on granted permissions
- **Mandatory Tracking**: Configured through the settings hierarchy (global/group/user)

### Tracking Control Permissions
- **Own Device Control**:
  - Users control when their clients send location data (if allowed by admin settings)
  - Each client operates independently - user can have tracking on mobile but off on desktop/web
  - If admin disables user control, all clients track based on server configuration
- **Remote Control**:
  - Users with permission can remotely start/stop tracking for specific users
  - Permissions are granted per target user (User A can control User B's tracking)
  - Admin configures who can control whom
  - Control actions are logged with timestamp and initiating user
  - **Remote control mechanism**: User changes tracking status on server; target's clients receive push notification to check configuration immediately, or discover change on next poll
- **Admin Control**:
  - Admin can start/stop tracking for any user
  - Admin can enable/disable user's ability to control their own tracking
  - Admin grants/revokes tracking control permissions between users

### Tracking State Synchronization
- **Client Polling**:
  - Clients check tracking configuration on startup
  - Periodic checks based on configured polling interval
  - User-specific polling interval takes precedence over system-wide
  - Checks when returning from background/resuming
  - **Push-triggered polling**: Push notifications trigger immediate configuration check
- **Configuration Source**:
  - Server maintains authoritative tracking state for each user
  - Clients discover updates via polling (immediate when pushed, periodic otherwise)
  - Polling Optimization: Check Redis cache before database

### Tracking Intervals & Accuracy
- **Configurable Intervals**:
  - **Standard intervals**: 60, 300, 600, 900, 1800, 3600 seconds (1 min to 1 hour)
  - **System-wide default**: Configurable by admin
  - **User-specific override**: Admin can set per-user interval
  - **Precedence**: User-specific settings always override system-wide defaults
- **Location Accuracy Settings**:
  - **High accuracy**: GPS + Network + WiFi (default)
  - **Balanced**: Network + WiFi (battery saving)
  - **Low power**: Network only (maximum battery saving)
  - **Web client**: Uses browser geolocation API (accuracy depends on device)
- **Movement Detection**:
  - **Distance filter**: 10 meters (only update if moved more than threshold)
  - **Stationary detection**: Reduce update frequency when not moving

### Client-Specific Implementation

#### Mobile Client (Android)
- **Native Java Service**: Dedicated background service handles all location tracking
- **Service Architecture**:
  - Java service runs independently of UI
  - Minimal JS-Java boundary crossing for battery efficiency
  - Service handles location acquisition and server transmission
  - UI only needed for configuration and status display
- **Persistent Notification**: Shows "Location Tracking Active" with app icon
- **Battery Optimization**: Exempt from battery optimization when tracking

#### Desktop Client
- **Background Process**: Runs as system service/daemon
- **Location Source**: WiFi positioning, IP geolocation, or connected GPS
- **System Tray**: Shows tracking status icon
- **Auto-start**: Option to start with system boot

#### Web Client
- **Browser Geolocation**: Uses HTML5 Geolocation API
- **Limitations**:
  - Requires HTTPS in production
  - User must keep browser tab open
  - May prompt for permission repeatedly
  - Limited background tracking capability
- **Status Indicator**: Shows tracking status in browser tab/window

### Data Transmission Rules
- **Security**:
  - All location data sent over HTTPS in production
  - Include session token with each transmission
  - Reject unauthorized location updates
- **Data Format**:
  - Latitude/Longitude (decimal degrees, 6 decimal precision)
  - Timestamp (UTC)
  - Accuracy radius (meters)
  - Speed (optional, m/s)
  - Bearing (optional, degrees)
  - Altitude (optional, meters)
  - **Client Type**: Identifies source (mobile/desktop/web)
  - **Client ID**: Unique identifier for multi-client tracking
- **Transmission Method**:
  - HTTP POST to `/api/locations` endpoint
  - Clients poll `/api/users/{id}/tracking-config` for configuration updates
  - Push notifications trigger immediate polling for critical changes

### Battery & Performance Management
- **Mobile Client**:
  - Native service optimized for battery efficiency
  - Show battery usage warning when high-frequency tracking
  - Automatic degradation to lower accuracy if battery < 20%
  - Stop tracking if battery < 5% (user configurable)
- **Desktop Client**:
  - Minimal CPU usage when stationary
  - Configurable resource limits
- **Web Client**:
  - Throttle requests based on browser limitations
  - Warn user about keeping tab active
- **Adaptive Tracking** (all clients):
  - Increase interval when stationary
  - Decrease interval when moving fast
  - Night mode: Reduce frequency during configured hours

### Error Handling & Recovery
- **Location Failures**:
  - Retry location acquisition 3 times
  - Fall back to last known location if GPS fails
  - Alert user if location services disabled
- **Network Failures**:
  - Queue locations locally when offline
  - Exponential backoff for failed transmissions
  - **Backoff upper boundary**: Configurable in global settings
  - **Retry attempts**: Indefinite (never give up)
  - Backoff interval never exceeds configured upper boundary
- **Permission Issues**:
  - Prompt user to re-enable permissions
  - Log permission denials for admin review
  - Provide settings shortcut for permission management

### Tracking Status & Indicators
- **All Clients Interface**:
  - Clear tracking status display
  - Shows if tracking is mandatory (can't be stopped)
  - Last update timestamp
  - Current accuracy level
  - Connection status to server
- **Management Interface** (web/desktop):
  - View tracking status for all permitted users
  - Start/Stop tracking buttons for users with control permission
  - Last known location and timestamp
  - Tracking history and logs
  - Multi-client status view (see all active clients per user)
- **Admin Dashboard**:
  - View/control tracking for all users
  - See all active clients per user
  - Tracking compliance report
  - View tracking control logs

*Full Location Tracking Rules specification continues as defined in the separate document*

## Database Models

### Core Models

#### Group Model
- **ID**: UUID (primary key)
- **Name**: String (required, max 255 chars)
- **IsDefault**: Boolean (default: false)
- **CreatedAt**: Timestamp
- **UpdatedAt**: Timestamp
- **DeletedAt**: Timestamp (soft delete)

#### User Model
- **ID**: UUID (primary key)
- **GroupID**: UUID (foreign key, required)
- **Username**: String (unique, required, max 100 chars)
- **Email**: String (required, max 255 chars)
- **PasswordHash**: String (required, max 255 chars)
- **FullName**: String (required, max 255 chars)
- **Phone**: String (required, max 50 chars)
- **Photo**: Binary/BLOB (optional)
- **JobTitle**: String (optional, max 100 chars)
- **Role**: String (required, default: 'user', values: 'admin'|'user')
- **IsActive**: Boolean (default: true)
- **InviteToken**: String (optional, max 255 chars)
- **FirstLogin**: Boolean (default: true)
- **CreatedAt**: Timestamp
- **UpdatedAt**: Timestamp
- **DeletedAt**: Timestamp (soft delete)

#### Session Model
- **Token**: String (primary key)
- **UserID**: UUID (foreign key, required)
- **CreatedAt**: Timestamp
- **LastActivity**: Timestamp
- **IPAddress**: String
- **UserAgent**: String
- **ClientType**: String ('web'|'mobile'|'desktop')
- **IsActive**: Boolean (default: true)

### Settings & Permissions Models

#### GlobalSettings Model
- **ID**: UUID (primary key)
- **TrackingInterval**: Integer (seconds, default: 300)
- **PollingInterval**: Integer (seconds, default: 60)
- **AccuracyMode**: String (default: 'high')
- **DataRetentionDays**: Integer (nullable - null = forever)
- **SessionMaxDuration**: Integer (days, default: 90)
- **SessionActivityExtension**: Integer (days, default: 14)
- **MandatoryTracking**: Boolean (default: false)
- *[All other configurable parameters]*
- **CreatedAt**: Timestamp
- **UpdatedAt**: Timestamp

#### GroupSettings Model
- **ID**: UUID (primary key)
- **GroupID**: UUID (foreign key, unique)
- **Settings**: JSON (overrides for any GlobalSettings fields)
- **CreatedAt**: Timestamp
- **UpdatedAt**: Timestamp

#### UserSettings Model
- **ID**: UUID (primary key)
- **UserID**: UUID (foreign key, unique)
- **Settings**: JSON (overrides for any GlobalSettings fields)
- **TrackingEnabled**: Boolean (current tracking state)
- **LastModifiedBy**: UUID
- **LastModifiedAt**: Timestamp
- **CreatedAt**: Timestamp
- **UpdatedAt**: Timestamp

#### Permission Model
- **ID**: UUID (primary key)
- **UserID**: UUID (foreign key, nullable)
- **GroupID**: UUID (foreign key, nullable)
- **PermissionType**: String
- **TargetType**: String ('self'|'group'|'specific_users'|'all')
- **TargetUsers**: JSON array of UUIDs
- **GrantedBy**: UUID
- **CreatedAt**: Timestamp
- **UpdatedAt**: Timestamp

### Location Models

#### Location Model
- **ID**: UUID (primary key)
- **UserID**: UUID (foreign key, required)
- **ClientID**: String
- **ClientType**: String ('mobile'|'desktop'|'web')
- **Latitude**: Decimal (precision 6)
- **Longitude**: Decimal (precision 6)
- **Accuracy**: Float (meters, optional)
- **Altitude**: Float (meters, optional)
- **Speed**: Float (m/s, optional)
- **Bearing**: Float (degrees, optional)
- **RecordedAt**: Timestamp
- **ReceivedAt**: Timestamp
- **BatteryLevel**: Integer (percent, optional)
- **IsStationary**: Boolean (default: false)
- **CreatedAt**: Timestamp

#### ActiveClient Model
- **ID**: UUID (primary key)
- **UserID**: UUID (foreign key)
- **ClientID**: String (unique per instance)
- **ClientType**: String
- **LastSeen**: Timestamp
- **TrackingActive**: Boolean
- **AppVersion**: String (optional)
- **DeviceInfo**: JSON (optional)
- **CreatedAt**: Timestamp
- **UpdatedAt**: Timestamp

### Audit Models

#### AuditLog Model
- **ID**: UUID (primary key)
- **UserID**: UUID (user who performed action)
- **Action**: String
- **EntityType**: String
- **EntityID**: UUID
- **Changes**: JSON (before/after values)
- **IPAddress**: String
- **UserAgent**: String
- **CreatedAt**: Timestamp

## Configuration Resolution

### Settings Resolution
When determining a setting value for a user:
1. Check UserSettings for override
2. If not found, check GroupSettings for user's group
3. If not found, use GlobalSettings value
4. Apply the first value found

### Permission Resolution
When checking if a user can perform an action:
1. If user is admin role, allow all actions (admins have immutable full permissions)
2. If user is trying to modify permissions, must be admin role
3. For regular users, check user-specific permissions
4. Check group permissions for user's group
5. If no permission found, deny action
6. Validate permission target matches requested action

## 4. Caching Strategy

### Redis Cache Implementation
- **Technology**: Redis
- **Purpose**: Reduce database load for frequently accessed data
- **Connection pooling**: Maintain persistent connections

### Cached Data Types

#### Session Cache
- **Key format**: `session:{token}`
- **Data**: User ID, permissions, last activity
- **TTL**: Match session timeout
- **Invalidation**: On logout or session termination

#### Permission Resolution Cache
- **Key format**: `perm:{userID}:{permissionType}:{targetType}`
- **Data**: Boolean result of permission check
- **TTL**: 5 minutes
- **Invalidation**: On any permission grant/revoke

#### Settings Resolution Cache
- **Key format**: `settings:{userID}:{settingName}`
- **Data**: Resolved setting value (after hierarchy resolution)
- **TTL**: 5 minutes
- **Invalidation**: On any setting change at any level

#### User Basic Info Cache
- **Key format**: `user:{userID}`
- **Data**: Username, role, group ID, active status
- **TTL**: 10 minutes
- **Invalidation**: On user update

### Cache Invalidation Rules
- **Cascading invalidation**: 
  - Global setting change → Clear all user settings cache
  - Group setting change → Clear all group member settings cache
  - Permission change → Clear affected user permission cache
- **Transaction safety**: Clear cache AFTER database commit
- **Fallback**: Always check database if cache miss

## 5. Push Notification System

### Push Notification Implementation
- **Providers**: 
  - Firebase Cloud Messaging (FCM) for Android
  - Apple Push Notification Service (APNS) for iOS
  - Web Push for browser clients (optional)
- **Purpose**: Trigger immediate client action without persistent connections
- **Reliability**: System works without push notifications via polling

### Notification Types

#### Configuration Change Notification
- **Trigger**: Any tracking configuration change
- **Payload**: `{"type": "config_change", "action": "poll"}`
- **Client Action**: Immediately poll `/api/users/me/tracking-config`
- **No sensitive data**: Notification doesn't contain actual changes

#### Tracking Control Notification
- **Trigger**: Tracking started/stopped by another user
- **Payload**: `{"type": "tracking_control", "action": "poll"}`
- **Client Action**: Poll configuration and update tracking state
- **User Feedback**: Show who initiated the change after polling

#### Emergency Alert Notification
- **Trigger**: Geofence violations, panic button, admin alerts
- **Payload**: `{"type": "emergency", "alertId": "uuid"}`
- **Client Action**: Poll `/api/alerts/{id}` for details
- **Priority**: High priority notification with sound

#### Session Termination Notification
- **Trigger**: Admin terminates session or permissions revoked
- **Payload**: `{"type": "session_terminated"}`
- **Client Action**: Clear local session, redirect to login
- **Immediate**: No polling needed, just logout

### Push Notification Architecture
- **Token Management**: 
  - Clients register device tokens on login
  - Update tokens when refreshed
  - Remove tokens on logout
- **Delivery**: 
  - Queue notifications in Redis
  - Retry failed deliveries with exponential backoff
  - Log delivery failures for monitoring
- **Fallback**: 
  - Polling ensures all changes discovered eventually
  - Default poll interval: 60 seconds (configurable)
  - Push notifications only optimize latency

### Privacy & Security
- **No sensitive data in push**: Only notification type and action
- **All data fetched via authenticated API**: After receiving push
- **Encryption**: Use platform-provided encryption
- **Token validation**: Verify device tokens belong to user

## Next Requirements to Define

1. **Geofencing System** (zones, triggers, alerts)
2. **Alert & Notification System** (push notifications, email, in-app)
3. **Dashboard & Monitoring** (maps, status indicators, activity feeds)
4. **Reporting System** (tracking reports, compliance reports)

## Status

✅ **Completed**: User Management & Registration functional requirements
✅ **Completed**: Authentication & Sessions functional requirements  
✅ **Completed**: Location Tracking Rules functional requirements
✅ **Completed**: Core permission and settings architecture
✅ **Completed**: Caching strategy with Redis
✅ **Completed**: Push notification system for time-sensitive updates
🚧 **Next**: Geofencing System
⏳ **Pending**: Alert System details, Dashboard & Monitoring, Reporting