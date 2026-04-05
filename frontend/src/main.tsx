import React, { Suspense } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App'
import 'virtual:uno.css'
import './assets/styles/theme.css'
import './assets/styles/base.scss'
import { BrowserRouter, HashRouter } from 'react-router-dom'
import { initI18n } from './i18n'

const container = document.getElementById('root')

const root = createRoot(container!)

// Ensure i18n initialization is complete before rendering the application
initI18n().then(() => {
  root.render(
    <BrowserRouter>
      <Suspense>
        <App />
      </Suspense>
    </BrowserRouter>
  )
})
