import { Link } from 'react-router-dom'
import Box from '@mui/material/Box';

export default function ToolBar() {
    return (
        <Box sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'start',
            minHeight: '48px',
            flexShrink: '0',
            width: '100%',
        }}>
            <Link to="/">Home</Link>
            <Link to="/other">Other</Link>
        </Box>
    )
}