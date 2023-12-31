import React from 'react'

import {createRoot} from 'react-dom/client'
import CssBaseline from '@mui/material/CssBaseline';
import '@fontsource/roboto/300.css'
import '@fontsource/roboto/400.css'
import '@fontsource/roboto/500.css'
import '@fontsource/roboto/700.css'
import './style.css'

import { BrowserRouter } from 'react-router-dom'

import App from './App'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <CssBaseline />
        <BrowserRouter>
            <App/>
        </BrowserRouter>
    </React.StrictMode>
)
