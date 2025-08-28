## Agent OS Documentation

### Product Context
- **Mission & Vision:** @.agent-os/product/mission.md
- **Technical Architecture:** @.agent-os/product/tech-stack.md
- **Development Roadmap:** @.agent-os/product/roadmap.md
- **Decision History:** @.agent-os/product/decisions.md

### Development Standards
- **Code Style:** @~/.agent-os/standards/code-style.md
- **Best Practices:** @~/.agent-os/standards/best-practices.md

### Project Management
- **Active Specs:** @.agent-os/specs/
- **Spec Planning:** Use `@~/.agent-os/instructions/create-spec.md`
- **Tasks Execution:** Use `@~/.agent-os/instructions/execute-tasks.md`

## Workflow Instructions

When asked to work on this codebase:

1. **First**, check @.agent-os/product/roadmap.md for current priorities
2. **Then**, follow the appropriate instruction file:
   - For new features: @.agent-os/instructions/create-spec.md
   - For tasks execution: @.agent-os/instructions/execute-tasks.md
3. **Always**, adhere to the standards in the files listed above

## Important Notes

- Product-specific files in `.agent-os/product/` override any global standards
- User's specific instructions override (or amend) instructions found in `.agent-os/specs/...`
- Always adhere to established patterns, code style, and best practices documented above.

## Project-Specific Guidelines

### Ferry Boarding Management System

This is a comprehensive boarding management system for ferry operations featuring:

- **Core Systems:** Ticket purchasing, scheduling, seat management
- **Operations:** Point of sales, refund processing, rescheduling
- **Management:** Operator management, port management
- **Customer Service:** Helpdesk system and AI-powered chatbot

### Technology Stack

- **Backend:** Go (Golang) with Gin/Fiber framework
- **Frontend:** React 18+ with Vite
- **Database:** PostgreSQL 17+
- **Authentication:** JWT with Redis session management
- **Real-time:** WebSockets for live updates

### Development Priorities

Follow the roadmap phases:
1. Core Foundation - Database, API, and basic management
2. Ticketing & Scheduling - Complete booking flow
3. Point of Sales & Operations - Terminal and onboard sales
4. Customer Service & AI - Chatbot and helpdesk
5. Advanced Features & Scale - Enterprise capabilities

### Code Organization

- `/backend` - Go API and business logic
- `/frontend` - React application
- `/database` - PostgreSQL schemas and migrations
- `/docs` - API documentation and guides
- `/.agent-os` - Agent OS configuration and specs