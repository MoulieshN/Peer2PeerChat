# Overview

The Peer-to-Peer Chat Application is a lightweight real-time chat system built with WebSocket technology. The project enables multiple users to connect to a chat server, broadcast messages, and view the list of online users.

# Key Features

Real-Time Communication: Users can send and receive messages instantly.

User Presence Tracking: Displays the list of currently connected users.

Broadcast Messaging: Messages are sent to all connected users.

Graceful Disconnection: Users leaving the chat are removed from the online users list.

---

# Project Structure

## Backend

Language: Go

Frameworks/Libraries:

Gorilla WebSocket: For WebSocket connections.

CloudyKit Jet: For template rendering.

BMizerany Pat: For routing.

### Key Files

handler.go: Implements WebSocket connection handling, message broadcasting, and user management.

routes.go: Defines HTTP and WebSocket routes.

main.go: Starts the server and listens for WebSocket channel events.

## Frontend

HTML Template: home.jet

CSS Framework: Bootstrap for styling.

JavaScript: WebSocket client handling user interactions and real-time updates.

---