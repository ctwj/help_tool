import {useState} from 'react';
import {Greet} from "../../wailsjs/go/main/App";

import Button from '@mui/material/Button';
import Box from '@mui/material/Box';

export default function Main() {
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
        }}>
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <Button onClick={greet} variant="contained">Greet</Button>;
            </div>
        </Box>
    )
}