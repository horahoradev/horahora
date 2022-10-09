import { useEffect, useState } from "react";
import { useRouter } from "next/router";

import { Page } from "#components/page";
import { getArchivalEvents, getArchivalRequests } from "#api/archives";
import { type IArchivalEvent } from "#codegen/schema/001_interfaces";
import { CardList } from "#components/lists";
import { LoadingBar } from "#components/loading-bar";
import { EventCard } from "#entities/event";
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

function ArchivalEventsPage(params) {
  const [events, changeEvents] = useState<IArchivalEvent[]>();
  // I think this is a hack? looks okay to me though!
  const [timerVal, setTimerVal] = useState(0);

  const router = useRouter();
  const { query, isReady } = router;
  const { downloadID } = params;


  let download_id = query.download_id || downloadID;


  function reloadPage() {
    setTimerVal((timerVal) => timerVal + 1);
  }

  // TODO(ivan): Make a nicer page fetch hook that accounts for failure states
  useEffect(() => {
    let ignore = false;

    if (!isReady) {
      return
    }

    let fetchData = async () => {
      let subscriptionData = await getArchivalEvents(download_id);

      // videos.map((video, idx) => video.progress = videoInProgressDataset && videoInProgressDataset[idx] ? videoInProgressDataset[idx].progress : 0);

      // TODO: diff downloads in progress vs old downloads state, and unsubscribe!
      if (!ignore) {
        changeEvents(subscriptionData.ArchivalEvents);
      }
    };

    fetchData();
    return () => {
      ignore = true;
    };
  }, [timerVal, isReady]);

  useEffect(() => {
    const interval = setInterval(() => {
      reloadPage();
    }, 30000);

    return () => clearInterval(interval);
  }, []);


  return (
    <ThemeProvider theme={darkTheme}>
    <CssBaseline />
    <TableContainer component={Paper}>
    <Table sx={{ minWidth: 650 }} size="small" aria-label="a dense table">
        <TableHead>
        <TableRow>
            <TableCell align="left">Date</TableCell>
            <TableCell align="left">Video</TableCell>
            <TableCell align="left">Parent URL</TableCell>
            <TableCell align="left">Message</TableCell>
        </TableRow>
        </TableHead>
        <TableBody>
        {events ? events.map((row) => (
            <TableRow
            key={row.name}
            sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
            <TableCell align="left">{row.timestamp}</TableCell>
            <TableCell align="left">{row.video_url}</TableCell>
            <TableCell align="left">{row.parent_url}</TableCell>
            <TableCell align="left">{row.message}</TableCell>
            </TableRow>
        )): null }
        </TableBody>
    </Table>
    </TableContainer>
    </ThemeProvider>
  )
}

export default ArchivalEventsPage;
