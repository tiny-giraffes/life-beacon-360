# LB360 MVP Functional Requirements

## Project Overview

Life Beacon 360 (LB360) is a location tracking application designed for families and companies to track the location of group members with geofencing capabilities and alert system.

## Tech Stack

### Backend
- **Framework**: Go with Echo v4
- **Database**: PostgreSQL with GORM
- **Authentication**: Server-side sessions with simple tokens

### Web Client
- **Framework**: Vue 3 + Vuetify 3
- **Build Tool**: Vite
- **State Management**: Pinia

### Mobile Client
- **Frontend**: Vue 3 + Vuetify 3 (Capacitor)
- **Background Service**: Native Java service for location tracking

## 1. User Management & Registration

### Admin Bootstrap
- **Setup Wizard**: First-time app setup creates initial admin when database is empty
- One-time setup process accessible via **both web and mobile interfaces**
- Requires: email, name, full name, phone, password, **group**
- **System automatically creates initial group** "Main" during setup
- **Initial admin is assigned to the "Main" group**

### Group Management Rules
- **Initial group**: System creates "Main" group during first setup
- **Admin permissions on groups**:
  - Create new groups
  - Edit group names/details
  - Delete groups (only if no users assigned)
  - **Cannot delete the "Main" group** (system protection)
- **Group structure**: Flat, no hierarchies
- **User assignment**: Each user belongs to exactly one group

### User Registration Flow
1. **Admin creates user account** with:
   - Required: email, name, full name, phone, **group selection**
   - Optional: photo, job title
2. **System sends invitation email** to user with temporary credentials
3. **User logs in with temporary password** 
4. **System suggests password change** on first login (optional, not mandatory)
5. **User account activated** after user logs in

### Role System
- **Two roles only**: `admin` and `user`
- **Admin protection**: **Cannot delete or demote the last admin** (system must have at least one admin)
- **Admin permissions**: 
  - Create/edit/delete users
  - Create other admins
  - **Manage groups** (create/edit/delete)
  - Create geofences and alerts
  - **View location data of ANY user/admin** (full visibility)
  - **Configure location visibility permissions** for users
- **User permissions**:
  - **View own location data** (always)
  - **View other users/groups location data** (only if permitted by admin)
  - Update own profile (except role/group)
  - Receive alerts

### Location Visibility Rules
- **Admin**: Can see all users' locations (no restrictions)
- **User**: Can see own location + other users/groups as configured by admin
- **Permission granularity**: Admin can grant user access to specific users or entire groups

### User Profile Fields
- **Required**: email, name, full name, phone, password
- **Optional**: photo (stored as binary data in database), job title
- **System fields**: role, group assignment, active status, first login tracking

## Database Models Specification

### Group Model
- **ID**: UUID (primary key)
- **Name**: String (required, max 255 chars)
- **IsDefault**: Boolean (default: false, protects "Main" group from deletion)
- **CreatedAt**: Timestamp (auto-generated)
- **UpdatedAt**: Timestamp (auto-updated)
- **DeletedAt**: Timestamp (nullable, soft delete)

### User Model
- **ID**: UUID (primary key)
- **GroupID**: UUID (foreign key, required)
- **Username**: String (unique, required, max 100 chars)
- **Email**: String (required, max 255 chars)
- **PasswordHash**: String (required, max 255 chars)
- **FullName**: String (required, max 255 chars)
- **Phone**: String (required, max 50 chars)
- **Photo**: Binary/BLOB (optional, stores image data in database)
- **JobTitle**: String (optional, max 100 chars)
- **Role**: String (required, default: 'user', values: 'admin'|'user', max 50 chars)
- **IsActive**: Boolean (default: true)
- **InviteToken**: String (optional, max 255 chars, for invitation flow)
- **FirstLogin**: Boolean (default: true, for password change suggestion)
- **SessionMaxDuration**: Integer (days, per-user session absolute maximum, uses system default if null)
- **SessionActivityExtension**: Integer (days, per-user session activity extension, uses system default if null)
- **CreatedAt**: Timestamp (auto-generated)
- **UpdatedAt**: Timestamp (auto-updated)
- **DeletedAt**: Timestamp (nullable, soft delete)

### Session Model
- **Token**: String (primary key, session identifier)
- **UserID**: UUID (foreign key, required)
- **CreatedAt**: Timestamp (session creation time)
- **LastActivity**: Timestamp (last request time)
- **IPAddress**: String (client IP address)
- **UserAgent**: String (device/browser info)
- **ClientType**: String (values: 'web'|'mobile')
- **IsActive**: Boolean (default: true)

## Next Requirements to Define

1. **Geofencing System** (zones, triggers, alerts)
2. **Alert & Notification System** (push notifications, email, in-app)
3. **Location Tracking Rules** (privacy settings, tracking intervals)
4. **Dashboard & Monitoring** (maps, status indicators, activity feeds)

## 2. Authentication & Sessions

### Session Management Approach
- **Simple session tokens**: Random UUID/string tokens (not JWTs)
- **Server-side sessions**: All session data stored on server
- **Unified session store**: Same session store for web and mobile clients
- **Instant revocation**: Admin can terminate any session immediately
- **Multiple concurrent sessions**: Users can be logged in from multiple devices simultaneously
- **Hybrid session lifetime**: Per-user configurable absolute maximum + activity-based extension

### Client-Specific Implementation
- **Web client**: Session cookies (httpOnly, secure, SameSite)
- **Mobile client**: Session token in Authorization header (`Authorization: Session <token>`)
- **Token storage**: Mobile stores token in secure storage (Keychain/Keystore)

### Session Features
- **Admin session control**:
  - View all active sessions per user
  - Revoke individual sessions
  - Revoke all sessions for a user
  - Force logout on role/permission changes
  - **Configure session lifetime per user** (override system defaults or use defaults)
- **Session tracking**:
  - Login time, last activity timestamp
  - Device/browser information
  - IP address
  - Client type (web/mobile)
- **Session lifetime**: Hybrid approach per user
  - **Absolute maximum**: Hard limit from creation (e.g., 90 days max)
  - **Activity-based extension**: Extends with each request (e.g., 14 days from last activity)
  - **Configuration**: System default values with admin override per user
- **Session persistence**: Survives app restarts on mobile

### Authentication Flow
1. **Login**: User submits credentials
2. **Server validation**: Verify credentials against database
3. **Session creation**: Generate random token, store session data
4. **Token delivery**: Return token to client (cookie for web, response body for mobile)
5. **Request validation**: Every API request validates token against session store
6. **Logout**: Remove session from server store

### Session Data Model
- **Complete Session Model specification included in Database Models section above**

## Next Requirements to Define

1. **Geofencing System** (zones, triggers, alerts)
2. **Alert & Notification System** (push notifications, email, in-app)
3. **Location Tracking Rules** (privacy settings, tracking intervals)
4. **Dashboard & Monitoring** (maps, status indicators, activity feeds)

## Status

✅ **Completed**: User Management & Registration functional requirements
✅ **Completed**: Authentication & Sessions functional requirements
🚧 **In Progress**: Geofencing System (next conversation)
⏳ **Pending**: Alerts, Location Tracking, Dashboard