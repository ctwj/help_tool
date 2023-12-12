import Box from '@mui/material/Box';

// import AppBar from './component/AppBar'
import ToolBar from './component/ToolBar'
import AppContent from './component/AppContent';

function App() {
    return (
        <Box sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            minHeight: '100vh',
            backgroundColor: 'rgba(255, 255, 255, 0.8)',
            border: '8px solid rgba(255, 255, 255, 0.45)',
            // borderTop: '32px solid rgba(255, 255, 255, 0.45)',
            backdropFilter: 'blur(147px);'
        }}>
            {/* <AppBar /> */}
            <ToolBar />
            <AppContent />
        </Box>
    )
}

export default App
