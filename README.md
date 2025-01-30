# Task Management System

## Overview

This project is a **Task Management System** that allows users to manage tasks, leave comments, and track task history. The system includes authentication. 

The UI is designed to be similar to **Jira**, with a clean, user-friendly interface for managing tasks and their associated comments and histories.

---

## Features

- **Task Management**: Create and manage tasks with different statuses (TODO, IN_PROGRESS, DONE, ARCHIVED).
- **Commenting**: Users can leave comments on tasks. Only the creator of a comment can modify or delete it.
- **History Tracking**: Tracks changes made to tasks, such as updates to the title and status.
- **User Roles**: Authentication and authorization using user roles (admin, user) (planned but not implemented).
  
---

## Tech Stack

- **Backend**: Go with the [Fiber](https://github.com/gofiber/fiber) framework
- **Database**: Supabase (PostgreSQL)
- **Authentication**: JWT tokens for user authentication
- **Docker**: For containerization of both the API service and the database

---

## Database Schema

### Tables (for Relational DB - Supabase)

![robinhood-assignment](https://github.com/user-attachments/assets/e91ca353-1b88-4d12-a267-49e883a0b848)

1. **Tasks**
   - `id` (varchar, Primary Key)
   - `title` (text)
   - `status` (Enum: TODO, IN_PROGRESS, DONE, ARCHIVED)
   - `createdAt` (timestamp)
   - `updatedAt` (timestamp)
   - `createdBy` (varchar, Foreign Key to Users)
   - `updatedBy` (varchar, Foreign Key to Users)

2. **Users**
   - `id` (varchar, Primary Key)
   - `name` (varchar) (not implemented)
   - `email` (varchar, Unique)
   - `role` (enum: admin, user) (not implemented)

3. **Comments**
   - `id` (varchar, Primary Key)
   - `task_id` (varchar, Foreign Key to Tasks)
   - `content` (text)
   - `createdAt` (timestamp)
   - `updatedAt` (timestamp)
   - `createdBy` (varchar, Foreign Key to Users)

4. **History**
   - `id` (varchar, Primary Key)
   - `task_id` (varchar, Foreign Key to Tasks)
   - `updatedBy` (varchar, Foreign Key to Users)
   - `changes` (json)
   - `updatedAt` (timestamp)

---

## API Endpoints

See: https://www.notion.so/API-Docs-188a4ff4b2a7803e8e9bebff2e8dde5d?pvs=4
