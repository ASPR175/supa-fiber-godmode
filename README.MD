# Supabase + Fiber Auth Backend

A JWT-auth backend using Go Fiber, Supabase RLS, and Redis for token blacklisting.

## Features
- Supabase Row Level Security (RLS)
- Redis token blacklist for logout
- Secure env handling

## Setup
1. Copy `.env.example` to `.env`
2. Fill in your keys
3. `go run main.go`