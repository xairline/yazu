import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'
import {ConfigProvider} from 'antd'
import {HashRouter} from "react-router-dom";

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <HashRouter basename={"/"}>
            <ConfigProvider theme={{token: {colorPrimary: '#006363'}}}>
                <App/>
            </ConfigProvider>
        </HashRouter>
    </React.StrictMode>
)
