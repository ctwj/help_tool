import { useNavigate, useLocation } from 'react-router-dom'

import { Box, Toolbar, AppBar, Typography, 
    IconButton, Stack, createTheme
} from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu';


export default function ToolBar() {
    const navigate = useNavigate();
    const location = useLocation();
    

    const toRouter = (router:string) => {
        navigate(router);
    }

    const MenuItem = (props: {name: string, to:string}) => {
        const {name, to} = props
        const isActive = location.pathname === to;

        return (
            <Typography 
                variant="h6" 
                color={isActive ? 'secondary' : 'inherit'}
                component="div" 
                onClick={() => toRouter(to)}>
                {name}
            </Typography>
        )
    }

    return (
        <Box sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'start',
            minHeight: '48px',
            flexShrink: '0',
            width: '100%',
        }}>
            <AppBar position="static">
                <Toolbar variant="dense">
                    <IconButton edge="start" color="inherit" aria-label="menu" sx={{ mr: 2 }}>
                    <MenuIcon />
                    </IconButton>

                    <Stack direction="row" spacing={1}>
                        <MenuItem name="Home" to="/" />
                        <MenuItem name="Other" to="/other" />
                    </Stack>
                    
                </Toolbar>
            </AppBar>
        </Box>
    )
}