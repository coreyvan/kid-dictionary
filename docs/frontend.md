# Kid Dictionary: Frontend Architecture

## Tech Stack

- **Framework**: Vue 3 (Composition API)
- **Build**: Vite
- **Styling**: Tailwind CSS
- **State**: Pinia
- **Routing**: Vue Router
- **HTTP**: Fetch API or axios
- **PWA**: vite-plugin-pwa

## Project Structure

```
src/
├── components/
│   ├── chat/
│   │   ├── ChatMessage.vue
│   │   ├── ChatInput.vue
│   │   └── AgeBracketSelector.vue
│   └── common/
├── views/
│   ├── HomeView.vue
│   ├── LoginView.vue
│   ├── RegisterView.vue
│   ├── ConversationsView.vue
│   └── ChatView.vue
├── stores/
│   ├── auth.js
│   └── conversations.js
├── composables/
│   └── useApi.js
├── router/
│   └── index.js
└── App.vue
```

## Key Screens

1. **Login/Register** - Simple email/password forms
2. **Conversations List** - Past chats, sorted by recent
3. **Chat View** - Message history + input, age bracket selector
4. **Settings** - Default age bracket, account management

## Mobile-First Considerations

- Touch-friendly tap targets (min 44px)
- Bottom-positioned input (thumb-reachable)
- Swipe gestures for navigation (optional)
- Offline indicator + cached conversations
- Pull-to-refresh on conversation list

## PWA Features

- Installable on home screen
- Offline fallback page
- Cache conversation list for offline viewing
- Background sync for failed messages (stretch goal)

## Design Tokens (for UX collaboration)

```css
/* Colors - to be defined with UX */
--color-primary: TBD;
--color-bracket-little: TBD;    /* Warm, playful */
--color-bracket-growing: TBD;   /* Curious, bright */
--color-bracket-preteen: TBD;   /* Mature, confident */

/* Spacing */
--spacing-xs: 4px;
--spacing-sm: 8px;
--spacing-md: 16px;
--spacing-lg: 24px;
```

## TODO

- [ ] Create Figma/wireframes with UX partner
- [ ] Decide on component library (Headless UI? Build custom?)
- [ ] Plan loading states and skeleton screens
- [ ] Design empty states (no conversations yet)
- [ ] Error handling UX (network failures, rate limits)