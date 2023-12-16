import React from 'react';
import { useLocation } from 'react-router-dom'

import {
    Typography, Box,
    Tabs, Tab,
    Fab, Zoom,
} from '@mui/material';
import { PlayArrow as PlayArrowIcon } from '@mui/icons-material'
import { useTheme } from '@mui/material/styles';
import { StartProxy } from '../../../wailsjs/go/proxy/Proxy'

interface TabPanelProps {
    children?: React.ReactNode;
    index: number;
    value: number;
}

function TabPanel(props: TabPanelProps) {
    const location = useLocation();
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{ p: 3 }}>
                    <Typography>{children}</Typography>
                </Box>
            )}
        </div>
    );
}

function a11yProps(index: number) {
    return {
        id: `simple-tab-${index}`,
        'aria-controls': `simple-tabpanel-${index}`,
    };
}

const fabStyle = {
    position: 'absolute',
    bottom: 16,
    right: 16,
};

const fab = {
    color: 'primary' as 'primary',
    sx: fabStyle,
    icon: <PlayArrowIcon />,
    label: 'Add',
}



export default function Main() {
    const [value, setValue] = React.useState(1)
    const theme = useTheme();

    const transitionDuration = {
        enter: theme.transitions.duration.enteringScreen,
        exit: theme.transitions.duration.leavingScreen,
    };

    const handleChange = (event: React.SyntheticEvent, newValue: string) => {
        setValue(Number(newValue));
    }

    return (
        <Box sx={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'start',
            justifyContent: 'start',
            flexGrow: 1,
            typography: 'body1',
            width: '100%',
            height: '100%',
        }}>
            <Box sx={{ borderBottom: 1, borderColor: 'divider', width: 1 }}>
                <Tabs value={value} onChange={handleChange} aria-label="basic tabs example">
                    <Tab label="Item One" {...a11yProps(0)} />
                    <Tab label="Item Two" {...a11yProps(1)} />
                    <Tab label="Item Three" {...a11yProps(2)} />
                </Tabs>
            </Box>
            <Box sx={{ flexGrow: 1, width: 1 }}>
                <TabPanel value={value} index={0}>
                    Item One
                </TabPanel>
                <TabPanel value={value} index={1}>
                    Item Two
                </TabPanel>
                <TabPanel value={value} index={2}>
                    Item Three
                </TabPanel>
            </Box>
            
            <Zoom
                key={fab.color}
                in={location.pathname === '/'}
                timeout={transitionDuration}
                style={{
                    transitionDelay: `${transitionDuration.exit}ms`,
                }}
                unmountOnExit
                >
                <Fab sx={fab.sx} aria-label={fab.label} color={fab.color}
                    onClick={StartProxy}>
                    {fab.icon}
                </Fab>
            </Zoom>
        </Box>
    )
}