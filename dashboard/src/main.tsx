import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { ModuleRegistry } from 'ag-grid-community';
import { ClientSideRowModelModule } from 'ag-grid-community';
import './index.css'
import App from './App.tsx'

// AG Grid 모듈 등록
ModuleRegistry.registerModules([ClientSideRowModelModule]);

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
