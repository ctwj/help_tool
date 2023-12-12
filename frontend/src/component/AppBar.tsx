import Box from '@mui/material/Box';
import Stack from '@mui/material/Stack';
import { OpenInFull, CloseFullscreen, Clear, Minimize, Maximize } from '@mui/icons-material';

// 关闭按钮元素
function ItemButton ({ bgColor, borderColor, children }: { bgColor: string, borderColor: string, children: any }) {
    const css = {
        height: '18px',
        width: '18px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        background: bgColor,
        boxSizing: 'border-box',
        borderRadius: '100%',
        border: `1.32px solid ${borderColor}`,
        boxShadow: 3,
        ":hover": {
            boxShadow: 6,
        },
    }
    return (
        <Box sx={css}>
            {children}
        </Box>
    )
}

export default function AppBar() {
    return (
        <Box sx={{
            display: 'flex',
            position: 'absolute',
            alignItems: 'center',
            justifyContent: 'start',
            top: -25,
            left: 0,
        }}>
            <Stack direction="row" spacing={1}>
                <ItemButton bgColor="rgba(237, 106, 94, 0.8)" borderColor="rgba(206, 83, 71, 0.8)">
                    <Clear sx={{fontSize:'12px'}}/>
                </ItemButton>
                <ItemButton bgColor="rgba(246, 190, 79, 0.8)" borderColor="rgba(214, 162, 67, 0.8)" >
                    <Minimize sx={{fontSize:'12px'}} />
                </ItemButton>
                <ItemButton bgColor="rgba(98, 197, 84, 0.8)" borderColor="rgba(88, 169, 66, 0.8)" >
                    <OpenInFull sx={{fontSize:'12px'}}/>
                </ItemButton>
            </Stack>
        </Box>
    )
}