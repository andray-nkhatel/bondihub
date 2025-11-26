# Vue Frontend Template

A Vue 3 frontend template with PrimeVue, authentication, and role-based access control.

## Features

- **Vue 3** with Composition API
- **PrimeVue 4** UI components
- **Tailwind CSS** for styling
- **Vuex** for state management
- **Vue Router** with navigation guards
- **JWT Authentication** with token refresh
- **Role-based Access Control (RBAC)**
- **Responsive Layout** with sidebar navigation

## Getting Started

### Prerequisites

- Node.js 18+
- npm or yarn

### Installation

```bash
npm install
```

### Configuration

Create a `.env` file in the root directory:

```env
VITE_API_BASE_URL=http://localhost:5000/api
```

### Development

```bash
npm run dev
```

### Build

```bash
npm run build
```

## Project Structure

```
src/
├── assets/          # Styles and static assets
├── components/      # Reusable Vue components
├── layout/          # Layout components (AppLayout, AppMenu, etc.)
├── router/          # Vue Router configuration
│   └── guard/       # Navigation guards
├── service/         # API services
├── store/           # Vuex store
│   └── modules/     # Store modules (auth)
└── views/           # Page components
    ├── dashboard/   # Dashboard views
    ├── pages/       # Auth pages (login, register)
    └── users/       # User management views
```

## Authentication

The template includes a complete authentication system:

- Login/Register pages
- JWT token management with auto-refresh
- Protected routes with navigation guards
- Role-based menu filtering

### Adding Protected Routes

```javascript
{
    path: '/admin',
    component: AdminPage,
    meta: {
        requiresAuth: true,
        roles: ['Admin']
    }
}
```

### Checking Roles in Components

```javascript
import { useStore } from 'vuex';

const store = useStore();
const isAdmin = store.getters['auth/hasRole']('Admin');
const hasAnyRole = store.getters['auth/hasAnyRole'](['Admin', 'User']);
```

## API Configuration

Update `src/service/api.service.js` to add your API endpoints:

```javascript
export const myService = {
    async getAll() {
        const response = await apiClient.get('/my-endpoint');
        return response.data;
    },
    // Add more methods...
};
```

## License

MIT
