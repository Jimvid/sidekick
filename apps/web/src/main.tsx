import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { ClerkProvider } from '@clerk/clerk-react'
import reportWebVitals from './reportWebVitals.ts'
import './styles.css'
import App from './App.tsx'

// Scream if clerk key is missing
const PUBLISHABLE_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY

if (!PUBLISHABLE_KEY) {
  throw new Error('Missing Publishable Key')
}

const rootElement = document.getElementById('app')
if (rootElement && !rootElement.innerHTML) {
  const root = ReactDOM.createRoot(rootElement)
  root.render(
    <StrictMode>
      <ClerkProvider publishableKey={PUBLISHABLE_KEY}>
        <App />
      </ClerkProvider>
    </StrictMode>,
  )
}

reportWebVitals()

// Register service worker (production only, handled by vite-plugin-pwa)
if ('serviceWorker' in navigator) {
  import('virtual:pwa-register').then(({ registerSW }) => {
    registerSW({
      immediate: false,
      onRegisteredSW(_swUrl, registration) {
        if (registration) {
          // Check for updates every hour
          setInterval(
            () => {
              registration.update()
            },
            60 * 60 * 1000,
          )
        }
      },
    })
  })
}
