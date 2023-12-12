import {useState} from 'react';
import {Greet} from "../../wailsjs/go/main/App";
import {Routes, Route} from 'react-router-dom'
import Button from '@mui/material/Button';
import Box from '@mui/material/Box';

import Main from '../pages/main/Main';
import Other from '../pages/other/Other';

export default function AppContent() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e: any) => setName(e.target.value);
    const updateResultText = (result: string) => setResultText(result);

    function greet() {
        Greet(name).then(updateResultText);
    }


    return (
        <Box sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'start',
            flexGrow: 1,
            width: '100%',
            height: '100%',
        }}>
            <Routes>
                <Route path="/" element={<Main />}></Route>    
                <Route path="/other" element={<Other />}></Route>    
            </Routes>
        </Box>
    )
}